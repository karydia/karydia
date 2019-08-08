// Copyright (C) 2019 SAP SE or an SAP affiliate company. All rights reserved.
// This file is licensed under the Apache Software License, v. 2 except as
// noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	karydiav1alpha1 "github.com/karydia/karydia/pkg/apis/karydia/v1alpha1"
	versioned "github.com/karydia/karydia/pkg/client/clientset/versioned"
	internalinterfaces "github.com/karydia/karydia/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/karydia/karydia/pkg/client/listers/karydia/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// KarydiaNetworkPolicyInformer provides access to a shared informer and lister for
// KarydiaNetworkPolicies.
type KarydiaNetworkPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.KarydiaNetworkPolicyLister
}

type karydiaNetworkPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewKarydiaNetworkPolicyInformer constructs a new informer for KarydiaNetworkPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKarydiaNetworkPolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKarydiaNetworkPolicyInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredKarydiaNetworkPolicyInformer constructs a new informer for KarydiaNetworkPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKarydiaNetworkPolicyInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KarydiaV1alpha1().KarydiaNetworkPolicies().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KarydiaV1alpha1().KarydiaNetworkPolicies().Watch(options)
			},
		},
		&karydiav1alpha1.KarydiaNetworkPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *karydiaNetworkPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKarydiaNetworkPolicyInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *karydiaNetworkPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&karydiav1alpha1.KarydiaNetworkPolicy{}, f.defaultInformer)
}

func (f *karydiaNetworkPolicyInformer) Lister() v1alpha1.KarydiaNetworkPolicyLister {
	return v1alpha1.NewKarydiaNetworkPolicyLister(f.Informer().GetIndexer())
}
