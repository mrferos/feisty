package client

import (
	"github.com/mrferos/feisty/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var appConfigResource = "applicationconfigs"

type ApplicationConfigInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ApplicationConfigList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.ApplicationConfig, error)
	Create(applicationConfig *v1alpha1.ApplicationConfig) (*v1alpha1.ApplicationConfig, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Update(applicationConfig *v1alpha1.ApplicationConfig) (*v1alpha1.ApplicationConfig, error)
}

type applicationConfigClient struct {
	restClient rest.Interface
	ns         string
}

func (c *applicationConfigClient) List(opts metav1.ListOptions) (*v1alpha1.ApplicationConfigList, error) {
	result := v1alpha1.ApplicationConfigList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(appConfigResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *applicationConfigClient) Get(name string, options metav1.GetOptions) (*v1alpha1.ApplicationConfig, error) {
	result := v1alpha1.ApplicationConfig{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(appConfigResource).
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *applicationConfigClient) Create(applicationConfig *v1alpha1.ApplicationConfig) (*v1alpha1.ApplicationConfig, error) {
	result := v1alpha1.ApplicationConfig{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(appConfigResource).
		Body(applicationConfig).
		Do().
		Into(&result)

	return &result, err
}

func (c *applicationConfigClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(appConfigResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

func (c *applicationConfigClient) Update(applicationConfig *v1alpha1.ApplicationConfig) (*v1alpha1.ApplicationConfig, error) {
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Name(applicationConfig.Name).
		Resource(appConfigResource).
		Body(applicationConfig).
		Do().
		Into(applicationConfig)

	return applicationConfig, err
}
