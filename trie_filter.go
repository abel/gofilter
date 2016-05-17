package gofilter

type TrieSlot struct {
	end   bool
	key   byte
	next  int32
	value TrieNode
}

type TrieNode struct {
	buckets []int32
	slots   []TrieSlot
}

// increaseCapacity
func (self *TrieNode) increaseCapacity(capacity int32) int32 {
	prime := GetPrimeInt32(capacity)
	self.buckets = make([]int32, prime, prime)
	//retset buckets and next
	size := int32(len(self.slots))
	for i := int32(0); i < size; i++ {
		index := int32(self.slots[i].key) % prime
		self.slots[i].next = self.buckets[index]
		self.buckets[index] = i + 1
	}
	return prime
}

// getNodeIndex 1: first index, 0: not find
func (self *TrieNode) getNodeIndex(key byte) int32 {
	capacity := len(self.buckets)
	if capacity > 0 {
		size := int32(len(self.slots))
		index := int(key) % capacity
		i := self.buckets[index]
		for i > 0 && i <= size {
			i--
			if self.slots[i].key == key {
				return i + 1
			}
			i = self.slots[i].next
		}
	}
	return 0
}

func (self *TrieNode) addChar(key byte) {
	//(1:1.5 = 67%)
	size := int32(len(self.slots))
	capacity := int32(len(self.buckets))
	cp := size + (size >> 1)
	if capacity <= cp {
		capacity = self.increaseCapacity(cp)
	}
	index := int32(key) % capacity
	//如果有冲突,则指向前一个的位置,如果无冲突,则为无效索引0
	//(因为buckets中保存的位置从1开始.next保存的位置从1开始)
	next := self.buckets[index]
	self.slots = append(self.slots, TrieSlot{key: key, next: next})
	self.buckets[index] = 1 + size
}

func (self *TrieNode) AddKeyword(keys []byte, trans []byte) {
	key_len := len(keys)
	if key_len == 0 {
		return
	}
	key := keys[0]
	if key == 0 {
		return
	}
	tran := trans[key]
	if tran != 0 {
		key = tran
	}
	i := self.getNodeIndex(key)
	if i == 0 {
		self.addChar(key)
		i = int32(len(self.slots))
	}
	i--
	if key_len == 1 {
		self.slots[i].end = true
		return
	}
	self.slots[i].value.AddKeyword(keys[1:], trans)
}

func (self *TrieNode) ExistKeyword(keys []byte, trans []byte) (find bool, depth int) {
	key_len := len(keys)
	if key_len == 0 {
		return false, 0
	}
	key := keys[0]
	tran := trans[key]
	var ignore bool
	if tran == 0 {
		ignore = true
	} else {
		ignore = false
		key = tran
	}
	i := self.getNodeIndex(key)
	if i == 0 {
		if ignore {
			find, depth = self.ExistKeyword(keys[1:], trans)
		}
	} else {
		if self.slots[i-1].end {
			find = true
		} else {
			find, depth = self.slots[i-1].value.ExistKeyword(keys[1:], trans)
		}
	}
	depth++
	return
}

const (
	CharCount = 256
)

type TrieFilter struct {
	transition [CharCount]byte
	rootNode   TrieNode
}

// SetFilter
func (self *TrieFilter) SetFilter(ignoreCase bool) {
	for i := 0; i < CharCount; i++ {
		self.transition[i] = byte(i)
	}
	if ignoreCase {
		for a := 'a'; a <= 'z'; a++ {
			self.transition[a] = byte(a) - 32 //(a:97,A:65);
		}
	}
}

// AddIgnoreChars
func (self *TrieFilter) AddIgnoreChars(passChars []byte) {
	for i := 0; i < len(passChars); i++ {
		src := passChars[i]
		self.transition[src] = 0
	}
}

// AddReplaceChars
func (self *TrieFilter) AddReplaceChars(srcChar []byte, replaceChar []byte) {
	count := len(srcChar)
	if count > len(replaceChar) {
		count = len(replaceChar)
	}
	for i := 0; i < count; i++ {
		src := srcChar[i]
		rep := replaceChar[i]
		self.transition[src] = rep
	}
}

// AddKey
func (self *TrieFilter) AddKey(key string) {
	self.rootNode.AddKeyword([]byte(key), self.transition[:])
}

// AddKeyword
func (self *TrieFilter) AddKeyword(key []byte) {
	self.rootNode.AddKeyword(key, self.transition[:])
}

// findKeywordIndex
func (self *TrieFilter) findKeywordIndex(text []byte) int {
	found, depth := self.rootNode.ExistKeyword(text, self.transition[:])
	if found {
		return depth
	} else {
		return 0
	}
}

// ExistKeyword
func (self *TrieFilter) ExistKeyword(text []byte) bool {
	for i := 0; i < len(text); i++ {
		index := self.findKeywordIndex(text[i:])
		if index > 0 {
			return true
		}
	}
	return false
}

// FindOne
func (self *TrieFilter) FindOne(text []byte) []byte {
	for i := 0; i < len(text); i++ {
		index := self.findKeywordIndex(text[i:])
		if index > 0 {
			return text[i : i+index]
		}
	}
	return nil
}

// FindAll
func (self *TrieFilter) FindAll(text []byte) [][]byte {
	var r [][]byte
	for i := 0; i < len(text); i++ {
		index := self.findKeywordIndex(text[i:])
		if index > 0 {
			r = append(r, text[i:i+index])
			i += (index - 1)
		}
	}
	return r
}

// Replace
func (self *TrieFilter) Replace(text []byte, mask byte) (int, []byte) {
	textLen := len(text)
	var outbuffer []byte
	findCount := 0
	for i := 0; i < len(text); i++ {
		index := self.findKeywordIndex(text[i:])
		if index > 0 {
			if findCount == 0 {
				outbuffer = make([]byte, 0, textLen)
				outbuffer = append(outbuffer, text[0:i]...)
			}
			findCount++
			outbuffer = append(outbuffer, mask)
			i += (index - 1)
		} else {
			if findCount > 0 {
				outbuffer = append(outbuffer, text[i])
			}
		}
	}
	if findCount == 0 {
		return 0, text
	}
	return findCount, outbuffer
}
