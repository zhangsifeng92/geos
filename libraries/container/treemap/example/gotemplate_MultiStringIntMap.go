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
package example

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

func assertMultiStringIntMapImplementation() {
	var _ container.Map = (*MultiStringIntMap)(nil)
}

// Map holds the elements in a red-black Tree
type MultiStringIntMap struct {
	*rbt.Tree
}

// NewWith instantiates a Tree map with the custom comparator.
func NewMultiStringIntMap() *MultiStringIntMap {
	return &MultiStringIntMap{Tree: rbt.NewWith(StringComparator, true)}
}

func CopyFromMultiStringIntMap(tm *MultiStringIntMap) *MultiStringIntMap {
	return &MultiStringIntMap{Tree: rbt.CopyFrom(tm.Tree)}
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiStringIntMap) Put(key string, value int) {
	m.Tree.Put(key, value)
}

func (m *MultiStringIntMap) Insert(key string, value int) IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{m.Tree.Insert(key, value)}
}

// Get searches the element in the map by key and returns its value or nil if key is not found in Tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiStringIntMap) Get(key string) IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{m.Tree.Get(key)}
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiStringIntMap) Remove(key string) {
	m.Tree.Remove(key)
}

// Keys returns all keys in-order
func (m *MultiStringIntMap) Keys() []string {
	keys := make([]string, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key().(string)
	}
	return keys
}

// Values returns all values in-order based on the key.
func (m *MultiStringIntMap) Values() []int {
	values := make([]int, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value().(int)
	}
	return values
}

// Each calls the given function once for each element, passing that element's key and value.
func (m *MultiStringIntMap) Each(f func(key string, value int)) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		f(Iterator.Key(), Iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (key,value) for which the function is true or nil,nil otherwise if no element
// matches the criteria.
func (m *MultiStringIntMap) Find(f func(key string, value int) bool) (k string, v int) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		if f(Iterator.Key(), Iterator.Value()) {
			return Iterator.Key(), Iterator.Value()
		}
	}
	return
}

// String returns a string representation of container
func (m MultiStringIntMap) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}

// Iterator holding the Iterator's state
type IteratorMultiStringIntMap struct {
	rbt.Iterator
}

// Iterator returns a stateful Iterator whose elements are key/value pairs.
func (m *MultiStringIntMap) Iterator() IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{Iterator: m.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (m *MultiStringIntMap) Begin() IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{m.Tree.Begin()}
}

// End returns End Iterator
func (m *MultiStringIntMap) End() IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{m.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the Iterator.
func (iterator IteratorMultiStringIntMap) Value() int {
	return iterator.Iterator.Value().(int)
}

// Key returns the current element's key.
// Does not modify the state of the Iterator.
func (iterator IteratorMultiStringIntMap) Key() string {
	return iterator.Iterator.Key().(string)
}

func (m *MultiStringIntMap) LowerBound(key string) IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{m.Tree.LowerBound(key)}
}

func (m *MultiStringIntMap) UpperBound(key string) IteratorMultiStringIntMap {
	return IteratorMultiStringIntMap{m.Tree.UpperBound(key)}

}

// ToJSON outputs the JSON representation of the map.
type pairMultiStringIntMap struct {
	Key string `json:"key"`
	Val int    `json:"val"`
}

func (m MultiStringIntMap) MarshalJSON() ([]byte, error) {
	elements := make([]pairMultiStringIntMap, 0, m.Size())
	it := m.Iterator()
	for it.Next() {
		elements = append(elements, pairMultiStringIntMap{it.Key(), it.Value()})
	}
	return json.Marshal(&elements)
}

// FromJSON populates the map from the input JSON representation.
func (m *MultiStringIntMap) UnmarshalJSON(data []byte) error {
	elements := make([]pairMultiStringIntMap, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		m.Tree = rbt.NewWith(StringComparator, true)
		for _, pair := range elements {
			m.Put(pair.Key, pair.Val)
		}
	}
	return err
}

func (m MultiStringIntMap) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(m.Size())...)
	m.Each(func(key string, value int) {
		rekey, _ := rlp.EncodeToBytes(key)
		re = append(re, rekey...)
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (m *MultiStringIntMap) Unpack(in []byte) (int, error) {
	m.Tree = rbt.NewWith(StringComparator, true)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		k, v := new(string), new(int)
		decoder.Decode(k)
		decoder.Decode(v)
		m.Put(*k, *v)
	}
	return decoder.GetPos(), nil
}
