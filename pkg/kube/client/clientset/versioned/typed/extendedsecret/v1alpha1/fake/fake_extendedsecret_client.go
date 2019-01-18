/*

Don't alter this file, it was generated.

*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "code.cloudfoundry.org/cf-operator/pkg/kube/client/clientset/versioned/typed/extendedsecret/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeExtendedsecretV1alpha1 struct {
	*testing.Fake
}

func (c *FakeExtendedsecretV1alpha1) ExtendedSecrets(namespace string) v1alpha1.ExtendedSecretInterface {
	return &FakeExtendedSecrets{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeExtendedsecretV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
