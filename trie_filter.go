package gofilter

type TrieSlot struct {
	key   byte
	next  int32
	value *TrieNode
}

type TrieNode struct {
	capacity int32
	size     int32
	buckets  []int32
	slots    []TrieSlot
	end      bool
}

func (node *TrieNode) increaseCapacity(capacity int32) {
	size := node.size
	prime := GetPrimeInt32(capacity)
	newSlots := make([]TrieSlot, prime, prime)
	if node.slots != nil {
		copy(newSlots, node.slots)
	}
	newBuckets := make([]int32, prime, prime)
	//重新设置buckets和next
	for i := int32(0); i < size; i++ {
		index := int32(newSlots[i].key) % prime
		newSlots[i].next = newBuckets[index]
		newBuckets[index] = i + 1
	}
	node.slots = newSlots
	node.buckets = newBuckets
	node.capacity = prime
}

func (node *TrieNode) addKey(key byte, trieNode *TrieNode) {
	//容量因子(1:1.5 = 67%)
	cp := node.size + (node.size >> 1)
	if node.capacity <= cp {
		node.increaseCapacity(cp)
	}
	index := int32(key) % node.capacity
	curSlot := &node.slots[node.size]
	curSlot.key = key
	curSlot.value = trieNode
	//如果有冲突,则指向前一个的位置,如果无冲突,则为无效索引0
	//(因为buckets中保存的位置从1开始.next保存的位置从1开始)
	curSlot.next = node.buckets[index]
	node.size++
	node.buckets[index] = node.size
}

func (node *TrieNode) GetValue(key byte) *TrieNode {
	if node.slots != nil {
		var curSlot *TrieSlot
		index := int32(key) % node.capacity
		for i := node.buckets[index]; i > 0; i = curSlot.next {
			curSlot = &node.slots[i-1]
			if curSlot.key == key {
				return curSlot.value
			}
		}
	}
	return nil
}

func (node *TrieNode) GetValueOrNew(key byte) *TrieNode {
	trieNode := node.GetValue(key)
	if trieNode == nil {
		trieNode = new(TrieNode)
		node.addKey(key, trieNode)
	}
	return trieNode
}

const (
	CharCount = 256
)

type TrieFilter struct {
	transition [CharCount]byte
	root_node  TrieNode //根节点
	//nodes      [CharCount]*TrieNode
}

func (filter *TrieFilter) SetFilter(ignoreCase bool, ignoreSimpTrad bool) {
	for i := 0; i < CharCount; i++ {
		filter.transition[i] = byte(i)
	}
	//将小写转为大写字母
	if ignoreCase {
		for a := 'a'; a <= 'z'; a++ {
			filter.transition[a] = byte(a) - 32 //(a:97,A:65);
		}
	}
	//简繁转换.暂未实现
	//if (ignoreSimpTrad)
	//{
	//	AddReplaceChars(zh_TW, zh_CN);
	//}
}

func GetNonzeroByte(a byte, b byte) byte {
	if a != 0 {
		return a
	}
	return b
}

// 查找关键字的位置
func (filter *TrieFilter) FindBadWordIndex(text []byte) int {
	textLen := len(text)
	node := &filter.root_node
	for index := 0; index < textLen && node != nil; index++ {
		src := text[index]
		if src == 0 {
			break
		}
		tranc := filter.transition[src]
		if tranc != 0 {
			nextNode := node.GetValue(tranc)
			if nextNode != nil {
				node = nextNode
			} else {
				break
			}
		} else {
			//被忽略的字符.使用原始值再次查找
			nextNode := node.GetValue(src)
			if nextNode != nil {
				node = nextNode
			}
		}
		if node.end {
			return index + 1
		}
	}
	return -1
}

//增加忽略字符
func (filter *TrieFilter) AddIgnoreChars(passChars []byte) {
	for i := 0; i < len(passChars); i++ {
		src := passChars[i]
		filter.transition[src] = 0
	}
}

//增加替换字符
func (filter *TrieFilter) AddReplaceChars(srcChar []byte, replaceChar []byte) {
	count := len(srcChar)
	if count > len(replaceChar) {
		count = len(replaceChar)
	}
	for i := 0; i < count; i++ {
		src := srcChar[i]
		rep := replaceChar[i]
		filter.transition[src] = rep
	}
}

//添加关键字
func (filter *TrieFilter) AddKey(key []byte) {
	keyLen := len(key)
	node := &filter.root_node
	for i := 0; i < keyLen; i++ {
		src := key[i]
		if src == 0 {
			break
		}
		index := GetNonzeroByte(filter.transition[src], src)
		node = node.GetValueOrNew(index)
	}
	node.end = true
	filter.root_node.end = false
}

// 存在过滤字
func (filter *TrieFilter) HasBadWord(text []byte) bool {
	for i := 0; i < len(text); i++ {
		index := filter.FindBadWordIndex(text[i:])
		if index > 0 {
			return true
		}
	}
	return false
}

//查找1个
func (filter *TrieFilter) FindOne(text []byte) []byte {
	for i := 0; i < len(text); i++ {
		index := filter.FindBadWordIndex(text[i:])
		if index > 0 {
			return text[i : i+index]
		}
	}
	return nil
}

func (filter *TrieFilter) FindAll(text []byte) [][]byte {
	var r [][]byte
	for i := 0; i < len(text); i++ {
		index := filter.FindBadWordIndex(text[i:])
		if index > 0 {
			r = append(r, text[i:i+index])
			i += (index - 1)
		}
	}
	return r
}

func (filter *TrieFilter) Replace(text []byte, mask byte) (int, []byte) {
	textLen := len(text)
	var outbuffer []byte
	findCount := 0
	for i := 0; i < len(text); i++ {
		index := filter.FindBadWordIndex(text[i:])
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
