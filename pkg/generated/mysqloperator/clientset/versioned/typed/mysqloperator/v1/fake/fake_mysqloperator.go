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

package fake

import (
	"context"

	mysqloperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/mysqloperator/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMysqlOperators implements MysqlOperatorInterface
type FakeMysqlOperators struct {
	Fake *FakeNevercaseV1
	ns   string
}

var mysqloperatorsResource = schema.GroupVersionResource{Group: "nevercase.io", Version: "v1", Resource: "mysqloperators"}

var mysqloperatorsKind = schema.GroupVersionKind{Group: "nevercase.io", Version: "v1", Kind: "MysqlOperator"}

// Get takes name of the mysqlOperator, and returns the corresponding mysqlOperator object, and an error if there is any.
func (c *FakeMysqlOperators) Get(ctx context.Context, name string, options v1.GetOptions) (result *mysqloperatorv1.MysqlOperator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(mysqloperatorsResource, c.ns, name), &mysqloperatorv1.MysqlOperator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mysqloperatorv1.MysqlOperator), err
}

// List takes label and field selectors, and returns the list of MysqlOperators that match those selectors.
func (c *FakeMysqlOperators) List(ctx context.Context, opts v1.ListOptions) (result *mysqloperatorv1.MysqlOperatorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(mysqloperatorsResource, mysqloperatorsKind, c.ns, opts), &mysqloperatorv1.MysqlOperatorList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &mysqloperatorv1.MysqlOperatorList{ListMeta: obj.(*mysqloperatorv1.MysqlOperatorList).ListMeta}
	for _, item := range obj.(*mysqloperatorv1.MysqlOperatorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested mysqlOperators.
func (c *FakeMysqlOperators) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(mysqloperatorsResource, c.ns, opts))

}

// Create takes the representation of a mysqlOperator and creates it.  Returns the server's representation of the mysqlOperator, and an error, if there is any.
func (c *FakeMysqlOperators) Create(ctx context.Context, mysqlOperator *mysqloperatorv1.MysqlOperator, opts v1.CreateOptions) (result *mysqloperatorv1.MysqlOperator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(mysqloperatorsResource, c.ns, mysqlOperator), &mysqloperatorv1.MysqlOperator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mysqloperatorv1.MysqlOperator), err
}

// Update takes the representation of a mysqlOperator and updates it. Returns the server's representation of the mysqlOperator, and an error, if there is any.
func (c *FakeMysqlOperators) Update(ctx context.Context, mysqlOperator *mysqloperatorv1.MysqlOperator, opts v1.UpdateOptions) (result *mysqloperatorv1.MysqlOperator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(mysqloperatorsResource, c.ns, mysqlOperator), &mysqloperatorv1.MysqlOperator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mysqloperatorv1.MysqlOperator), err
}

// Delete takes name of the mysqlOperator and deletes it. Returns an error if one occurs.
func (c *FakeMysqlOperators) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(mysqloperatorsResource, c.ns, name), &mysqloperatorv1.MysqlOperator{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMysqlOperators) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(mysqloperatorsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &mysqloperatorv1.MysqlOperatorList{})
	return err
}

// Patch applies the patch and returns the patched mysqlOperator.
func (c *FakeMysqlOperators) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *mysqloperatorv1.MysqlOperator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(mysqloperatorsResource, c.ns, name, pt, data, subresources...), &mysqloperatorv1.MysqlOperator{})

	if obj == nil {
		return nil, err
	}
	return obj.(*mysqloperatorv1.MysqlOperator), err
}
