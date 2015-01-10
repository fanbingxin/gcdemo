package gc

import "container/list"

type Object interface {
        Next() Object
}

type Allocator struct {
        free  *list.List
        root  map[Object]bool
        all   map[Object]bool
        alloc func() Object
}

func NewAllocator(alloc func() Object) *Allocator {
        return &Allocator{
                free:  list.New(),
                root:  make(map[Object]bool),
                all:   make(map[Object]bool),
                alloc: alloc,
        }
}

func (l *Allocator) Malloc() Object {
        n := l.getFreeObject()
        if n != nil {
                return n
        }
        // no more free nodes, gc
        l.GC()

        // try again
        n = l.getFreeObject()
        if n != nil {
                return n
        }

        // still no free node
        n = l.alloc()
        l.all[n] = false
        return n
}

func (l *Allocator) Root(n Object) {
        _, ok := l.root[n]
        if ok {
                return
        }
        l.root[n] = true
}

func (l *Allocator) Unroot(n Object) {
        delete(l.root, n)
}

func (l *Allocator) GC() {
        // first clear all marks
        for k := range l.all {
                l.all[k] = false
        }

        // mark root objects recursively
        for n := range l.root {
                l.mark(n)
        }

        // sweep unmarked objects
        for k := range l.all {
                if !l.all[k] {
                        l.free.PushBack(k)
                }
        }
}

func (l *Allocator) getFreeObject() Object {
        elem := l.free.Back()
        if elem == nil {
                return nil
        }
        l.free.Remove(elem)
        return elem.Value.(Object)
}

func (l *Allocator) mark(n Object) {
        if n == nil {
                return
        }
        // we have marked this object
        if l.all[n] {
                return
        }
        l.all[n] = true
        l.mark(n.Next())
}
