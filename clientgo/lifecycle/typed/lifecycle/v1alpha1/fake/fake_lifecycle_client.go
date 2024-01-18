// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/ironcore-dev/lifecycle-manager/clientgo/lifecycle/typed/lifecycle/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeLifecycleV1alpha1 struct {
	*testing.Fake
}

func (c *FakeLifecycleV1alpha1) Machines(namespace string) v1alpha1.MachineInterface {
	return &FakeMachines{c, namespace}
}

func (c *FakeLifecycleV1alpha1) MachineTypes(namespace string) v1alpha1.MachineTypeInterface {
	return &FakeMachineTypes{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeLifecycleV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}