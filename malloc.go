package gc

import "container/list"

type Node struct {
        Random *Node
        Next   *Node
}

type Allocator struct {
        free *list.List
        root map[*Node]bool
        all  map[*Node]bool
}

func NewAllocator() *Allocator {
        return &Allocator{
                free: list.New(),
                root: make(map[*Node]bool),
                all:  make(map[*Node]bool),
        }
}

func (l *Allocator) Malloc() *Node {
        n := l.getFreeNode()
        if n != nil {
                return n
        }
        // no more free nodes, gc
        l.GC()

        // try again
        n = l.getFreeNode()
        if n != nil {
                return n
        }

        // still no free node
        n = new(Node)
        l.all[n] = false
        return n
}

func (l *Allocator) Root(n *Node) {
        _, ok := l.root[n]
        if ok {
                return
        }
        l.root[n] = true
}

func (l *Allocator) Unroot(n *Node) {
        delete(l.root, n)
}

func (l *Allocator) GC() {
        // first clear all marks
        for k := range l.all {
                l.all[k] = false
        }

        // mark root nodes recursively
        for n := range l.root {
                l.mark(n)
        }

        // sweep unmarked nodes
        for k := range l.all {
                if !l.all[k] {
                        l.free.PushBack(k)
                }
        }
}

func (l *Allocator) getFreeNode() *Node {
        elem := l.free.Back()
        if elem == nil {
                return nil
        }
        l.free.Remove(elem)
        return elem.Value.(*Node)
}

func (l *Allocator) mark(n *Node) {
        if n == nil {
                return
        }
        // we have marked this node
        if l.all[n] {
                return
        }
        l.all[n] = true
        l.mark(n.Next)
}
