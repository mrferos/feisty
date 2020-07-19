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
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	feistyv1alpha1 "github.com/mrferos/feisty/api/v1alpha1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=feisty.paas.fesity.dev,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=feisty.paas.fesity.dev,resources=applications/status,verbs=get;update;patch

func (r *ApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("application", req.NamespacedName)

	var app feistyv1alpha1.Application
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		log.Error(err, "Unable to fetch Application")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if app.Spec.Image == "" {
		log.Info("No deployment action taken because no image was supplied")
		return ctrl.Result{}, nil
	}

	appLabel := map[string]string{
		"app": app.Name,
	}

	replicas := int32(1)
	deployment := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: appLabel,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   app.Name,
					Labels: appLabel,
				},
				Spec: v12.PodSpec{
					Containers: []v12.Container{{
						Name:  app.Name,
						Image: app.Spec.Image,
					}},
				},
			},
		},
	}

	if app.Spec.Port != 0 {
		deployment.Spec.Template.Spec.Containers[0].Ports = []v12.ContainerPort{{
			Name:          "http",
			ContainerPort: 80,
		}}
	}

	_ = ctrl.SetControllerReference(&app, &deployment, r.Scheme)

	if err := r.Create(ctx, &deployment); err != nil {
		log.Error(err, "Could not create deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&feistyv1alpha1.Application{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}
