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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	redisoperatorv1 "github.com/nevercase/k8s-controller-custom-resource/pkg/apis/redisoperator/v1"
	versioned "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/clientset/versioned"
	internalinterfaces "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/informers/externalversions/internalinterfaces"
	v1 "github.com/nevercase/k8s-controller-custom-resource/pkg/generated/redisoperator/listers/redisoperator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// RedisOperatorInformer provides access to a shared informer and lister for
// RedisOperators.
type RedisOperatorInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.RedisOperatorLister
}

type redisOperatorInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRedisOperatorInformer constructs a new informer for RedisOperator type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRedisOperatorInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRedisOperatorInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRedisOperatorInformer constructs a new informer for RedisOperator type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRedisOperatorInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NevercaseV1().RedisOperators(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NevercaseV1().RedisOperators(namespace).Watch(options)
			},
		},
		&redisoperatorv1.RedisOperator{},
		resyncPeriod,
		indexers,
	)
}

func (f *redisOperatorInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRedisOperatorInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *redisOperatorInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&redisoperatorv1.RedisOperator{}, f.defaultInformer)
}

func (f *redisOperatorInformer) Lister() v1.RedisOperatorLister {
	return v1.NewRedisOperatorLister(f.Informer().GetIndexer())
}
