package gc

import "testing"

type Node struct {
        val  int
        next *Node
}

func (n *Node) Children() []Object {
        return []Object{n.next}
}

func NewNode() Object {
        return new(Node)
}

func TestMalloc(t *testing.T) {
        l := NewAllocator(NewNode)
        root := l.Malloc().(*Node)
        l.Root(root)
        stat := l.Status()
        if stat.Total != 1 || stat.Free != 0 {
                t.Fatal("status error:%v", stat)
        }
        n := root
        for i := 0; i < 5; i++ {
                n.next = l.Malloc().(*Node)
                n.val = i
                n = n.next
        }
        stat = l.Status()
        if stat.Total != 6 || stat.Free != 0 {
                t.Fatal("status error:%v", stat)
        }
        root.next = nil
        l.GC()
        stat = l.Status()
        if stat.Total != 6 || stat.Free != 5 {
                t.Fatal("status error:%v", stat)
        }
}

func TestCircleRef(t *testing.T) {
        l := NewAllocator(NewNode)
        root := l.Malloc().(*Node)
        l.Root(root)

        n0 := l.Malloc().(*Node)
        root.next = n0
        n1 := l.Malloc().(*Node)
        n0.next = n1
        n1.next = n0

        root.next = nil
        l.GC()
        stat := l.Status()
        if stat.Total != 3 || stat.Free != 2 {
                t.Fatal("status error:%v", stat)
        }
}
