/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	v1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	scheme "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RedisOperatorsGetter has a method to return a RedisOperatorInterface.
// A group's client should implement this interface.
type RedisOperatorsGetter interface {
	RedisOperators(namespace string) RedisOperatorInterface
}

// RedisOperatorInterface has methods to work with RedisOperator resources.
type RedisOperatorInterface interface {
	Create(*v1.RedisOperator) (*v1.RedisOperator, error)
	Update(*v1.RedisOperator) (*v1.RedisOperator, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.RedisOperator, error)
	List(opts metav1.ListOptions) (*v1.RedisOperatorList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RedisOperator, err error)
	RedisOperatorExpansion
}

// redisOperators implements RedisOperatorInterface
type redisOperators struct {
	client rest.Interface
	ns     string
}

// newRedisOperators returns a RedisOperators
func newRedisOperators(c *NevercaseV1Client, namespace string) *redisOperators {
	return &redisOperators{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the redisOperator, and returns the corresponding redisOperator object, and an error if there is any.
func (c *redisOperators) Get(name string, options metav1.GetOptions) (result *v1.RedisOperator, err error) {
	result = &v1.RedisOperator{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisoperators").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RedisOperators that match those selectors.
func (c *redisOperators) List(opts metav1.ListOptions) (result *v1.RedisOperatorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.RedisOperatorList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested redisOperators.
func (c *redisOperators) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("redisoperators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a redisOperator and creates it.  Returns the server's representation of the redisOperator, and an error, if there is any.
func (c *redisOperators) Create(redisOperator *v1.RedisOperator) (result *v1.RedisOperator, err error) {
	result = &v1.RedisOperator{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("redisoperators").
		Body(redisOperator).
		Do().
		Into(result)
	return
}

// Update takes the representation of a redisOperator and updates it. Returns the server's representation of the redisOperator, and an error, if there is any.
func (c *redisOperators) Update(redisOperator *v1.RedisOperator) (result *v1.RedisOperator, err error) {
	result = &v1.RedisOperator{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redisoperators").
		Name(redisOperator.Name).
		Body(redisOperator).
		Do().
		Into(result)
	return
}

// Delete takes name of the redisOperator and deletes it. Returns an error if one occurs.
func (c *redisOperators) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisoperators").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *redisOperators) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisoperators").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched redisOperator.
func (c *redisOperators) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RedisOperator, err error) {
	result = &v1.RedisOperator{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("redisoperators").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
