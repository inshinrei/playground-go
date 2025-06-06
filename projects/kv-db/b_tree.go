package kv_db

import (
	"bytes"
	"encoding/binary"
)
import (
	"kv-db/util"
)

type Node struct {
	keys     [][]byte
	values   [][]byte
	children []*Node
}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VALUE_SIZE = 3000

func init() {
	//node1max := 4 + 1*8 + 1*2 + 4 + BTREE_MAX_KEY_SIZE + BTREE_MAX_VALUE_SIZE
	// assert(node1max <= BTREE_PAGE_SIZE)
}

type BNode []byte

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys())
}

func (node BNode) setHeader(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}

func (node BNode) getPtr(idx uint16) uint64 {
	util.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, value uint64) {
	util.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], value)
}

func (node BNode) setValue(idx uint16, value uint64) {
	util.Assert(idx < node.nkeys())
	pos := 4 + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], value)
}

func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	return binary.LittleEndian.Uint16(node[pos:])
}

func (node BNode) setOffset(idx uint16, offset uint16) {
	util.Assert(idx < node.nkeys())
	pos := 4 + 8*node.nkeys() + 2*(idx-1)
	binary.LittleEndian.PutUint16(node[pos:], offset)
}

func (node BNode) kvPos(idx uint16) uint16 {
	util.Assert(idx <= node.nkeys())
	return 4 + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	util.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos+4:][:klen]
}

func (node BNode) getValue(idx uint16) []byte {
	util.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos+0:])
	vlen := binary.LittleEndian.Uint16(node[pos+2:])
	return node[pos+4+klen:][:vlen]
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, value []byte) {
	new.setPtr(idx, ptr)
	pos := new.kvPos(idx)
	binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
	binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(value)))
	copy(new[pos+4:], key)
	copy(new[pos+4+uint16(len(key)):], value)
	new.setOffset(idx+1, new.getOffset(idx)+4+uint16(len(key)+len(value)))
}

func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	for i := uint16(0); i < n; i++ {
		dst, src := dstNew+i, srcOld+i
		nodeAppendKV(new, dst, old.getPtr(src), old.getKey(src), old.getValue(src))
	}
}

func leafInsert(new BNode, old BNode, idx uint16, key []byte, value []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys()+1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, value)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

func leafUpdate(new BNode, old BNode, idx uint16, key []byte, value []byte) {
	new.setHeader(BNODE_LEAF, old.nkeys())
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, value)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-(idx+1))
}

func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	var i uint16
	for i = 0; i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)
		if cmp == 0 {
			return i
		}
		if cmp > 0 {
			return i - 1
		}
	}
	return i - 1
}

func nodeSplit2(left BNode, right BNode, old BNode) {
	util.Assert(old.nkeys() >= 2)

	nleft := old.nkeys() / 2
	left_bytes := func() uint16 {
		return 4 + 8*nleft + 2*nleft + old.getOffset(nleft)
	}
	for left_bytes() > BTREE_PAGE_SIZE {
		nleft--
	}

	util.Assert(nleft >= 1)

	right_bytes := func() uint16 {
		return old.nbytes() - left_bytes() + 4
	}
	for right_bytes() > BTREE_PAGE_SIZE {
		nleft++
	}

	util.Assert(nleft < old.nkeys())

	nright := old.nkeys() - nleft
	left.setHeader(old.btype(), nleft)
	right.setHeader(old.btype(), nright)
	nodeAppendRange(left, old, 0, 0, nleft)
	nodeAppendRange(right, old, 0, nleft, nright)

	util.Assert(right.nbytes() <= BTREE_PAGE_SIZE)
}

func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}
	left := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	right := BNode(make([]byte, BTREE_PAGE_SIZE))
	nodeSplit2(left, right, old)
	if left.nbytes() <= BTREE_PAGE_SIZE {
		left = left[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}
	leftleft := BNode(make([]byte, BTREE_PAGE_SIZE))
	middle := BNode(make([]byte, BTREE_PAGE_SIZE))
	nodeSplit2(leftleft, middle, left)

	util.Assert(left.nbytes() <= BTREE_PAGE_SIZE)

	return 3, [3]BNode{leftleft, middle, right}
}

type BTree struct {
	root uint64
	get  func(uint64) []byte
	new  func([]byte) uint64
	del  func(uint64)
}

func (tree *BTree) Insert(key []byte, value []byte) error {
	//if err := checkLimit(key, value); err != nil {
	//	return err
	//}
	if tree.root == 0 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeader(BNODE_LEAF, 2)
		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, value)
		tree.root = tree.new(root)
		return nil
	}
	node := treeInsert(tree, tree.get(tree.root), key, value)
	nsplit, split := nodeSplit3(node)
	tree.del(tree.root)
	if nsplit > 1 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeader(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
	} else {
		tree.root = tree.new(split[0])
	}
	return nil
}

var HEADER uint16 = 1

func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode) {
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}
	if idx > 0 {
		sibling := BNode(tree.get(node.getPtr(idx - 1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}
	}
	if idx+1 < node.nkeys() {
		sibling := BNode(tree.get(node.getPtr(idx + 1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return +1, sibling
		}
	}
	return 0, BNode{}
}

func nodeDelete(tree *BTree, node BNode, idx uint16, key []byte) BNode {

}

func nodeReplaceKidN(tree *BTree, new BNode, old BNode, idx uint16, kids ...BNode) {
	inc := uint16(len(kids))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}

func treeInsert(tree *BTree, node BNode, key []byte, value []byte) BNode {
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	idx := nodeLookupLE(node, key)
	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(key, node.getKey(idx)) {
			leafUpdate(new, node, idx, key, value)
		} else {
			leafInsert(new, node, idx+1, key, value)
		}
	case BNODE_NODE:
		kptr := node.getPtr(idx)
		knode := treeInsert(tree, tree.get(kptr), key, value)
		nsplit, split := nodeSplit3(knode)
		tree.del(kptr)
		nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
	}
	return new
}
