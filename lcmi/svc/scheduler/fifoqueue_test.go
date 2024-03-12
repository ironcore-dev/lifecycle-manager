// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"strconv"
	"sync"

	lifecyclev1alpha1 "github.com/ironcore-dev/lifecycle-manager/api/lifecycle/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FIFO Queue", func() {
	Context("On push", func() {
		It("Should not allow to enqueue more items than queue capacity", func() {
			var ok bool
			queue := NewFIFOQueue[*lifecyclev1alpha1.Machine](2)
			ok = queue.Push(Task[*lifecyclev1alpha1.Machine]{Key: "1"})
			Expect(ok).To(BeTrue())
			Expect(queue.Has("1")).To(BeTrue())
			ok = queue.Push(Task[*lifecyclev1alpha1.Machine]{Key: "2"})
			Expect(ok).To(BeTrue())
			Expect(queue.Has("2")).To(BeTrue())
			Expect(queue.IsFull()).To(BeTrue())
			ok = queue.Push(Task[*lifecyclev1alpha1.Machine]{Key: "3"})
			Expect(ok).To(BeFalse())
			Expect(queue.Has("3")).To(BeFalse())
		})
	})

	Context("On pop", func() {
		It("Should pop from head", func() {
			var (
				item Task[*lifecyclev1alpha1.Machine]
				ok   bool
			)
			queue := NewFIFOQueue[*lifecyclev1alpha1.Machine](2)
			ok = queue.Push(Task[*lifecyclev1alpha1.Machine]{Key: "1"})
			Expect(ok).To(BeTrue())
			Expect(queue.Has("1")).To(BeTrue())
			ok = queue.Push(Task[*lifecyclev1alpha1.Machine]{Key: "2"})
			Expect(ok).To(BeTrue())
			Expect(queue.Has("2")).To(BeTrue())

			item, ok = queue.Pop()
			Expect(ok).To(BeTrue())
			Expect(item.Key).To(Equal("1"))
			item, ok = queue.Pop()
			Expect(ok).To(BeTrue())
			Expect(item.Key).To(Equal("2"))
			item, ok = queue.Pop()
			Expect(ok).To(BeFalse())
			Expect(item.Key).To(BeEmpty())
		})
	})

	Context("On parallel push/pop", func() {
		It("Should succeed", func() {
			var (
				wg sync.WaitGroup
				ok bool
			)
			queue := NewFIFOQueue[*lifecyclev1alpha1.Machine](200)
			for i := range 200 {
				wg.Add(1)
				key := i
				go func() {
					ok = queue.TryPush(Task[*lifecyclev1alpha1.Machine]{Key: strconv.Itoa(key)})
					Expect(ok).To(BeTrue())
					wg.Done()
				}()
			}
			wg.Wait()
			Expect(queue.IsFull()).To(BeTrue())

			for range 200 {
				wg.Add(1)
				go func() {
					_, ok = queue.Pop()
					Expect(ok).To(BeTrue())
					wg.Done()
				}()
			}
			wg.Wait()
			Expect(queue.IsEmpty()).To(BeTrue())
		})
	})
})
