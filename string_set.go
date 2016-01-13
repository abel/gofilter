package gofilter

import (
	"math"
)

var (
	primes = []int{3, 7, 11, 0x11, 0x17, 0x1d, 0x25, 0x2f, 0x3b, 0x47, 0x59, 0x6b, 0x83, 0xa3, 0xc5, 0xef,
		0x125, 0x161, 0x1af, 0x209, 0x277, 0x2f9, 0x397, 0x44f, 0x52f, 0x63d, 0x78b, 0x91d, 0xaf1, 0xd2b, 0xfd1, 0x12fd,
		0x16cf, 0x1b65, 0x20e3, 0x2777, 0x2f6f, 0x38ff, 0x446f, 0x521f, 0x628d, 0x7655, 0x8e01, 0xaa6b, 0xcc89, 0xf583, 0x126a7, 0x1619b,
		0x1a857, 0x1fd3b, 0x26315, 0x2dd67, 0x3701b, 0x42023, 0x4f361, 0x5f0ed, 0x72125, 0x88e31, 0xa443b, 0xc51eb, 0xec8c1, 0x11bdbf, 0x154a3f, 0x198c4f,
		0x1ea867, 0x24ca19, 0x2c25c1, 0x34fa1b, 0x3f928f, 0x4c4987, 0x5b8b6f, 0x6dda89}
)

func IsPrime(candidate int) bool {
	if candidate&1 == 0 {
		return candidate == 2
	}
	num := int(math.Sqrt(float64(candidate)))
	for i := 3; i <= num; i += 2 {
		if (candidate % i) == 0 {
			return false
		}
	}
	return true
}

func GetPrime(min int) int {
	//在已知数据中查找
	for _, p := range primes {
		if p >= min {
			return p
		}
	}
	//通过计算获得
	for i := min | 1; i < 0x7fffffff; i += 2 {
		if IsPrime(i) {
			return i
		}
	}
	return 0
}

func GetPrimeInt32(min int32) int32 {
	return int32(GetPrime(int(min)))
}

func internalGetHashCode(key []byte) int {
	count := len(key)
	h := count
	for i := 1; i < count; i += 2 {
		h = (h << 5) - h + int(key[i-1])
		h = (h << 5) - h + int(key[i])
	}
	if count%2 == 1 {
		h = (h << 5) - h + int(key[count-1])
	}
	return h & 0x7fffffff
}

func StringEquals(s1 []byte, s2 []byte) bool {
	count := len(s1)
	if count == len(s2) {
		for i := 0; i < count; i++ {
			if s1[i] != s2[i] {
				return false
			}
		}
		return true
	}
	return false
}

type Slot struct {
	hashCode int
	next     int
	value    []byte
}

type StringSet struct {
	capacity int
	size     int
	buckets  []int
	slots    []Slot
}

func (set *StringSet) increaseCapacity(capacity int) {
	size := set.size
	prime := GetPrime(capacity)

	newSlots := make([]Slot, prime, prime)
	if size > 0 {
		copy(newSlots, set.slots)
	}

	newBuckets := make([]int, prime, prime)
	for i := 0; i < size; i++ {
		index := newSlots[i].hashCode % prime
		newSlots[i].next = newBuckets[index] - 1
		newBuckets[index] = i + 1
	}
	set.slots = newSlots
	set.buckets = newBuckets
	set.capacity = prime
}

func (set *StringSet) Contains(key []byte) bool {
	if set.buckets != nil {
		hashCode := internalGetHashCode(key)
		var curSlot *Slot
		for i := set.buckets[hashCode%set.capacity] - 1; i >= 0; i = curSlot.next {
			curSlot = &set.slots[i]
			if curSlot.hashCode == hashCode && StringEquals(curSlot.value, key) {
				return true
			}
		}
	}
	return false
}

func (set *StringSet) Add(key []byte) bool {
	cp := set.size + (set.size >> 1)
	if set.capacity <= cp {
		set.increaseCapacity(cp)
	}

	hashCode := internalGetHashCode(key)
	index := hashCode % set.capacity
	var curSlot *Slot
	for i := set.buckets[index] - 1; i >= 0; i = curSlot.next {
		curSlot = &set.slots[i]
		if curSlot.hashCode == hashCode && StringEquals(curSlot.value, key) {
			return false
		}
	}

	curSlot = &set.slots[set.size]
	curSlot.hashCode = hashCode
	curSlot.value = key
	curSlot.next = set.buckets[index] - 1
	set.size++
	set.buckets[index] = set.size
	return true
}
