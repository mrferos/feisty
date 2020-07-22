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
	"github.com/go-logr/logr"
	"github.com/mrferos/feisty/constants"
	"github.com/mrferos/feisty/revisions"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	feistyv1alpha1 "github.com/mrferos/feisty/api/v1alpha1"
)

var (
	defaultExposedPort             = int32(80)
	restartDeploymentAnnotationKey = constants.FeistyAnnotationPrefix + "restart-time"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=feisty.paas.feisty.dev,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=feisty.paas.feisty.dev,resources=applications/status,verbs=get;update;patch

func getAppLabels(app feistyv1alpha1.Application) map[string]string {
	return map[string]string{
		"app": app.Name,
	}
}

func (r *ApplicationReconciler) upsertDeployment(app feistyv1alpha1.Application, req ctrl.Request, ctx context.Context) (ctrl.Result, error) {
	log := r.Log.WithValues("application", req.NamespacedName)
	appLabels := getAppLabels(app)

	replicas := int32(0)
	if app.Spec.Replicas > 0 {
		replicas = int32(app.Spec.Replicas)
	}

	doCreate := false
	var deployment v1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &deployment); err != nil {
		doCreate = true
		deployment = v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
			},
			Spec: v1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: appLabels,
				},
				Template: v12.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:   app.Name,
						Labels: appLabels,
					},
					Spec: v12.PodSpec{
						Containers: []v12.Container{{
							Name: app.Name,
						}},
					},
				},
			},
		}
	}

	if app.Spec.RestartTime != "" {
		if deployment.Spec.Template.ObjectMeta.Annotations == nil {
			deployment.Spec.Template.ObjectMeta.Annotations = map[string]string{}
		}

		deployment.Spec.Template.ObjectMeta.Annotations[restartDeploymentAnnotationKey] = app.Spec.RestartTime
	}

	deployment.Spec.Replicas = &replicas
	deployment.Spec.Template.Spec.Containers[0].Image = app.Spec.Image

	if app.Spec.Port != 0 {
		deployment.Spec.Template.Spec.Containers[0].Ports = []v12.ContainerPort{{
			Name:          "http",
			ContainerPort: int32(app.Spec.Port),
		}}
	}

	if app.Spec.AppConfigRef != "" {
		deployment.Spec.Template.Spec.Containers[0].EnvFrom = []v12.EnvFromSource{{
			SecretRef: &v12.SecretEnvSource{
				LocalObjectReference: v12.LocalObjectReference{
					Name: app.Spec.AppConfigRef,
				},
			},
		}}
	}

	if doCreate {
		_ = ctrl.SetControllerReference(&app, &deployment, r.Scheme)
		if err := r.Create(ctx, &deployment); err != nil {
			log.Error(err, "Could not create deployment")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Update(ctx, &deployment); err != nil {
			log.Error(err, "Could not update deployment")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) upsertService(app feistyv1alpha1.Application, req ctrl.Request, ctx context.Context) (ctrl.Result, error) {
	log := r.Log.WithValues("application", req.NamespacedName)
	appLabels := getAppLabels(app)

	doCreate := false
	var svc v12.Service
	if err := r.Get(ctx, req.NamespacedName, &svc); err != nil {
		doCreate = true
		svc = v12.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
			},
			Spec: v12.ServiceSpec{
				Ports: []v12.ServicePort{{
					Name:     "http",
					Protocol: "TCP",
					Port:     defaultExposedPort,
					TargetPort: intstr.IntOrString{
						IntVal: int32(app.Spec.Port),
					},
				}},
				Selector: appLabels,
				Type:     "ClusterIP",
			},
		}
	}

	if doCreate {
		_ = ctrl.SetControllerReference(&app, &svc, r.Scheme)
		if err := r.Create(ctx, &svc); err != nil {
			log.Error(err, "Could not create svc")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Update(ctx, &svc); err != nil {
			log.Error(err, "Could not update svc")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) upsertIngress(app feistyv1alpha1.Application, req ctrl.Request, ctx context.Context) (ctrl.Result, error) {
	log := r.Log.WithValues("application", req.NamespacedName)

	doCreate := false
	var ingress v1beta1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ingress); err != nil {
		doCreate = true
		ingress = v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
			},
		}
	}

	// TODO: add code here to deal with default domain
	var rules []v1beta1.IngressRule
	domains := app.Spec.Domains
	for _, domain := range domains {
		rules = append(rules, v1beta1.IngressRule{
			Host: domain.Host,
			IngressRuleValue: v1beta1.IngressRuleValue{
				HTTP: &v1beta1.HTTPIngressRuleValue{
					Paths: []v1beta1.HTTPIngressPath{{
						Path: "/",
						Backend: v1beta1.IngressBackend{
							ServiceName: app.Name,
							ServicePort: intstr.IntOrString{
								IntVal: defaultExposedPort,
							},
						},
					}},
				},
			},
		})
	}

	ingress.Spec.Rules = rules

	if doCreate {
		_ = ctrl.SetControllerReference(&app, &ingress, r.Scheme)
		if err := r.Create(ctx, &ingress); err != nil {
			log.Error(err, "Could not create ingress")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Update(ctx, &ingress); err != nil {
			log.Error(err, "Could not update ingress")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("application", req.NamespacedName)
	rev := revisions.Revision{
		Client: r.Client,
		Log:    r.Log,
	}

	var app feistyv1alpha1.Application
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		log.Error(err, "Unable to fetch Application")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploymentExist := false
	if app.Spec.Image == "" {
		log.Info("No deployment action taken because no image was supplied")
		return ctrl.Result{}, nil
	} else {
		if res, err := r.upsertDeployment(app, req, ctx); err != nil {
			log.Error(err, "There was an error doing deployment handling")
			return res, err
		} else {
			deploymentExist = true
		}
	}

	svcExists := false
	if app.Spec.Port != 0 && deploymentExist {
		if res, err := r.upsertService(app, req, ctx); err != nil {
			log.Error(err, "There was an error doing service handling")
			return res, err
		} else {
			svcExists = true
		}
	}

	if app.Spec.RoutingEnabled && svcExists {
		if res, err := r.upsertIngress(app, req, ctx); err != nil {
			log.Error(err, "There was an error doing ingress handling")
			return res, err
		}
	}

	_ = rev.CreateIfNeeded(req.NamespacedName, ctx)

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&feistyv1alpha1.Application{}).
		Owns(&v1.Deployment{}).
		WithEventFilter(predicate.Funcs{
			UpdateFunc: revisions.RevisionWatchFilter,
		}).
		Complete(r)
}
