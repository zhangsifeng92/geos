// Code generated by gotemplate. DO NOT EDIT.

package node_transaction

import (
	"fmt"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/libraries/container"
	"github.com/eosspark/eos-go/libraries/multiindex"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply)

// OrderedIndex holds elements of the red-black tree
type ById struct {
	super *ByExpiry             // index on the OrderedIndex, IndexBase is the last super index
	final *NodeTransactionIndex // index under the OrderedIndex, MultiIndex is the final index

	Root *ByIdNode
	size int
}

func (tree *ById) init(final *NodeTransactionIndex) {
	tree.final = final
	tree.super = &ByExpiry{}
	tree.super.init(final)
}

func (tree *ById) clear() {
	tree.Clear()
	tree.super.clear()
}

/*generic class*/

/*generic class*/

// OrderedIndexNode is a single element within the tree
type ByIdNode struct {
	Key    common.TransactionIdType
	super  *ByExpiryNode
	final  *NodeTransactionIndexNode
	color  colorById
	Left   *ByIdNode
	Right  *ByIdNode
	Parent *ByIdNode
}

/*generic class*/

/*generic class*/

func (node *ByIdNode) value() *multi_index.NodeTransactionState {
	return node.super.value()
}

type colorById bool

const (
	blackById, redById colorById = true, false
)

func (tree *ById) Insert(v multi_index.NodeTransactionState) (IteratorById, bool) {
	fn, res := tree.final.insert(v)
	if res {
		return tree.makeIterator(fn), true
	}
	return tree.End(), false
}

func (tree *ById) insert(v multi_index.NodeTransactionState, fn *NodeTransactionIndexNode) (*ByIdNode, bool) {
	key := ByIdFunc(v)

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

func (tree *ById) Erase(iter IteratorById) (itr IteratorById) {
	itr = iter
	itr.Next()
	tree.final.erase(iter.node.final)
	return
}

func (tree *ById) Erases(first, last IteratorById) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *ById) erase(n *ByIdNode) {
	tree.remove(n)
	tree.super.erase(n.super)
	n.super = nil
	n.final = nil
}

func (tree *ById) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(IteratorById); ok {
		tree.Erase(itr)
	} else {
		tree.super.erase_(iter)
	}
}

