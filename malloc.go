package gc

import (
        "container/list"
        "log"
)

type Object interface {
        Children() []Object
}

type Stat struct {
        Free  int
        Total int
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

func (l *Allocator) Status() *Stat {
        return &Stat{
                Free:  l.free.Len(),
                Total: len(l.all),
        }
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

        free := 0
        // sweep unmarked objects
        for k := range l.all {
                if !l.all[k] {
                        free++
                        l.free.PushBack(k)
                }
        }
        log.Printf("GC, free %d objects, current %d objects", free, len(l.all))
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
        marked, ok := l.all[n]
        if !ok || marked {
                return
        }
        l.all[n] = true
        for _, child := range n.Children() {
                l.mark(child)
        }
}
