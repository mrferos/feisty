/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	feistyv1alpha1 "github.com/mrferos/feisty/api/v1alpha1"
)

// ApplicationConfigReconciler reconciles a ApplicationConfig object
type ApplicationConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=feisty.paas.feisty.dev,resources=applicationconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=feisty.paas.feisty.dev,resources=applicationconfigs/status,verbs=get;update;patch

func (r *ApplicationConfigReconciler) getSecretId(cfg feistyv1alpha1.ApplicationConfig) (string, error) {
	keyValueJson, err := json.Marshal(cfg.Spec.KeyValuePairs)
	if err != nil {
		return "", err
	}

	md5bytes := md5.Sum([]byte(keyValueJson))

	return fmt.Sprintf("%x", md5bytes), nil
}

func (r *ApplicationConfigReconciler) upsertSecret(cfg feistyv1alpha1.ApplicationConfig, req ctrl.Request, ctx context.Context) (v1.Secret, error) {
	log := r.Log.WithValues("application", req.NamespacedName)

	secretId, err := r.getSecretId(cfg)
	if err != nil {
		return v1.Secret{}, err
	}

	secretName := cfg.Name + "-" + secretId
	objKey := client.ObjectKey{
		Namespace: cfg.Namespace,
		Name:      secretName,
	}

	doCreate := false
	var secret v1.Secret
	if err := r.Get(ctx, objKey, &secret); err != nil {
		doCreate = true
		secret = v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: cfg.Namespace,
			},
		}
	}

	secretData := map[string][]byte{}
	for k, v := range cfg.Spec.KeyValuePairs {
		secretData[k] = []byte(v)
	}

	secret.Data = secretData

	if doCreate {
		_ = ctrl.SetControllerReference(&cfg, &secret, r.Scheme)
		if err := r.Create(ctx, &secret); err != nil {
			log.Error(err, "Could not create secret")
			return secret, err
		}
	} else {
		if err := r.Update(ctx, &secret); err != nil {
			log.Error(err, "Could not update secret")
			return secret, err
		}
	}

	return secret, nil
}

func (r *ApplicationConfigReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("applicationconfig", req.NamespacedName)

	var cfg feistyv1alpha1.ApplicationConfig
	if err := r.Get(ctx, req.NamespacedName, &cfg); err != nil {
		log.Error(err, "Unable to fetch ApplicationConfig")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if secret, err := r.upsertSecret(cfg, req, ctx); err != nil {
		log.Error(err, "There was an error doing secret handling")
		return ctrl.Result{}, err
	} else {
		var app feistyv1alpha1.Application
		if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
			log.Error(err, "Unable to fetch Application")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		app.Spec.AppConfigRef = secret.Name
		if err := r.Update(ctx, &app); err != nil {
			log.Error(err, "Could not update Application with target ApplicationConfig")
		}
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&feistyv1alpha1.ApplicationConfig{}).
		Owns(&v1.Secret{}).
		Complete(r)
}
