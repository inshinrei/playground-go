package kv_db

import "encoding/binary"
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

func main() {
	old := BNode(make([]byte, BTREE_PAGE_SIZE))
	new := BNode(make([]byte, BTREE_PAGE_SIZE))
	new.setHeader(BNODE_LEAF, 3)

	nodeAppendKV(new, 0, 0, old.getKey(0), old.getValue(0))
	nodeAppendKV(new, 1, 0, []byte("k2"), []byte("b"))
	nodeAppendKV(new, 2, 0, old.getKey(2), old.getValue(2))

	new = make([]byte, BTREE_PAGE_SIZE)
	new.setHeader(BNODE_LEAF, 2)
	nodeAppendKV(new, 0, 0, old.getKey(0), old.getValue(0))
	nodeAppendKV(new, 1, 0, old.getKey(2), old.getValue(2))

	new = make([]byte, 2*BTREE_PAGE_SIZE)
	new.setHeader(BNODE_LEAF, 4)
	nodeAppendKV(new, 0, 0, []byte("a"), []byte("b"))
	nodeAppendKV(new, 1, 0, old.getKey(0), old.getValue(0))
	nodeAppendKV(new, 1, 0, old.getKey(1), old.getValue(1))
	nodeAppendKV(new, 3, 0, old.getKey(2), old.getValue(2))
}
