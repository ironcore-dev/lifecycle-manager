// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/clientgo/applyconfiguration/lifecycle/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMachineTypes implements MachineTypeInterface
type FakeMachineTypes struct {
	Fake *FakeLifecycleV1alpha1
	ns   string
}

var machinetypesResource = v1alpha1.SchemeGroupVersion.WithResource("machinetypes")

var machinetypesKind = v1alpha1.SchemeGroupVersion.WithKind("MachineType")

// Get takes name of the machineType, and returns the corresponding machineType object, and an error if there is any.
func (c *FakeMachineTypes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MachineType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(machinetypesResource, c.ns, name), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}

// List takes label and field selectors, and returns the list of MachineTypes that match those selectors.
func (c *FakeMachineTypes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MachineTypeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(machinetypesResource, machinetypesKind, c.ns, opts), &v1alpha1.MachineTypeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MachineTypeList{ListMeta: obj.(*v1alpha1.MachineTypeList).ListMeta}
	for _, item := range obj.(*v1alpha1.MachineTypeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested machineTypes.
func (c *FakeMachineTypes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(machinetypesResource, c.ns, opts))

}

// Create takes the representation of a machineType and creates it.  Returns the server's representation of the machineType, and an error, if there is any.
func (c *FakeMachineTypes) Create(ctx context.Context, machineType *v1alpha1.MachineType, opts v1.CreateOptions) (result *v1alpha1.MachineType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(machinetypesResource, c.ns, machineType), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}

// Update takes the representation of a machineType and updates it. Returns the server's representation of the machineType, and an error, if there is any.
func (c *FakeMachineTypes) Update(ctx context.Context, machineType *v1alpha1.MachineType, opts v1.UpdateOptions) (result *v1alpha1.MachineType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(machinetypesResource, c.ns, machineType), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeMachineTypes) UpdateStatus(ctx context.Context, machineType *v1alpha1.MachineType, opts v1.UpdateOptions) (*v1alpha1.MachineType, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(machinetypesResource, "status", c.ns, machineType), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}

// Delete takes name of the machineType and deletes it. Returns an error if one occurs.
func (c *FakeMachineTypes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(machinetypesResource, c.ns, name, opts), &v1alpha1.MachineType{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMachineTypes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(machinetypesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.MachineTypeList{})
	return err
}

// Patch applies the patch and returns the patched machineType.
func (c *FakeMachineTypes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MachineType, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(machinetypesResource, c.ns, name, pt, data, subresources...), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied machineType.
func (c *FakeMachineTypes) Apply(ctx context.Context, machineType *lifecyclev1alpha1.MachineTypeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MachineType, err error) {
	if machineType == nil {
		return nil, fmt.Errorf("machineType provided to Apply must not be nil")
	}
	data, err := json.Marshal(machineType)
	if err != nil {
		return nil, err
	}
	name := machineType.Name
	if name == nil {
		return nil, fmt.Errorf("machineType.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(machinetypesResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeMachineTypes) ApplyStatus(ctx context.Context, machineType *lifecyclev1alpha1.MachineTypeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MachineType, err error) {
	if machineType == nil {
		return nil, fmt.Errorf("machineType provided to Apply must not be nil")
	}
	data, err := json.Marshal(machineType)
	if err != nil {
		return nil, err
	}
	name := machineType.Name
	if name == nil {
		return nil, fmt.Errorf("machineType.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(machinetypesResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.MachineType{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MachineType), err
}
