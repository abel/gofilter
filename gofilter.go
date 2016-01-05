package gofilter

import (
	"io/ioutil"
)

var (
	trie      TrieFilter
	trie_name TrieFilter
)

func init() {
	trie.SetFilter(true, true)
	trie_name.SetFilter(true, true)
}

func (filter *TrieFilter) LoadMaskFile(path string) bool {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	index := 0
	size := 0
	for i := 0; i < len(content); i++ {
		if content[i] == 0x0A {
			filter.AddKey(content[index : index+size])
			size = 0
			index = i + 1
		} else {
			size++
		}
	}
	if size > 0 && index < len(content) {
		filter.AddKey(content[index : index+size])
	}
	return true
}

func LoadMaskWordFile(path string) {
	trie.AddIgnoreChars([]byte(" *&^%$#@!~,.:[]{}?+-~\"\\"))
	trie.LoadMaskFile(path)
}

func LoadMaskNameFile(path string) {
	trie_name.LoadMaskFile(path)
}

func TrieHasBadWord(text string) bool {
	return trie.HasBadWord([]byte(text))
}

func TrieHasBadName(text string) bool {
	t := []byte(text)
	return trie.HasBadWord(t) || trie_name.HasBadWord(t)
}

func TrieReplaceBadWord(text string) string {
	count, outbuffer := trie.Replace([]byte(text), '*')
	if count == 0 {
		return text
	}
	return string(outbuffer)
}