func (tree *ById) Modify(iter IteratorById, mod func(*multi_index.NodeTransactionState)) bool {
	if _, b := tree.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (tree *ById) modify(n *ByIdNode) (*ByIdNode, bool) {
	n.Key = ByIdFunc(*n.value())

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

func (tree *ById) modify_(iter multiindex.IteratorType, mod func(*multi_index.NodeTransactionState)) bool {
	if itr, ok := iter.(IteratorById); ok {
		return tree.Modify(itr, mod)
	} else {
		return tree.super.modify_(iter, mod)
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *ById) Find(key common.TransactionIdType) IteratorById {
	if false {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && ByIdCompare(key, lower.Key()) == 0 {
			return lower
		}
		return tree.End()
	} else {
		if node := tree.lookup(key); node != nil {
			return IteratorById{tree, node, betweenById}
		}
		return tree.End()
	}
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (tree *ById) LowerBound(key common.TransactionIdType) IteratorById {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByIdCompare(key, node.Key) > 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenById
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
func (tree *ById) UpperBound(key common.TransactionIdType) IteratorById {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if ByIdCompare(key, node.Key) >= 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = betweenById
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
func (tree *ById) Remove(key common.TransactionIdType) {
	if false {
		for lower := tree.LowerBound(key); lower.position != endById; {
			if ByIdCompare(lower.Key(), key) == 0 {
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

func (tree *ById) put(key common.TransactionIdType) (*ByIdNode, bool) {
	var insertedNode *ByIdNode
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		ByIdCompare(key, key)
		tree.Root = &ByIdNode{Key: key, color: redById}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		if false {
			for loop {
				compare := ByIdCompare(key, node.Key)
				switch {
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByIdNode{Key: key, color: redById}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare >= 0:
					if node.Right == nil {
						node.Right = &ByIdNode{Key: key, color: redById}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := ByIdCompare(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					return node, false
				case compare < 0:
					if node.Left == nil {
						node.Left = &ByIdNode{Key: key, color: redById}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare > 0:
					if node.Right == nil {
						node.Right = &ByIdNode{Key: key, color: redById}
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

func (tree *ById) swapNode(node *ByIdNode, pred *ByIdNode) {
	if node == pred {
		return
	}

	tmp := ByIdNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

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

func (tree *ById) remove(node *ByIdNode) {
	var child *ByIdNode
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
		if node.color == blackById {
			node.color = nodeColorById(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = blackById
		}
	}
	tree.size--
}

func (tree *ById) lookup(key common.TransactionIdType) *ByIdNode {
	node := tree.Root
	for node != nil {
		compare := ByIdCompare(key, node.Key)
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
func (tree *ById) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *ById) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *ById) Keys() []common.TransactionIdType {
	keys := make([]common.TransactionIdType, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *ById) Values() []multi_index.NodeTransactionState {
	values := make([]multi_index.NodeTransactionState, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *ById) Left() *ByIdNode {
	var parent *ByIdNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *ById) Right() *ByIdNode {
	var parent *ByIdNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *ById) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *ById) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		outputById(tree.Root, "", true, &str)
	}
	return str
}

func (node *ByIdNode) String() string {
	if !node.color {
		return fmt.Sprintf("(%v,%v)", node.Key, "red")
	}
	return fmt.Sprintf("(%v)", node.Key)
}

func outputById(node *ByIdNode, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		outputById(node.Right, newPrefix, false, str)
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
		outputById(node.Left, newPrefix, true, str)
	}
}

func (node *ByIdNode) grandparent() *ByIdNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *ByIdNode) uncle() *ByIdNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *ByIdNode) sibling() *ByIdNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (node *ByIdNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *ById) rotateLeft(node *ByIdNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *ById) rotateRight(node *ByIdNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *ById) replaceNode(old *ByIdNode, new *ByIdNode) {
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

func (tree *ById) insertCase1(node *ByIdNode) {
	if node.Parent == nil {
		node.color = blackById
	} else {
		tree.insertCase2(node)
	}
}

func (tree *ById) insertCase2(node *ByIdNode) {
	if nodeColorById(node.Parent) == blackById {
		return
	}
	tree.insertCase3(node)
}

func (tree *ById) insertCase3(node *ByIdNode) {
	uncle := node.uncle()
	if nodeColorById(uncle) == redById {
		node.Parent.color = blackById
		uncle.color = blackById
		node.grandparent().color = redById
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *ById) insertCase4(node *ByIdNode) {
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

func (tree *ById) insertCase5(node *ByIdNode) {
	node.Parent.color = blackById
	grandparent := node.grandparent()
	grandparent.color = redById
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *ByIdNode) maximumNode() *ByIdNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *ById) deleteCase1(node *ByIdNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *ById) deleteCase2(node *ByIdNode) {
	sibling := node.sibling()
	if nodeColorById(sibling) == redById {
		node.Parent.color = redById
		sibling.color = blackById
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *ById) deleteCase3(node *ByIdNode) {
	sibling := node.sibling()
	if nodeColorById(node.Parent) == blackById &&
		nodeColorById(sibling) == blackById &&
		nodeColorById(sibling.Left) == blackById &&
		nodeColorById(sibling.Right) == blackById {
		sibling.color = redById
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *ById) deleteCase4(node *ByIdNode) {
	sibling := node.sibling()
	if nodeColorById(node.Parent) == redById &&
		nodeColorById(sibling) == blackById &&
		nodeColorById(sibling.Left) == blackById &&
		nodeColorById(sibling.Right) == blackById {
		sibling.color = redById
		node.Parent.color = blackById
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *ById) deleteCase5(node *ByIdNode) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColorById(sibling) == blackById &&
		nodeColorById(sibling.Left) == redById &&
		nodeColorById(sibling.Right) == blackById {
		sibling.color = redById
		sibling.Left.color = blackById
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColorById(sibling) == blackById &&
		nodeColorById(sibling.Right) == redById &&
		nodeColorById(sibling.Left) == blackById {
		sibling.color = redById
		sibling.Right.color = blackById
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *ById) deleteCase6(node *ByIdNode) {
	sibling := node.sibling()
	sibling.color = nodeColorById(node.Parent)
	node.Parent.color = blackById
	if node == node.Parent.Left && nodeColorById(sibling.Right) == redById {
		sibling.Right.color = blackById
		tree.rotateLeft(node.Parent)
	} else if nodeColorById(sibling.Left) == redById {
		sibling.Left.color = blackById
		tree.rotateRight(node.Parent)
	}
}

func nodeColorById(node *ByIdNode) colorById {
	if node == nil {
		return blackById
	}
	return node.color
}

//////////////iterator////////////////

func (tree *ById) makeIterator(fn *NodeTransactionIndexNode) IteratorById {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*ByIdNode); ok {
			return IteratorById{tree: tree, node: n, position: betweenById}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

// Iterator holding the iterator's state
type IteratorById struct {
	tree     *ById
	node     *ByIdNode
	position positionById
}

type positionById byte

const (
	beginById, betweenById, endById positionById = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *ById) Iterator() IteratorById {
	return IteratorById{tree: tree, node: nil, position: beginById}
}

func (tree *ById) Begin() IteratorById {
	itr := IteratorById{tree: tree, node: nil, position: beginById}
	itr.Next()
	return itr
}

func (tree *ById) End() IteratorById {
	return IteratorById{tree: tree, node: nil, position: endById}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *IteratorById) Next() bool {
	if iterator.position == endById {
		goto end
	}
	if iterator.position == beginById {
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
	iterator.position = endById
	return false

between:
	iterator.position = betweenById
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *IteratorById) Prev() bool {
	if iterator.position == beginById {
		goto begin
	}
	if iterator.position == endById {
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
	iterator.position = beginById
	return false

between:
	iterator.position = betweenById
	return true
}

func (iterator IteratorById) HasNext() bool {
	return iterator.position != endById
}

func (iterator *IteratorById) HasPrev() bool {
	return iterator.position != beginById
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorById) Value() multi_index.NodeTransactionState {
	return *iterator.node.value()
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator IteratorById) Key() common.TransactionIdType {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *IteratorById) Begin() {
	iterator.node = nil
	iterator.position = beginById
}

func (iterator IteratorById) IsBegin() bool {
	return iterator.position == beginById
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *IteratorById) End() {
	iterator.node = nil
	iterator.position = endById
}

func (iterator IteratorById) IsEnd() bool {
	return iterator.position == endById
}

// Delete remove the node which pointed by the iterator
// Modifies the state of the iterator.
func (iterator *IteratorById) Delete() {
	node := iterator.node
	//iterator.Prev()
	iterator.tree.remove(node)
}

func (tree *ById) inPlace(n *ByIdNode) bool {
	prev := IteratorById{tree, n, betweenById}
	next := IteratorById{tree, n, betweenById}
	prev.Prev()
	next.Next()

	var (
		prevResult int
		nextResult int
	)

	if prev.IsBegin() {
		prevResult = 1
	} else {
		prevResult = ByIdCompare(n.Key, prev.Key())
	}

	if next.IsEnd() {
		nextResult = -1
	} else {
		nextResult = ByIdCompare(n.Key, next.Key())
	}

	return (false && prevResult >= 0 && nextResult <= 0) ||
		(!false && prevResult > 0 && nextResult < 0)
}
