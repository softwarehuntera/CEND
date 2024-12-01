package Runetree

import (
	"cend/database/tokenizer"
	"cend/database/document"
)


// Introduction to notation/terminology
// text: Any kind of string
// token: A subset of text that we use for data storage

// Goal 1: Implement a very simple trie-like structure in memory
// Goal 2: Implement a file store version of the trie

// The idea of implementing an 'analyzer' is that all strings are broken up
// into a predictable set of Runes.
// 
type RuneTree struct {
	rootNodes map[rune]RuneTreeNode
}

func New() RuneTree {
	return RuneTree{
		rootNodes: map[rune]RuneTreeNode{},
	}
}

// Idea: provide the entire search 'text' so we can search multiple tokens at once.
func (rt *RuneTree) Search(text string) []RuneTreeNode {
	nGrams := tokenizer.NGrams(text, 3)
	
	matches := []string
	for _, nGram := range nGrams {

	}
	// sort by tokens and weight by frequency

	// search recursively, giving all search tokens that match the prefix string
	
}

func (rt *RuneTree) Add(token string) {
	// search for current matching token
	rt.Search(token)

}


type RuneTreeNode struct {
	Rune string
	Size int32
	DocIDs *document.DocumentIDs
}

func (rtn *RuneTreeNode) Search(tokens []string) {

}

// Almost trie structure or Rune index

// How do I implement concurrency to read multiple references of Rune_tree at once?