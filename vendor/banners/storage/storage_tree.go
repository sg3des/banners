package storage

import (
	"banners/banners"
	"math/rand"
	"sync/atomic"
)

type branchCallback func(i leaf) bool

type leaf *banners.Banner

type branch struct {
	branches [256]*branch
	leafs    []leaf
}

func (b *branch) add(l leaf, key []byte) {
	b.leafs = append(b.leafs, l)

	if len(key) > 0 {
		if b.branches[key[0]] == nil {
			b.branches[key[0]] = new(branch)
		}
		b.branches[key[0]].add(l, key[1:])
	}
}

func (b *branch) lookupAny(call branchCallback) (leaf, bool) {
	iter := NewRandIter(len(b.leafs))
	for iter.Next() {
		l := b.leafs[iter.i]
		if call(l) {
			return l, true
		}
	}

	return nil, false
}

func (b *branch) lookup(key []byte, call branchCallback) (leaf, bool) {
	if len(key) == 0 {
		return b.lookupAnyLeaf(call)
	}

	return b.branches[key[0]].lookup(key[1:], call)
}

func (b *branch) lookupAnyLeaf(call branchCallback) (leaf, bool) {
	if len(b.leafs) > 0 {
		iter := NewRandIter(len(b.leafs))
		for iter.Next() {
			if call(b.leafs[iter.i]) {
				return b.leafs[iter.i], true
			}
		}
	}

	return nil, false
}

type RandIter struct {
	limit int
	x     int
	i     int
}

func NewRandIter(size int) *RandIter {
	iter := &RandIter{limit: size - 1}
	iter.x = rand.Intn(size)
	iter.i = iter.x - 1

	return iter
}

func (iter *RandIter) Next() bool {
	if iter.i < iter.x {
		iter.i--
	} else {
		iter.i++
	}

	if iter.i < 0 {
		iter.i = iter.x + 1
	}

	if iter.i > iter.limit {
		return false
	}

	return true
}

//
//
//

type TreeStorage struct {
	root *branch
}

func NewTreeStorage() *TreeStorage {
	ts := new(TreeStorage)
	ts.root = new(branch)

	return ts
}

func (ts *TreeStorage) AppendBanner(b *banners.Banner) {
	for _, category := range b.Categories {
		ts.root.add(b, []byte(category))
	}
}

func (ts *TreeStorage) LookupBanner(categories []string) (b *banners.Banner, ok bool) {

	if len(categories) == 0 {
		return ts.root.lookupAny(ts.CheckCount)
	}

	iter := NewRandIter(len(categories))
	for iter.Next() {
		cat := categories[iter.i]
		if l, ok := ts.root.lookup([]byte(cat), ts.CheckCount); ok {
			return l, ok
		}
	}

	return
}

func (ts *TreeStorage) CheckCount(l leaf) bool {
	return atomic.LoadInt64(&l.Count) > 0
}

//GetCount returned total count of banners in rotation
func (ts *TreeStorage) GetCount() int {
	return len(ts.root.leafs)
}
