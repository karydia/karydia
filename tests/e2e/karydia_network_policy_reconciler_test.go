// Copyright 2019 Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
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

package e2e

import (
	"bytes"
	"testing"
	"time"

	"github.com/karydia/karydia/pkg/apis/karydia/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultNetworkPolicyName   = "karydia-default-network-policy"
	defaultNetworkPolicyL2Name = "karydia-default-network-policy-l2"
)

func TestCreateKarydiaNetworkPolicyForNewNamespace(t *testing.T) {
	defaultNetworkPolicy := &networkingv1.NetworkPolicy{}
	defaultNetworkPolicy.Name = defaultNetworkPolicyName
	defaultNetworkPolicy.Spec = networkingv1.NetworkPolicySpec{
		PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeEgress},
	}
	namespace, err := f.CreateTestNamespace()
	if err != nil {
		t.Fatalf("failed to create test namespace: %v", err)
	}

	timeout := 3000 * time.Millisecond
	if err := f.WaitNetworkPolicyCreatedCreated(namespace.GetName(), defaultNetworkPolicyName, timeout); err != nil {
		t.Fatalf("failed to create default network policy for new namespace: %v", err)
	}

	namespaceNetworkPolicy, err := f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Get(defaultNetworkPolicyName, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to create default network policy for new namespace: %v", err)
	}

	if !networkPoliciesAreEqual(namespaceNetworkPolicy, defaultNetworkPolicy) {
		t.Fatalf("Network policy for created namespace is not equal to the default network policy: %v", err)
	}

	err = f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Delete(defaultNetworkPolicyName, &meta_v1.DeleteOptions{})
	if err != nil {
		t.Fatalf("failed to delete default network policy for new namespace: %v", err)
	}

	if err := f.WaitNetworkPolicyCreatedCreated(namespace.GetName(), defaultNetworkPolicyName, timeout); err != nil {
		t.Fatalf("failed to create default network policy for new namespace: %v", err)
	}

	namespaceNetworkPolicy, err = f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Get(defaultNetworkPolicyName, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("Reconciler failed to recreate default network policy for new namespace: %v", err)
	}
	//Update NP
	namespaceNetworkPolicy.Spec = networkingv1.NetworkPolicySpec{
		PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress}}

	updatedNetworkPolicy, err := f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Update(namespaceNetworkPolicy)
	if err != nil {
		t.Fatalf("failed to update default network policy for new namespace: %v", err)
	}

	duration := 3 * time.Second
	time.Sleep(duration)

	updatedNetworkPolicy, err = f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Get(defaultNetworkPolicyName, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get default network policy for new namespace: %v", err)
	}
	if !networkPoliciesAreEqual(updatedNetworkPolicy, defaultNetworkPolicy) {
		t.Fatalf("Reconcilation failed after network policy has changed ")
	}
}

func TestCreateKarydiaNetworkPolicyForAnnotatedNamespace(t *testing.T) {
	defaultNetworkPolicy := &networkingv1.NetworkPolicy{}
	defaultNetworkPolicy.Name = defaultNetworkPolicyL2Name
	defaultNetworkPolicy.Spec = networkingv1.NetworkPolicySpec{
		PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
	}

	defaultKarydiaNetworkPolicyL2 := &v1alpha1.KarydiaNetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "karydia-e2e-test",
			},
			Name: defaultNetworkPolicyL2Name,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PolicyTypes: []networkingv1.PolicyType{networkingv1.PolicyTypeIngress},
		},
	}

	_, err := f.KarydiaClientset.KarydiaV1alpha1().KarydiaNetworkPolicies().Create(defaultKarydiaNetworkPolicyL2)
	if err != nil {
		t.Fatalf("failed to create: %v", defaultKarydiaNetworkPolicyL2)
	}

	annotations := make(map[string]string)
	annotations["karydia.gardener.cloud/networkPolicy"] = defaultNetworkPolicyL2Name
	namespace, err := f.CreateTestNamespaceWithAnnotation(annotations)
	if err != nil {
		t.Fatalf("failed to create test namespace: %v", err)
	}

	timeout := 3000 * time.Millisecond
	if err := f.WaitNetworkPolicyCreatedCreated(namespace.GetName(), defaultNetworkPolicyL2Name, timeout); err != nil {
		t.Fatalf("failed to create default network policy for new namespace: %v", err)
	}

	_, err = f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Get(defaultNetworkPolicyName, meta_v1.GetOptions{})
	if err == nil {
		t.Fatalf("Default level 1 network policy should not be found")
	}

	namespaceNetworkPolicy, err := f.KubeClientset.NetworkingV1().NetworkPolicies(namespace.GetName()).Get(defaultNetworkPolicyL2Name, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to create default network policy for new namespace: %v", err)
	}

	if !networkPoliciesAreEqual(namespaceNetworkPolicy, defaultNetworkPolicy) {
		t.Fatalf("Network policy for created namespace is not equal to the default network policy: %v", err)
	}
}
func TestGetKarydiaNetworkPolicyForExcludedNamespace(t *testing.T) {
	if _, err := f.KubeClientset.NetworkingV1().NetworkPolicies("kube-system").Get(defaultNetworkPolicyName, meta_v1.GetOptions{}); err == nil {
		t.Fatalf("Default network policy should not be found for excluded namespace ")
	}

}

func networkPoliciesAreEqual(np1, np2 *networkingv1.NetworkPolicy) bool {
	np1Spec, _ := np1.Spec.Marshal()
	np2Spec, _ := np2.Spec.Marshal()
	return bytes.Equal(np1Spec, np2Spec)

}