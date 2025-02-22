// Code generated by gotemplate. DO NOT EDIT.

// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treemap implements a map backed by red-black Tree.
//
// Elements are ordered by key in the map.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
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

// template type Map(K,V,Compare,Multi)

func assertAccountNameUint64MapImplementation() {
	var _ container.Map = (*AccountNameUint64Map)(nil)
}

// Map holds the elements in a red-black Tree
type AccountNameUint64Map struct {
	*rbt.Tree
}

// NewWith instantiates a Tree map with the custom comparator.
func NewAccountNameUint64Map() *AccountNameUint64Map {
	return &AccountNameUint64Map{Tree: rbt.NewWith(common.CompareName, false)}
}

func CopyFromAccountNameUint64Map(tm *AccountNameUint64Map) *AccountNameUint64Map {
	return &AccountNameUint64Map{Tree: rbt.CopyFrom(tm.Tree)}
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *AccountNameUint64Map) Put(key common.AccountName, value uint64) {
	m.Tree.Put(key, value)
}

func (m *AccountNameUint64Map) Insert(key common.AccountName, value uint64) IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{m.Tree.Insert(key, value)}
}

// Get searches the element in the map by key and returns its value or nil if key is not found in Tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *AccountNameUint64Map) Get(key common.AccountName) IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{m.Tree.Get(key)}
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *AccountNameUint64Map) Remove(key common.AccountName) {
	m.Tree.Remove(key)
}

// Keys returns all keys in-order
func (m *AccountNameUint64Map) Keys() []common.AccountName {
	keys := make([]common.AccountName, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key().(common.AccountName)
	}
	return keys
}

// Values returns all values in-order based on the key.
func (m *AccountNameUint64Map) Values() []uint64 {
	values := make([]uint64, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value().(uint64)
	}
	return values
}

// Each calls the given function once for each element, passing that element's key and value.
func (m *AccountNameUint64Map) Each(f func(key common.AccountName, value uint64)) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		f(Iterator.Key(), Iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (key,value) for which the function is true or nil,nil otherwise if no element
// matches the criteria.
func (m *AccountNameUint64Map) Find(f func(key common.AccountName, value uint64) bool) (k common.AccountName, v uint64) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		if f(Iterator.Key(), Iterator.Value()) {
			return Iterator.Key(), Iterator.Value()
		}
	}
	return
}

// String returns a string representation of container
func (m AccountNameUint64Map) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}

// Iterator holding the Iterator's state
type IteratorAccountNameUint64Map struct {
	rbt.Iterator
}

// Iterator returns a stateful Iterator whose elements are key/value pairs.
func (m *AccountNameUint64Map) Iterator() IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{Iterator: m.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (m *AccountNameUint64Map) Begin() IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{m.Tree.Begin()}
}

// End returns End Iterator
func (m *AccountNameUint64Map) End() IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{m.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the Iterator.
func (iterator IteratorAccountNameUint64Map) Value() uint64 {
	return iterator.Iterator.Value().(uint64)
}

// Key returns the current element's key.
// Does not modify the state of the Iterator.
func (iterator IteratorAccountNameUint64Map) Key() common.AccountName {
	return iterator.Iterator.Key().(common.AccountName)
}

func (m *AccountNameUint64Map) LowerBound(key common.AccountName) IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{m.Tree.LowerBound(key)}
}

func (m *AccountNameUint64Map) UpperBound(key common.AccountName) IteratorAccountNameUint64Map {
	return IteratorAccountNameUint64Map{m.Tree.UpperBound(key)}

}

// ToJSON outputs the JSON representation of the map.
type pairAccountNameUint64Map struct {
	Key common.AccountName `json:"key"`
	Val uint64             `json:"val"`
}

func (m AccountNameUint64Map) MarshalJSON() ([]byte, error) {
	elements := make([]pairAccountNameUint64Map, 0, m.Size())
	it := m.Iterator()
	for it.Next() {
		elements = append(elements, pairAccountNameUint64Map{it.Key(), it.Value()})
	}
	return json.Marshal(&elements)
}

// FromJSON populates the map from the input JSON representation.
func (m *AccountNameUint64Map) UnmarshalJSON(data []byte) error {
	elements := make([]pairAccountNameUint64Map, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		m.Tree = rbt.NewWith(common.CompareName, false)
		for _, pair := range elements {
			m.Put(pair.Key, pair.Val)
		}
	}
	return err
}

func (m AccountNameUint64Map) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(m.Size())...)
	m.Each(func(key common.AccountName, value uint64) {
		rekey, _ := rlp.EncodeToBytes(key)
		re = append(re, rekey...)
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (m *AccountNameUint64Map) Unpack(in []byte) (int, error) {
	m.Tree = rbt.NewWith(common.CompareName, false)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		k, v := new(common.AccountName), new(uint64)
		decoder.Decode(k)
		decoder.Decode(v)
		m.Put(*k, *v)
	}
	return decoder.GetPos(), nil
}
