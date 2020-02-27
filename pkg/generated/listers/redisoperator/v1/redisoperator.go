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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RedisOperatorLister helps list RedisOperators.
type RedisOperatorLister interface {
	// List lists all RedisOperators in the indexer.
	List(selector labels.Selector) (ret []*v1.RedisOperator, err error)
	// RedisOperators returns an object that can list and get RedisOperators.
	RedisOperators(namespace string) RedisOperatorNamespaceLister
	RedisOperatorListerExpansion
}

// redisOperatorLister implements the RedisOperatorLister interface.
type redisOperatorLister struct {
	indexer cache.Indexer
}

// NewRedisOperatorLister returns a new RedisOperatorLister.
func NewRedisOperatorLister(indexer cache.Indexer) RedisOperatorLister {
	return &redisOperatorLister{indexer: indexer}
}

// List lists all RedisOperators in the indexer.
func (s *redisOperatorLister) List(selector labels.Selector) (ret []*v1.RedisOperator, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RedisOperator))
	})
	return ret, err
}

// RedisOperators returns an object that can list and get RedisOperators.
func (s *redisOperatorLister) RedisOperators(namespace string) RedisOperatorNamespaceLister {
	return redisOperatorNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RedisOperatorNamespaceLister helps list and get RedisOperators.
type RedisOperatorNamespaceLister interface {
	// List lists all RedisOperators in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.RedisOperator, err error)
	// Get retrieves the RedisOperator from the indexer for a given namespace and name.
	Get(name string) (*v1.RedisOperator, error)
	RedisOperatorNamespaceListerExpansion
}

// redisOperatorNamespaceLister implements the RedisOperatorNamespaceLister
// interface.
type redisOperatorNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RedisOperators in the indexer for a given namespace.
func (s redisOperatorNamespaceLister) List(selector labels.Selector) (ret []*v1.RedisOperator, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.RedisOperator))
	})
	return ret, err
}

// Get retrieves the RedisOperator from the indexer for a given namespace and name.
func (s redisOperatorNamespaceLister) Get(name string) (*v1.RedisOperator, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("redisoperator"), name)
	}
	return obj.(*v1.RedisOperator), nil
}
