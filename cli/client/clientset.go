package client

import (
	"github.com/mrferos/feisty/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FeistyV1Alpha1Interface interface {
	Applications(namespace string) ApplicationInterface
	ApplicationConfigs(namespace string) ApplicationConfigInterface
}

type FeistyV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*FeistyV1Alpha1Client, error) {
	config := *c
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()
	config.ContentConfig.GroupVersion = &schema.GroupVersion{
		Group:   v1alpha1.GroupVersion.Group,
		Version: v1alpha1.GroupVersion.Version,
	}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &FeistyV1Alpha1Client{restClient: client}, nil
}

func (c *FeistyV1Alpha1Client) Applications(namespace string) ApplicationInterface {
	return &applicationClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *FeistyV1Alpha1Client) ApplicationConfigs(namespace string) ApplicationConfigInterface {
	return &applicationConfigClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}