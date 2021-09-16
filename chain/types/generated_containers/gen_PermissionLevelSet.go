// Code generated by gotemplate. DO NOT EDIT.

// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treeset implements a Tree backed by a red-black Tree.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package generated

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zhangsifeng92/geos/common"
	"github.com/zhangsifeng92/geos/crypto/rlp"
	"github.com/zhangsifeng92/geos/libraries/container"
	rbt "github.com/zhangsifeng92/geos/libraries/container/redblacktree"
)

// template type Set(V,Compare,Multi)

func assertPermissionLevelSetImplementation() {
	var _ container.Set = (*PermissionLevelSet)(nil)
}

// Set holds elements in a red-black Tree
type PermissionLevelSet struct {
	*rbt.Tree
}

var itemExistsPermissionLevelSet = struct{}{}

// NewWith instantiates a new empty set with the custom comparator.

func NewPermissionLevelSet(Value ...common.PermissionLevel) *PermissionLevelSet {
	set := &PermissionLevelSet{Tree: rbt.NewWith(common.ComparePermissionLevel, false)}
	set.Add(Value...)
	return set
}

func CopyFromPermissionLevelSet(ts *PermissionLevelSet) *PermissionLevelSet {
	return &PermissionLevelSet{Tree: rbt.CopyFrom(ts.Tree)}
}

func PermissionLevelSetIntersection(a *PermissionLevelSet, b *PermissionLevelSet, callback func(elem common.PermissionLevel)) {
	aIterator := a.Iterator()
	bIterator := b.Iterator()

	if !aIterator.First() || !bIterator.First() {
		return
	}

	for aHasNext, bHasNext := true, true; aHasNext && bHasNext; {
		comp := common.ComparePermissionLevel(aIterator.Value(), bIterator.Value())
		switch {
		case comp > 0:
			bHasNext = bIterator.Next()
		case comp < 0:
			aHasNext = aIterator.Next()
		default:
			callback(aIterator.Value())
			aHasNext = aIterator.Next()
			bHasNext = bIterator.Next()
		}
	}
}

// Add adds the item one to the set.Returns false and the interface if it already exists
func (set *PermissionLevelSet) AddItem(item common.PermissionLevel) (bool, common.PermissionLevel) {
	itr := set.Tree.Insert(item, itemExistsPermissionLevelSet)
	if itr.IsEnd() {
		return false, item
	}
	return true, itr.Key().(common.PermissionLevel)
}

// Add adds the items (one or more) to the set.
func (set *PermissionLevelSet) Add(items ...common.PermissionLevel) {
	for _, item := range items {
		set.Tree.Put(item, itemExistsPermissionLevelSet)
	}
}

// Remove removes the items (one or more) from the set.
func (set *PermissionLevelSet) Remove(items ...common.PermissionLevel) {
	for _, item := range items {
		set.Tree.Remove(item)
	}

}

// Values returns all items in the set.
func (set *PermissionLevelSet) Values() []common.PermissionLevel {
	keys := make([]common.PermissionLevel, set.Size())
	it := set.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Value()
	}
	return keys
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *PermissionLevelSet) Contains(items ...common.PermissionLevel) bool {
	for _, item := range items {
		if iter := set.Get(item); iter.IsEnd() {
			return false
		}
	}
	return true
}

// String returns a string representation of container
func (set *PermissionLevelSet) String() string {
	str := "TreeSet\n"
	items := make([]string, 0)
	for _, v := range set.Tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
type IteratorPermissionLevelSet struct {
	rbt.Iterator
}

// Iterator holding the iterator's state
func (set *PermissionLevelSet) Iterator() IteratorPermissionLevelSet {
	return IteratorPermissionLevelSet{Iterator: set.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (set *PermissionLevelSet) Begin() IteratorPermissionLevelSet {
	return IteratorPermissionLevelSet{set.Tree.Begin()}
}

// End returns End Iterator
func (set *PermissionLevelSet) End() IteratorPermissionLevelSet {
	return IteratorPermissionLevelSet{set.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorPermissionLevelSet) Value() common.PermissionLevel {
	return iterator.Iterator.Key().(common.PermissionLevel)
}

// Each calls the given function once for each element, passing that element's index and value.
func (set *PermissionLevelSet) Each(f func(value common.PermissionLevel)) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (set *PermissionLevelSet) Find(f func(value common.PermissionLevel) bool) (v common.PermissionLevel) {
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Value()) {
			return iterator.Value()
		}
	}
	return
}

func (set *PermissionLevelSet) LowerBound(item common.PermissionLevel) IteratorPermissionLevelSet {
	return IteratorPermissionLevelSet{set.Tree.LowerBound(item)}
}

func (set *PermissionLevelSet) UpperBound(item common.PermissionLevel) IteratorPermissionLevelSet {
	return IteratorPermissionLevelSet{set.Tree.UpperBound(item)}
}

// ToJSON outputs the JSON representation of the set.
func (set PermissionLevelSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// FromJSON populates the set from the input JSON representation.
func (set *PermissionLevelSet) UnmarshalJSON(data []byte) error {
	elements := make([]common.PermissionLevel, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Tree = rbt.NewWith(common.ComparePermissionLevel, false)
		set.Add(elements...)
	}
	return err
}

func (set PermissionLevelSet) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(set.Size())...)
	set.Each(func(value common.PermissionLevel) {
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (set *PermissionLevelSet) Unpack(in []byte) (int, error) {
	set.Tree = rbt.NewWith(common.ComparePermissionLevel, false)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		v := new(common.PermissionLevel)
		decoder.Decode(v)
		set.Add(*v)
	}
	return decoder.GetPos(), nil
}
