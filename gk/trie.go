package gk

import (
	"fmt"
	"strings"
)

// symbol ':' matches parameter
// symbol '*' matches all described string
type node struct {
	pattern  string // router to be matched
	part     string // part of router
	children []*node
	isWild   bool // true while part has ':' or '*'
}

func (t *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", t.pattern, t.part, t.isWild)
}

// search returns node that matches parts
// it will search part by part
// when it reaches the end of parts, returns current node onyl if pattern is not nil
func (t *node) search(parts []string, height int) *node {
	// found
	if len(parts) == height || strings.HasPrefix(t.part, "*") {
		if t.pattern == "" {
			return nil
		}
		return t
	}

	part := parts[height]
	children := t.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// insert searches for a node that matches slice parts
// as it reaches the end of parts, set pattern to this node
// if there is no matched subnode, insert a new node and continue
func (t *node) insert(pattern string, parts []string, height int) {
	// found
	if len(parts) == height {
		t.pattern = pattern
		return
	}

	part := parts[height]
	child := t.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		t.children = append(t.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (t *node) travel(list *([]*node)) {
	if t.pattern != "" {
		*list = append(*list, t)
	}
	for _, child := range t.children {
		child.travel(list)
	}
}

// matchChild return the first matched node
func (t *node) matchChild(part string) *node {
	for _, child := range t.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren returns all matched nodes
func (t *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range t.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
