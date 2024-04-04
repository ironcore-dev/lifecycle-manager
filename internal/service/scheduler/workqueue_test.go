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

var _ = Describe("Ring Buffer Queue", func() {
	Context("On enqueue", func() {
		It("Should not allow to enqueue more items than queue capacity", func() {
			var ok bool
			queue := NewRingBufQueue[*lifecyclev1alpha1.Machine](2)
			ok = queue.Enqueue(Task[*lifecyclev1alpha1.Machine]{Key: "1"})
			<-queue.Enqueued
			Expect(ok).To(BeTrue())
			Expect(queue.Has("1")).To(BeTrue())
			ok = queue.Enqueue(Task[*lifecyclev1alpha1.Machine]{Key: "2"})
			<-queue.Enqueued
			Expect(ok).To(BeTrue())
			Expect(queue.Has("2")).To(BeTrue())
			Expect(queue.IsFull()).To(BeTrue())
			ok = queue.Enqueue(Task[*lifecyclev1alpha1.Machine]{Key: "3"})
			Expect(ok).To(BeFalse())
			Expect(queue.Has("3")).To(BeFalse())
		})
	})

	Context("On dequeue", func() {
		It("Should dequeue from head", func() {
			var (
				item Task[*lifecyclev1alpha1.Machine]
				ok   bool
			)
			queue := NewRingBufQueue[*lifecyclev1alpha1.Machine](2)
			ok = queue.Enqueue(Task[*lifecyclev1alpha1.Machine]{Key: "1"})
			<-queue.Enqueued
			Expect(ok).To(BeTrue())
			Expect(queue.Has("1")).To(BeTrue())
			ok = queue.Enqueue(Task[*lifecyclev1alpha1.Machine]{Key: "2"})
			<-queue.Enqueued
			Expect(ok).To(BeTrue())
			Expect(queue.Has("2")).To(BeTrue())

			item, ok = queue.Dequeue()
			Expect(ok).To(BeTrue())
			Expect(item.Key).To(Equal("1"))
			item, ok = queue.Dequeue()
			Expect(ok).To(BeTrue())
			Expect(item.Key).To(Equal("2"))
			item, ok = queue.Dequeue()
			Expect(ok).To(BeFalse())
			Expect(item.Key).To(BeEmpty())
		})
	})

	Context("On parallel enqueue/dequeue", func() {
		It("Should succeed", func() {
			var (
				wg sync.WaitGroup
				ok bool
			)
			queue := NewRingBufQueue[*lifecyclev1alpha1.Machine](20)
			for i := range 20 {
				wg.Add(1)
				key := i
				go func() {
					ok = queue.TryEnqueue(Task[*lifecyclev1alpha1.Machine]{Key: strconv.Itoa(key)})
					<-queue.Enqueued
					Expect(ok).To(BeTrue())
					wg.Done()
				}()
			}
			wg.Wait()
			Expect(queue.IsFull()).To(BeTrue())

			for range 20 {
				wg.Add(1)
				go func() {
					_, ok = queue.Dequeue()
					Expect(ok).To(BeTrue())
					wg.Done()
				}()
			}
			wg.Wait()
			Expect(queue.IsEmpty()).To(BeTrue())
		})
	})
})
