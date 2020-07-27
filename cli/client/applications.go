package client

import (
	"github.com/mrferos/feisty/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var appResource = "applications"

type ApplicationInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ApplicationList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Application, error)
	Create(application *v1alpha1.Application) (*v1alpha1.Application, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Update(application *v1alpha1.Application) (*v1alpha1.Application, error)
}

type applicationClient struct {
	restClient rest.Interface
	ns         string
}

func (c *applicationClient) List(opts metav1.ListOptions) (*v1alpha1.ApplicationList, error) {
	result := v1alpha1.ApplicationList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(appResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *applicationClient) Get(name string, options metav1.GetOptions) (*v1alpha1.Application, error) {
	result := v1alpha1.Application{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(appResource).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *applicationClient) Create(application *v1alpha1.Application) (*v1alpha1.Application, error) {
	result := v1alpha1.Application{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(appResource).
		Body(application).
		Do().
		Into(&result)

	return &result, err
}

func (c *applicationClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(appResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func (c *applicationClient) Update(application *v1alpha1.Application) (*v1alpha1.Application, error) {
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Name(application.Name).
		Resource(appResource).
		Body(application).
		Do().
		Into(application)

	return application, err
}
