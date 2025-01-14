// Code generated by gotemplate. DO NOT EDIT.

package peer_block_state

import (
	"github.com/zhangsifeng92/geos/libraries/container"
	"github.com/zhangsifeng92/geos/libraries/multiindex"
	"github.com/zhangsifeng92/geos/plugins/net_plugin/multi_index"
)

// template type MultiIndex(SuperIndex,SuperNode,Value)

type PeerBlockStateIndex struct {
	super *ById
	count int
}

func NewPeerBlockStateIndex() *PeerBlockStateIndex {
	m := &PeerBlockStateIndex{}
	m.super = &ById{}
	m.super.init(m)
	return m
}

/*generic class*/

type PeerBlockStateIndexNode struct {
	super *ByIdNode
}

/*generic class*/

//method for MultiIndex
func (m *PeerBlockStateIndex) GetSuperIndex() interface{} { return m.super }
func (m *PeerBlockStateIndex) GetFinalIndex() interface{} { return nil }

func (m *PeerBlockStateIndex) GetIndex() interface{} {
	return nil
}

func (m *PeerBlockStateIndex) Size() int {
	return m.count
}

func (m *PeerBlockStateIndex) Clear() {
	m.super.clear()
	m.count = 0
}

func (m *PeerBlockStateIndex) Insert(v multi_index.PeerBlockState) bool {
	_, res := m.insert(v)
	return res
}

func (m *PeerBlockStateIndex) insert(v multi_index.PeerBlockState) (*PeerBlockStateIndexNode, bool) {
	fn := &PeerBlockStateIndexNode{}
	n, res := m.super.insert(v, fn)
	if res {
		fn.super = n
		m.count++
		return fn, true
	}
	return nil, false
}

func (m *PeerBlockStateIndex) Erase(iter multiindex.IteratorType) {
	m.super.erase_(iter)
}

func (m *PeerBlockStateIndex) erase(n *PeerBlockStateIndexNode) {
	m.super.erase(n.super)
	m.count--
}

func (m *PeerBlockStateIndex) Modify(iter multiindex.IteratorType, mod func(*multi_index.PeerBlockState)) bool {
	return m.super.modify_(iter, mod)
}

func (m *PeerBlockStateIndex) modify(mod func(*multi_index.PeerBlockState), n *PeerBlockStateIndexNode) (*PeerBlockStateIndexNode, bool) {
	defer func() {
		if e := recover(); e != nil {
			container.Logger.Error("#multi modify failed: %v", e)
			m.erase(n)
			m.count--
			panic(e)
		}
	}()
	mod(n.value())
	if sn, res := m.super.modify(n.super); !res {
		m.count--
		return nil, false
	} else {
		n.super = sn
		return n, true
	}
}

func (n *PeerBlockStateIndexNode) GetSuperNode() interface{} { return n.super }
func (n *PeerBlockStateIndexNode) GetFinalNode() interface{} { return nil }

func (n *PeerBlockStateIndexNode) value() *multi_index.PeerBlockState {
	return n.super.value()
}

/// IndexBase
type PeerBlockStateIndexBase struct {
	final *PeerBlockStateIndex
}

type PeerBlockStateIndexBaseNode struct {
	final *PeerBlockStateIndexNode
	pv    *multi_index.PeerBlockState
}

func (i *PeerBlockStateIndexBase) init(final *PeerBlockStateIndex) {
	i.final = final
}

func (i *PeerBlockStateIndexBase) clear() {}

func (i *PeerBlockStateIndexBase) GetSuperIndex() interface{} { return nil }

func (i *PeerBlockStateIndexBase) GetFinalIndex() interface{} { return i.final }

func (i *PeerBlockStateIndexBase) insert(v multi_index.PeerBlockState, fn *PeerBlockStateIndexNode) (*PeerBlockStateIndexBaseNode, bool) {
	return &PeerBlockStateIndexBaseNode{fn, &v}, true
}

func (i *PeerBlockStateIndexBase) erase(n *PeerBlockStateIndexBaseNode) {
	n.pv = nil
}

func (i *PeerBlockStateIndexBase) erase_(iter multiindex.IteratorType) {
	container.Logger.Warn("erase iterator doesn't match all index")
}

func (i *PeerBlockStateIndexBase) modify(n *PeerBlockStateIndexBaseNode) (*PeerBlockStateIndexBaseNode, bool) {
	return n, true
}

func (i *PeerBlockStateIndexBase) modify_(iter multiindex.IteratorType, mod func(*multi_index.PeerBlockState)) bool {
	container.Logger.Warn("modify iterator doesn't match all index")
	return false
}

func (n *PeerBlockStateIndexBaseNode) value() *multi_index.PeerBlockState {
	return n.pv
}
