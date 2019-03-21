// Code generated by gotemplate. DO NOT EDIT.

package node_transaction

import (
	"fmt"

	"github.com/eosspark/eos-go/libraries/container"
	"github.com/eosspark/eos-go/libraries/multiindex"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply)

// OrderedIndex holds elements of the red-black tree
type ByBlockNum struct {
	super *NodeTransactionIndexBase // index on the OrderedIndex, IndexBase is the last super index
	final *NodeTransactionIndex     // index under the OrderedIndex, MultiIndex is the final index

	Root *ByBlockNumNode
	size int
}

func (tree *ByBlockNum) init(final *NodeTransactionIndex) {
	tree.final = final
	tree.super = &NodeTransactionIndexBase{}
	tree.super.init(final)
}

func (tree *ByBlockNum) clear() {
	tree.Clear()
	tree.super.clear()
}

/*generic class*/

/*generic class*/

// OrderedIndexNode is a single element within the tree
type ByBlockNumNode struct {
	Key    uint32
	super  *NodeTransactionIndexBaseNode
	final  *NodeTransactionIndexNode
	color  colorByBlockNum
	Left   *ByBlockNumNode
	Right  *ByBlockNumNode
	Parent *ByBlockNumNode
}

/*generic class*/

/*generic class*/

func (node *ByBlockNumNode) value() *multi_index.NodeTransactionState {
	return node.super.value()
}

type colorByBlockNum bool

const (
	blackByBlockNum, redByBlockNum colorByBlockNum = true, false
)

func (tree *ByBlockNum) Insert(v multi_index.NodeTransactionState) (IteratorByBlockNum, bool) {
	fn, res := tree.final.insert(v)
	if res {
		return tree.makeIterator(fn), true
	}
	return tree.End(), false
}

func (tree *ByBlockNum) insert(v multi_index.NodeTransactionState, fn *NodeTransactionIndexNode) (*ByBlockNumNode, bool) {
	key := ByBlockNumFunc(v)

	node, res := tree.put(key)
	if !res {
		container.Logger.Warn("#ordered index insert failed")
		return nil, false
	}
	sn, res := tree.super.insert(v, fn)
	if res {
		node.super = sn
		node.final = fn
		return node, true
	}
	tree.remove(node)
	return nil, false
}

func (tree *ByBlockNum) Erase(iter IteratorByBlockNum) (itr IteratorByBlockNum) {
	itr = iter
	itr.Next()
	tree.final.erase(iter.node.final)
	return
}

func (tree *ByBlockNum) Erases(first, last IteratorByBlockNum) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *ByBlockNum) erase(n *ByBlockNumNode) {
	tree.remove(n)
	tree.super.erase(n.super)
	n.super = nil
	n.final = nil
}

func (tree *ByBlockNum) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(IteratorByBlockNum); ok {
		tree.Erase(itr)
	} else {
		tree.super.erase_(iter)
	}
}

