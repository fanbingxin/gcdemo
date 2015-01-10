package gc

import "testing"

type Node struct {
        val  int
        next *Node
}

func (n *Node) Next() Object {
        if n.next != nil {
                return n.next
        }
        return nil
}

func NewNode() Object {
        return new(Node)
}

func TestMalloc(t *testing.T) {
        l := NewAllocator(NewNode)
        root := l.Malloc().(*Node)
        l.Root(root)
        stat := l.Status()
        if stat.Total != 1 && stat.Free != 0 {
                t.Fatal("status error:%v", stat)
        }
        n := root
        for i := 0; i < 5; i++ {
                n.next = l.Malloc().(*Node)
                n.val = i
                n = n.next
        }
        if stat.Total != 1 && stat.Free != 0 {
                t.Fatal("status error:%v", stat)
        }
        root.next = nil
        l.Malloc()
}
