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

func assertNamePairSetImplementation() {
	var _ container.Set = (*NamePairSet)(nil)
}

// Set holds elements in a red-black Tree
type NamePairSet struct {
	*rbt.Tree
}

var itemExistsNamePairSet = struct{}{}

// NewWith instantiates a new empty set with the custom comparator.

func NewNamePairSet(Value ...common.NamePair) *NamePairSet {
	set := &NamePairSet{Tree: rbt.NewWith(common.CompareNamePair, false)}
	set.Add(Value...)
	return set
}

func CopyFromNamePairSet(ts *NamePairSet) *NamePairSet {
	return &NamePairSet{Tree: rbt.CopyFrom(ts.Tree)}
}

func NamePairSetIntersection(a *NamePairSet, b *NamePairSet, callback func(elem common.NamePair)) {
	aIterator := a.Iterator()
	bIterator := b.Iterator()

	if !aIterator.First() || !bIterator.First() {
		return
	}

	for aHasNext, bHasNext := true, true; aHasNext && bHasNext; {
		comp := common.CompareNamePair(aIterator.Value(), bIterator.Value())
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
func (set *NamePairSet) AddItem(item common.NamePair) (bool, common.NamePair) {
	itr := set.Tree.Insert(item, itemExistsNamePairSet)
	if itr.IsEnd() {
		return false, item
	}
	return true, itr.Key().(common.NamePair)
}

// Add adds the items (one or more) to the set.
func (set *NamePairSet) Add(items ...common.NamePair) {
	for _, item := range items {
		set.Tree.Put(item, itemExistsNamePairSet)
	}
}

// Remove removes the items (one or more) from the set.
func (set *NamePairSet) Remove(items ...common.NamePair) {
	for _, item := range items {
		set.Tree.Remove(item)
	}

}

// Values returns all items in the set.
func (set *NamePairSet) Values() []common.NamePair {
	keys := make([]common.NamePair, set.Size())
	it := set.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Value()
	}
	return keys
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *NamePairSet) Contains(items ...common.NamePair) bool {
	for _, item := range items {
		if iter := set.Get(item); iter.IsEnd() {
			return false
		}
	}
	return true
}

// String returns a string representation of container
func (set *NamePairSet) String() string {
	str := "TreeSet\n"
	items := make([]string, 0)
	for _, v := range set.Tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
type IteratorNamePairSet struct {
	rbt.Iterator
}

// Iterator holding the iterator's state
func (set *NamePairSet) Iterator() IteratorNamePairSet {
	return IteratorNamePairSet{Iterator: set.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (set *NamePairSet) Begin() IteratorNamePairSet {
	return IteratorNamePairSet{set.Tree.Begin()}
}

// End returns End Iterator
func (set *NamePairSet) End() IteratorNamePairSet {
	return IteratorNamePairSet{set.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorNamePairSet) Value() common.NamePair {
	return iterator.Iterator.Key().(common.NamePair)
}

// Each calls the given function once for each element, passing that element's index and value.
func (set *NamePairSet) Each(f func(value common.NamePair)) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (set *NamePairSet) Find(f func(value common.NamePair) bool) (v common.NamePair) {
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Value()) {
			return iterator.Value()
		}
	}
	return
}

func (set *NamePairSet) LowerBound(item common.NamePair) IteratorNamePairSet {
	return IteratorNamePairSet{set.Tree.LowerBound(item)}
}

func (set *NamePairSet) UpperBound(item common.NamePair) IteratorNamePairSet {
	return IteratorNamePairSet{set.Tree.UpperBound(item)}
}

// ToJSON outputs the JSON representation of the set.
func (set NamePairSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// FromJSON populates the set from the input JSON representation.
func (set *NamePairSet) UnmarshalJSON(data []byte) error {
	elements := make([]common.NamePair, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Tree = rbt.NewWith(common.CompareNamePair, false)
		set.Add(elements...)
	}
	return err
}

func (set NamePairSet) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(set.Size())...)
	set.Each(func(value common.NamePair) {
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (set *NamePairSet) Unpack(in []byte) (int, error) {
	set.Tree = rbt.NewWith(common.CompareNamePair, false)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		v := new(common.NamePair)
		decoder.Decode(v)
		set.Add(*v)
	}
	return decoder.GetPos(), nil
}