func (tree *ByBlockNum) Modify(iter IteratorByBlockNum, mod func(*multi_index.NodeTransactionState)) bool {
	if _, b := tree.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (tree *ByBlockNum) modify(n *ByBlockNumNode) (*ByBlockNumNode, bool) {
	n.Key = ByBlockNumFunc(*n.value())

	if !tree.inPlace(n) {
		tree.remove(n)
		node, res := tree.put(n.Key)
		if !res {
			container.Logger.Warn("#ordered index modify failed")
			tree.super.erase(n.super)
			return nil, false
		}

		//n.Left = node.Left
		//if n.Left != nil {
		//	n.Left.Parent = n
		//}
		//n.Right = node.Right
		//if n.Right != nil {
		//	n.Right.Parent = n
		//}
		//n.Parent = node.Parent
		//if n.Parent != nil {
		//	if n.Parent.Left == node {
		//		n.Parent.Left = n
		//	} else {
		//		n.Parent.Right = n
		//	}
		//} else {
		//	tree.Root = n
		//}
		node.super = n.super
		node.final = n.final
		n = node
	}

	if sn, res := tree.super.modify(n.super); !res {
		tree.remove(n)
		return nil, false
	} else {
		n.super = sn
	}

	return n, true
}

func (tree *ByBlockNum) modify_(iter multiindex.IteratorType, mod func(*multi_index.NodeTransactionState)) bool {
	if itr, ok := iter.(IteratorByBlockNum); ok {
		return tree.Modify(itr, mod)
	} else {
		return tree.super.modify_(iter, mod)
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByBlockNum) Find(key uint32) IteratorByBlockNum {
	if true {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && ByBlockNumCompare(key, lower.Key()) == 0 {
			return lower
		}
		return tree.End()
	} else {
		if node := tree.lookup(key); node != nil {
			return IteratorByBlockNum{tree, node, betweenByBlockNum}
		}
		return tree.End()
	}
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (tree *ByBlockNum) LowerBound(key uint32) IteratorByBlockNum {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByBlockNumCompare(key, node.Key) > 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenByBlockNum
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (tree *ByBlockNum) UpperBound(key uint32) IteratorByBlockNum {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByBlockNumCompare(key, node.Key) >= 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenByBlockNum
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ByBlockNum) Remove(key uint32) {
	if true {
		for lower := tree.LowerBound(key); lower.position != endByBlockNum; {
			if ByBlockNumCompare(lower.Key(), key) == 0 {
				node := lower.node
				lower.Next()
				tree.remove(node)
			} else {
				break
			}
		}
	} else {
		node := tree.lookup(key)
		tree.remove(node)
	}
}

func (tree *ByBlockNum) put(key uint32) (*ByBlockNumNode, bool) {
	var insertedNode *ByBlockNumNode
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		ByBlockNumCompare(key, key)
		tree.Root = &ByBlockNumNode{Key: key, color: redByBlockNum}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		if true {
			for loop {
				compare := ByBlockNumCompare(key, node.Key)
				switch {
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByBlockNumNode{Key: key, color: redByBlockNum}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare >= 0:
					if node.Right == nil {
						node.Right = &ByBlockNumNode{Key: key, color: redByBlockNum}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := ByBlockNumCompare(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					return node, false
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByBlockNumNode{Key: key, color: redByBlockNum}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare > 0:
					if node.Right == nil {
						node.Right = &ByBlockNumNode{Key: key, color: redByBlockNum}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++

	return insertedNode, true
}

func (tree *ByBlockNum) swapNode(node *ByBlockNumNode, pred *ByBlockNumNode) {
	if node == pred {
		return
	}

	tmp := ByBlockNumNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

	pred.color = node.color
	node.color = tmp.color

	pred.Right = node.Right
	if pred.Right != nil {
		pred.Right.Parent = pred
	}
	node.Right = tmp.Right
	if node.Right != nil {
		node.Right.Parent = node
	}

	if pred.Parent == node {
		pred.Left = node
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}
		node.Parent = pred

	} else {
		pred.Left = node.Left
		if pred.Left != nil {
			pred.Left.Parent = pred
		}
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}

		node.Parent = tmp.Parent
		if node.Parent != nil {
			if node.Parent.Left == pred {
				node.Parent.Left = node
			} else {
				node.Parent.Right = node
			}
		} else {
			tree.Root = node
		}
	}
}

func (tree *ByBlockNum) remove(node *ByBlockNumNode) {
	var child *ByBlockNumNode
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		tree.swapNode(node, pred)
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == blackByBlockNum {
			node.color = nodeColorByBlockNum(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = blackByBlockNum
		}
	}
	tree.size--
}

func (tree *ByBlockNum) lookup(key uint32) *ByBlockNumNode {
	node := tree.Root
	for node != nil {
		compare := ByBlockNumCompare(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

// Empty returns true if tree does not contain any nodes
func (tree *ByBlockNum) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *ByBlockNum) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *ByBlockNum) Keys() []uint32 {
	keys := make([]uint32, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *ByBlockNum) Values() []multi_index.NodeTransactionState {
	values := make([]multi_index.NodeTransactionState, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *ByBlockNum) Left() *ByBlockNumNode {
	var parent *ByBlockNumNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *ByBlockNum) Right() *ByBlockNumNode {
	var parent *ByBlockNumNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *ByBlockNum) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *ByBlockNum) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		outputByBlockNum(tree.Root, "", true, &str)
	}
	return str
}

func (node *ByBlockNumNode) String() string {
	if !node.color {
		return fmt.Sprintf("(%v,%v)", node.Key, "red")
	}
	return fmt.Sprintf("(%v)", node.Key)
}

func outputByBlockNum(node *ByBlockNumNode, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		outputByBlockNum(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		outputByBlockNum(node.Left, newPrefix, true, str)
	}
}

func (node *ByBlockNumNode) grandparent() *ByBlockNumNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *ByBlockNumNode) uncle() *ByBlockNumNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *ByBlockNumNode) sibling() *ByBlockNumNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (node *ByBlockNumNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *ByBlockNum) rotateLeft(node *ByBlockNumNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *ByBlockNum) rotateRight(node *ByBlockNumNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *ByBlockNum) replaceNode(old *ByBlockNumNode, new *ByBlockNumNode) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *ByBlockNum) insertCase1(node *ByBlockNumNode) {
	if node.Parent == nil {
		node.color = blackByBlockNum
	} else {
		tree.insertCase2(node)
	}
}

func (tree *ByBlockNum) insertCase2(node *ByBlockNumNode) {
	if nodeColorByBlockNum(node.Parent) == blackByBlockNum {
		return
	}
	tree.insertCase3(node)
}

func (tree *ByBlockNum) insertCase3(node *ByBlockNumNode) {
	uncle := node.uncle()
	if nodeColorByBlockNum(uncle) == redByBlockNum {
		node.Parent.color = blackByBlockNum
		uncle.color = blackByBlockNum
		node.grandparent().color = redByBlockNum
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *ByBlockNum) insertCase4(node *ByBlockNumNode) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *ByBlockNum) insertCase5(node *ByBlockNumNode) {
	node.Parent.color = blackByBlockNum
	grandparent := node.grandparent()
	grandparent.color = redByBlockNum
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *ByBlockNumNode) maximumNode() *ByBlockNumNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *ByBlockNum) deleteCase1(node *ByBlockNumNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *ByBlockNum) deleteCase2(node *ByBlockNumNode) {
	sibling := node.sibling()
	if nodeColorByBlockNum(sibling) == redByBlockNum {
		node.Parent.color = redByBlockNum
		sibling.color = blackByBlockNum
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *ByBlockNum) deleteCase3(node *ByBlockNumNode) {
	sibling := node.sibling()
	if nodeColorByBlockNum(node.Parent) == blackByBlockNum &&
		nodeColorByBlockNum(sibling) == blackByBlockNum &&
		nodeColorByBlockNum(sibling.Left) == blackByBlockNum &&
		nodeColorByBlockNum(sibling.Right) == blackByBlockNum {
		sibling.color = redByBlockNum
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *ByBlockNum) deleteCase4(node *ByBlockNumNode) {
	sibling := node.sibling()
	if nodeColorByBlockNum(node.Parent) == redByBlockNum &&
		nodeColorByBlockNum(sibling) == blackByBlockNum &&
		nodeColorByBlockNum(sibling.Left) == blackByBlockNum &&
		nodeColorByBlockNum(sibling.Right) == blackByBlockNum {
		sibling.color = redByBlockNum
		node.Parent.color = blackByBlockNum
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *ByBlockNum) deleteCase5(node *ByBlockNumNode) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColorByBlockNum(sibling) == blackByBlockNum &&
		nodeColorByBlockNum(sibling.Left) == redByBlockNum &&
		nodeColorByBlockNum(sibling.Right) == blackByBlockNum {
		sibling.color = redByBlockNum
		sibling.Left.color = blackByBlockNum
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColorByBlockNum(sibling) == blackByBlockNum &&
		nodeColorByBlockNum(sibling.Right) == redByBlockNum &&
		nodeColorByBlockNum(sibling.Left) == blackByBlockNum {
		sibling.color = redByBlockNum
		sibling.Right.color = blackByBlockNum
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *ByBlockNum) deleteCase6(node *ByBlockNumNode) {
	sibling := node.sibling()
	sibling.color = nodeColorByBlockNum(node.Parent)
	node.Parent.color = blackByBlockNum
	if node == node.Parent.Left && nodeColorByBlockNum(sibling.Right) == redByBlockNum {
		sibling.Right.color = blackByBlockNum
		tree.rotateLeft(node.Parent)
	} else if nodeColorByBlockNum(sibling.Left) == redByBlockNum {
		sibling.Left.color = blackByBlockNum
		tree.rotateRight(node.Parent)
	}
}

func nodeColorByBlockNum(node *ByBlockNumNode) colorByBlockNum {
	if node == nil {
		return blackByBlockNum
	}
	return node.color
}

//////////////iterator////////////////

func (tree *ByBlockNum) makeIterator(fn *NodeTransactionIndexNode) IteratorByBlockNum {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*ByBlockNumNode); ok {
			return IteratorByBlockNum{tree: tree, node: n, position: betweenByBlockNum}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

// Iterator holding the iterator's state
type IteratorByBlockNum struct {
	tree     *ByBlockNum
	node     *ByBlockNumNode
	position positionByBlockNum
}

type positionByBlockNum byte

const (
	beginByBlockNum, betweenByBlockNum, endByBlockNum positionByBlockNum = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *ByBlockNum) Iterator() IteratorByBlockNum {
	return IteratorByBlockNum{tree: tree, node: nil, position: beginByBlockNum}
}

func (tree *ByBlockNum) Begin() IteratorByBlockNum {
	itr := IteratorByBlockNum{tree: tree, node: nil, position: beginByBlockNum}
	itr.Next()
	return itr
}

func (tree *ByBlockNum) End() IteratorByBlockNum {
	return IteratorByBlockNum{tree: tree, node: nil, position: endByBlockNum}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *IteratorByBlockNum) Next() bool {
	if iterator.position == endByBlockNum {
		goto end
	}
	if iterator.position == beginByBlockNum {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Left {
				goto between
			}
			node = iterator.node
		}
	}

end:
	iterator.node = nil
	iterator.position = endByBlockNum
	return false

between:
	iterator.position = betweenByBlockNum
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *IteratorByBlockNum) Prev() bool {
	if iterator.position == beginByBlockNum {
		goto begin
	}
	if iterator.position == endByBlockNum {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Right {
				goto between
			}
			node = iterator.node
			//if iterator.tree.Comparator(node.Key, iterator.node.Key) >= 0 {
			//	goto between
			//}
		}
	}

begin:
	iterator.node = nil
	iterator.position = beginByBlockNum
	return false

between:
	iterator.position = betweenByBlockNum
	return true
}

func (iterator IteratorByBlockNum) HasNext() bool {
	return iterator.position != endByBlockNum
}

func (iterator *IteratorByBlockNum) HasPrev() bool {
	return iterator.position != beginByBlockNum
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorByBlockNum) Value() multi_index.NodeTransactionState {
	return *iterator.node.value()
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator IteratorByBlockNum) Key() uint32 {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *IteratorByBlockNum) Begin() {
	iterator.node = nil
	iterator.position = beginByBlockNum
}

func (iterator IteratorByBlockNum) IsBegin() bool {
	return iterator.position == beginByBlockNum
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *IteratorByBlockNum) End() {
	iterator.node = nil
	iterator.position = endByBlockNum
}

func (iterator IteratorByBlockNum) IsEnd() bool {
	return iterator.position == endByBlockNum
}

// Delete remove the node which pointed by the iterator
// Modifies the state of the iterator.
func (iterator *IteratorByBlockNum) Delete() {
	node := iterator.node
	//iterator.Prev()
	iterator.tree.remove(node)
}

func (tree *ByBlockNum) inPlace(n *ByBlockNumNode) bool {
	prev := IteratorByBlockNum{tree, n, betweenByBlockNum}
	next := IteratorByBlockNum{tree, n, betweenByBlockNum}
	prev.Prev()
	next.Next()

	var (
		prevResult int
		nextResult int
	)

	if prev.IsBegin() {
		prevResult = 1
	} else {
		prevResult = ByBlockNumCompare(n.Key, prev.Key())
	}

	if next.IsEnd() {
		nextResult = -1
	} else {
		nextResult = ByBlockNumCompare(n.Key, next.Key())
	}

	return (true && prevResult >= 0 && nextResult <= 0) ||
		(!true && prevResult > 0 && nextResult < 0)
}
