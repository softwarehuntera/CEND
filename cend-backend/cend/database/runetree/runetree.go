package runetree

import (
	"cend/database/document"
	"cend/database/tokenizer"
	"fmt"
	"unicode/utf8"
)

// Introduction to notation/terminology
// text: Any kind of string
// token: A subset of text that we use for data storage

// Goal 1: Implement a very simple trie-like structure in memory
// Goal 2: Implement a file store version of the trie

// The idea of implementing an 'analyzer' is that all strings are broken up
// into a predictable set of Runes.
//

type SearchResult struct {
	ID int `json:"id"`
	Document string  `json:"document"`
	Score    float64 `json:"score"`
}

type Node interface {
	Children() map[rune]*RuneTreeNode
}

type RuneTree struct {
	children map[rune]*RuneTreeNode
}

func (rt *RuneTree) Children() map[rune]*RuneTreeNode {
	return rt.children
}


type RuneTreeNode struct {
	Rune rune
	Size int32
	DocIDs *document.DocumentIDs
	children map[rune]*RuneTreeNode
}

func (rt *RuneTreeNode) Children() map[rune]*RuneTreeNode {
	return rt.children
}
func New() RuneTree {
	return RuneTree{
		rootNodes: map[rune]RuneTreeNode{},
	}
}

// Idea: provide the entire search 'text' so we can search multiple tokens at once.
func Search(text string, rt *RuneTree, docs ) []RuneTreeNode {
	nGrams := tokenizer.NGrams(text, 3)
	var node *RuneTreeNode
	
	docIds := map[int]struct{}{}
	for _, token := range nGrams {
		for index, runeValue := range token {
			if index == 0 {
				if node, exists := rt.Children()[runeValue]; exists {
					addDocIds(node, docIds)
				} else {
					break
				}
			} else {
				if node, exists := node.Children()[runeValue]; exists {
					addDocIds(node, docIds)
				} else {
					break
				}
			}
		}
	}


	pq := PriorityQueue{}
	// sort by tokens and weight by frequency

	// search recursively, giving all search tokens that match the prefix string
	
}

func addDocIds(node *RuneTreeNode, docIds map[int]struct{}) {
	if node.DocIDs == nil || node.DocIDs.Ids == nil || len(node.DocIDs.Ids) != 0 {
		return
	}
	for id, _ := range node.DocIDs.Ids {
		docIds[id] = struct{}{}
	}
}

func (rt *RuneTree) Add(token string) {
	// search for current matching token
	rt.Search(token)

}


func (rtn *RuneTreeNode) Search(tokens []string) {

}

// Almost trie structure or Rune index

// How do I implement concurrency to read multiple references of Rune_tree at once?