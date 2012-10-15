// Copyright 2011 Numrotron Inc.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.
//
// Developed at www.stathat.com by Patrick Crosby
// Contact us on twitter with any questions:  twitter.com/stat_hat

// The treap package provides a balanced binary tree datastructure, expected
// to have logarithmic height.
package treap

import (
	"math"
	"math/rand"
	"time"
	"fmt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// A Tree is the data structure that stores everything
type Tree struct {
	less    LessFunc
	overlap OverlapFunc
	count   int
	root    *Node
	max_doc int
	min_doc int
}

// LessFunc returns true if a < b
type LessFunc func(a, b interface{}) bool

// OverlapFunc return true if a overlaps b.  Optional.  Only used by overlap trees.
type OverlapFunc func(a, b interface{}) bool

// Key can be anything.  It will use LessFunc to compare keys.
type Key interface{}

// Item can be anything.
type Item interface{}

// A Node in the Tree.
type Node struct {
	key      int
	item     int
	priority int
	delta int
	deltakey int
	left     *Node
	right    *Node
}

func newNode(key int, item int, priority int) *Node {
	result := new(Node)
	result.key = key
	result.item = item
	result.priority = priority
	return result
}

func (n *Node) GetKey() int {
	return n.key
}

func (n *Node) GetPriority() int {
	return n.priority
}

func (n *Node) GetLeft() *Node {
	return n.left
}
func (n *Node) GetRight() *Node {
	return n.right
}

// To create a Tree, you need to supply a LessFunc that can compare the
// keys in the Node.
func NewTree(lessfn LessFunc) *Tree {
	t := new(Tree)
	t.max_doc = 0
	t.min_doc = 1<<31 - 1
	t.less = lessfn
	return t
}

// To create a tree that also lets you iterate by key overlap, supply a LessFunc
// and an OverlapFunc
func NewOverlapTree(lessfn LessFunc, overfn OverlapFunc) *Tree {
	t := new(Tree)
	t.less = lessfn
	t.overlap = overfn
	return t
}

// Remove everything from the tree.
func (t *Tree) Reset() {
	t.root = nil
	t.count = 0
}

// The length of the tree (number of nodes).
func (t *Tree) Len() int {
	return t.count
}

func (t *Tree) GetRootKey() int {
	return t.root.key
}
func (t *Tree) GetRootPriority() int {
	return t.root.priority
}
func (t *Tree) GetRoot() *Node {
	return t.root
}

// Get an Item in the tree.
func (t *Tree) Get(key int) Item {
	return t.get(t.root, key)
}

func (t *Tree) CompressPriority(node *Node) {
	if node == nil {
		return
	}
	if node == t.root {
		t.root.delta = t.root.priority
		if (node.left != nil ) { 
			node.left.delta = node.left.priority - node.priority	
		} else if (node.right != nil ) { 
			node.right.delta = node.right.priority - node.priority
		}
	} else {
		if (node.left != nil ) { 
			node.left.delta = node.left.priority - node.priority	
		} else if (node.right != nil ) { 
			node.right.delta = node.right.priority - node.priority
		}
	}
	t.CompressPriority(node.left)
	t.CompressPriority(node.right)
}

func (t *Tree) CompressKeys(node *Node) {
	if node == nil {
		return
	}
	if node == t.root {
		t.root.deltakey = t.root.key
		if (node.left != nil ) { 
			node.left.deltakey = node.left.key - node.key	
		} else if (node.right != nil ) { 
			node.right.deltakey = node.right.key - node.key
		}
	} else {
		if (node.left != nil ) { 
			node.left.deltakey = node.left.key - node.key	
		} else if (node.right != nil ) { 
			node.right.deltakey = node.right.key - node.key
		}
	}
	t.CompressKeys(node.left)
	t.CompressKeys(node.right)
}


func (t *Tree) PrintPriorities(node *Node,height int) {
	if node == nil {
		return
	}
	fmt.Println("node:",node.key, " priority = ", node.priority, "height=",height)
	height+=1
	t.PrintPriorities(node.left,height)
	t.PrintPriorities(node.right,height)
}
func (t *Tree) PrintDeltas(node *Node,height int,acum []int) {
	if node == nil {
		return
	}
	//fmt.Println("node:",node.key, " priority = ", node.priority,"delta=",node.delta, "height=",height)
	acum = append(acum,node.delta)
	height+=1
	t.PrintDeltas(node.left,height,acum)
	t.PrintDeltas(node.right,height,acum)
}
func (t *Tree) GetKeyPriorityDelta(node *Node,key int) int {
	if t.less(key,node.key) {
		if (node.left != nil) {
			return t.GetKeyPriorityDelta(node.left,key)
		}
	}
	if t.less(node.key,key) {
		if (node.right != nil) {
			return t.GetKeyPriorityDelta(node.right,key)
		}
	}
	return node.delta
}

func (t *Tree) GetKeyPriority(node *Node,key int) int {
	if t.less(key,node.key) {
		if (node.left != nil) {
			return t.GetKeyPriority(node.left,key)
		}
	}
	if t.less(node.key,key) {
		if (node.right != nil) {
			return t.GetKeyPriority(node.right,key)
		}
	}
	return node.priority
}

func (t *Tree) GetKeyDelta(node *Node,key int) int {
	if t.less(key,node.key) {
		if (node.left != nil) {
			return t.GetKeyDelta(node.left,key)
		}
	}
	if t.less(node.key,key) {
		if (node.right != nil) {
			return t.GetKeyDelta(node.right,key)
		}
	}
	return node.deltakey
}

func (t *Tree) get(node *Node, key int) Item {
	if node == nil {
		return nil
	}
	if t.less(key, node.key) {
		return t.get(node.left, key)
	}
	if t.less(node.key, key) {
		return t.get(node.right, key)
	}
	return node.item
}

// Returns true if there is an item in the tree with this key.
func (t *Tree) Exists(key Key) bool {
	return t.exists(t.root, key)
}

func (t *Tree) exists(node *Node, key Key) bool {
	if node == nil {
		return false
	}
	if t.less(key, node.key) {
		return t.exists(node.left, key)
	}
	if t.less(node.key, key) {
		return t.exists(node.right, key)
	}
	return true
}

// Insert an item into the tree.
func (t *Tree) Insert(key int, item int) {
	if t.max_doc < int(key) {
		t.max_doc = int(key)
	}
	if t.min_doc > int(key) {
		t.min_doc = int(key)
	}
	priority := int(item)
	t.root = t.insert(t.root, key, item, priority)
}

func (t *Tree) insert(node *Node, key int, item int, priority int) *Node {
	if node == nil {
		t.count++
		return newNode(key, item, priority)
	}
	if t.less(key, node.key) {
		node.left = t.insert(node.left, key, item, priority)
		if node.left.priority > node.priority {
			return t.leftRotate(node)
		}
		return node
	}
	if t.less(node.key, key) {
		node.right = t.insert(node.right, key, item, priority)
		if node.right.priority > node.priority {
			return t.rightRotate(node)
		}
		return node
	}

	// equal: replace the value
	node.item = item
	return node
}

func (t *Tree) leftRotate(node *Node) *Node {
	result := node.left
	x := result.right
	result.right = node
	node.left = x
	return result
}

func (t *Tree) rightRotate(node *Node) *Node {
	result := node.right
	x := result.left
	result.left = node
	node.right = x
	return result
}

// Split the tree by creating a tree with a node of priority -1 so it will be the root
func (t *Tree) Split(key int) (*Node, *Node) {
	inserted := t.insert(t.root, key, -1, -1)
	return inserted.left, inserted.right
}

// Merge two trees together by supplying the root node of each tree.
func (t *Tree) Merge(left, right *Node) *Node {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	if left.priority < right.priority {
		result := left
		x := left.right
		result.right = t.Merge(x, right)
		return result
	}

	result := right
	x := right.left
	result.left = t.Merge(x, left)
	return result
}

// Delete the item from the tree that has this key.
func (t *Tree) Delete(key Key) {
	if t.Exists(key) == false {
		return
	}
	t.root = t.delete(t.root, key)
}

func (t *Tree) delete(node *Node, key Key) *Node {
	if node == nil {
		panic("key not found")
	}

	if t.less(key, node.key) {
		result := node
		x := node.left
		result.left = t.delete(x, key)
		return result
	}
	if t.less(node.key, key) {
		result := node
		x := node.right
		result.right = t.delete(x, key)
		return result
	}
	t.count--
	return t.Merge(node.left, node.right)
}

func (t *Tree) GetHeightTree(node *Node) int {
	if node == nil {
		return 0
	}
	return int(1 + math.Max(float64(t.GetHeightTree(node.left)), float64(t.GetHeightTree(node.right))))
}

// Returns the height (depth) of the tree.
func (t *Tree) Height(key Key) int {
	return t.height(t.root, key)
}

func (t *Tree) height(node *Node, key Key) int {
	if node == nil {
		return 0
	}
	if t.less(key, node.key) {
		depth := t.height(node.left, key)
		return depth + 1
	}
	if t.less(node.key, key) {
		depth := t.height(node.right, key)
		return depth + 1
	}
	return 0
}

// Returns a channel of Items whose keys are in ascending order.
func (t *Tree) IterAscend() <-chan Item {
	c := make(chan Item)
	go func() {
		iterateInOrder(t.root, c)
		close(c)
	}()
	return c
}

func iterateInOrder(h *Node, c chan<- Item) {
	if h == nil {
		return
	}
	iterateInOrder(h.left, c)
	c <- h.item
	iterateInOrder(h.right, c)
}

// Returns a channel of keys in ascending order.
func (t *Tree) IterKeysAscend() <-chan Key {
	c := make(chan Key)
	go func() {
		iterateKeysInOrder(t.root, c)
		close(c)
	}()
	return c
}

func iterateKeysInOrder(h *Node, c chan<- Key) {
	if h == nil {
		return
	}
	iterateKeysInOrder(h.left, c)
	c <- h.key
	iterateKeysInOrder(h.right, c)
}

// Returns a channel of items that overlap key.
func (t *Tree) IterateOverlap(key Key) <-chan Item {
	c := make(chan Item)
	go func() {
		if t.overlap != nil {
			t.iterateOverlap(t.root, key, c)
		}
		close(c)
	}()
	return c
}

func (t *Tree) iterateOverlap(h *Node, key Key, c chan<- Item) {
	if h == nil {
		return
	}
	lessThanLower := t.overlap(h.key, key)
	greaterThanUpper := t.overlap(key, h.key)

	if !lessThanLower {
		t.iterateOverlap(h.left, key, c)
	}
	if !lessThanLower && !greaterThanUpper {
		c <- h.item
	}
	if !greaterThanUpper {
		t.iterateOverlap(h.right, key, c)
	}
}

// Returns the minimum item in the tree.
func (t *Tree) Min() Item {
	return min(t.root)
}

func min(h *Node) Item {
	if h == nil {
		return nil
	}
	if h.left == nil {
		return h.item
	}
	return min(h.left)
}

// Returns the maximum item in the tree.
func (t *Tree) Max() Item {
	return max(t.root)
}
func (t *Tree) GetMaxDoc() int {
	return t.max_doc
}

func (t *Tree) GetMinDoc() int {
	return t.min_doc
}

func max(h *Node) Item {
	if h == nil {
		return nil
	}
	if h.right == nil {
		return h.item
	}
	return max(h.right)
}

func IntLess(p, q interface{}) bool {
	return p.(int) < q.(int)
}

type TreapList struct {
	Treap []*Tree
}

func (t *Tree) GetCount() int {
	return t.count
}

func (s *TreapList) Len() int {
	return len(s.Treap)
}
func (s *TreapList) Less(i, j int) bool {
	return s.Treap[i].count < s.Treap[j].count
}
func (s *TreapList) Swap(i, j int) {
	s.Treap[i], s.Treap[j] = s.Treap[j], s.Treap[i]
}

func Intersect(terms *TreapList, y0 int, y1 int, min_freq int) *Tree {
	new_max := int(math.Min(float64(terms.Treap[0].GetMaxDoc()), float64(terms.Treap[1].GetMaxDoc())))
	new_min := int(math.Max(float64(terms.Treap[0].GetMinDoc()), float64(terms.Treap[1].GetMinDoc())))

	for i := 2; i < len(terms.Treap); i++ {
		new_max = int(math.Min(float64(new_max), float64(terms.Treap[i].GetMaxDoc())))
		new_min = int(math.Max(float64(new_min), float64(terms.Treap[i].GetMinDoc())))
	}

	if new_min > new_max {
		return nil
	}

	res := NewTree(IntLess)
	IntersectNode(terms.Treap[0].GetRoot(), terms.Treap[1].GetRoot(), terms.Treap[0].GetRoot().GetKey(), y0, y1, 0, min_freq, new_min, new_max, res)
	for i := 2; i < len(terms.Treap)-2; i++ {
		if res.GetRoot() == nil {
			return nil
		}
		res_aux := NewTree(IntLess)
		IntersectNode(res.GetRoot(), terms.Treap[i].GetRoot(), res.GetRoot().GetKey(), y0, y1, 0, min_freq, new_min, new_max, res_aux)

		res = res_aux
	}
	return res
}

func IntersectNode(Node1 *Node, Node2 *Node, key int, y0 int, y1 int, depth int, min_freq int, min_doc int, max_doc int, result *Tree) {
	if Node1 == nil {
		return
	}

	if Node2 == nil {
		return
	}

	// if (levelfreq[depth] < min_freq) // if the frequency of the level is not enough
	// 	return

	if key < min_doc {
		return
	}

	if key > max_doc {
		return
	}

	if key == Node2.GetKey() { // if we found a result 
		result.Insert(key, Node1.GetPriority()+Node2.GetPriority()) // insert in the auxilary tree, the sum of the frequencies

		// we have to continue to the next value of the tree
		if Node1.GetLeft() != nil {
			if Node1.GetLeft().GetKey() < y1 { // if left is not out of range
				depth++
				IntersectNode(Node1.GetLeft(), Node2, Node1.GetLeft().GetKey(), y0, y1, depth, min_freq, min_doc, max_doc, result)
			}
		}

		if Node1.GetRight() != nil {
			if Node1.GetRight().GetKey() > y0 { // if right is not out of range
				depth++
				IntersectNode(Node1.GetRight(), Node2, Node1.GetRight().GetKey(), y0, y1, depth, min_freq, min_doc, max_doc, result)
			}
		}
	}

	if Node2.GetLeft() != nil {
		if key < Node2.GetKey() { // if key is smaller than the key of the other tree
			if Node2.GetLeft().GetKey() < y1 { // if we are not out of bounds
				depth++
				IntersectNode(Node1, Node2.GetLeft(), key, y0, y1, depth, min_freq, min_doc, max_doc, result) // move the node of the other tree
			}
		}
	}
	if Node2.GetRight() != nil {
		if key > Node2.GetKey() { // if the key is greater than the key of the other tree
			if Node2.GetRight().GetKey() > y0 { // if we are not out of bounds
				depth++
				IntersectNode(Node1, Node2.GetRight(), key, y0, y1, depth, min_freq, min_doc, max_doc, result) // mode the node of the other tree
			}
		}
	}
}
