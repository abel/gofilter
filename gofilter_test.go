package gofilter

import (
	//"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringSet(t *testing.T) {
	set := StringSet{}
	set.Add([]byte("abc"))
	set.Add([]byte("efg"))

	b1 := set.Contains([]byte("abc"))
	b2 := set.Contains([]byte("xd"))

	Convey("random attribute count should equal count", func() {
		So(b1, ShouldEqual, true)
		So(b2, ShouldEqual, false)
	})
}

func TestTrieFilter(t *testing.T) {
	filter := TrieFilter{}
	filter.SetFilter(true, true)

	filter.LoadMaskFile("maskWord.txt")
	filter.AddKey([]byte("abc"))
	filter.AddKey([]byte("efg"))

	i, r := filter.Replace([]byte("zzeeabcdefgeffgabc"), byte('*'))
	allr := filter.FindAll([]byte("zzeeabcdefgeffgabc"))
	Convey("Subject： TestSelectRandomAttribute", t, func() {
		Convey("random attribute count should equal count", func() {
			So(i, ShouldEqual, 3)
			So(string(r), ShouldEqual, "zzee*d*effg*")
			So(len(allr), ShouldEqual, 3)
		})
	})
}

func TestTrieFilterLoad(t *testing.T) {
	filter := TrieFilter{}
	filter.SetFilter(true, true)

	filter.LoadMaskFile("maskWord.txt")

	i, r := filter.Replace([]byte("zzeeab毛泽东cdef占领中环geffgabc"), byte('*'))
	Convey("Subject： TestSelectRandomAttribute", t, func() {
		Convey("random attribute count should equal count", func() {
			So(i, ShouldEqual, 2)
			So(string(r), ShouldEqual, "zzeeab*cdef*geffgabc")
		})
	})
}

func TestTrieFilterFile(t *testing.T) {
	LoadMaskWordFile("maskWord.txt")
	LoadMaskNameFile("maskSpecial.txt")

	t1 := "Subject： TestSelectRandomAttribute"
	p1 := TrieReplaceBadWord(t1)

	t2 := "zzeeab毛泽东cdef占领中环geffgabc"
	p2 := TrieReplaceBadWord(t2)

	Convey("Subject： TestSelectRandomAttribute", t, func() {
		Convey("random attribute count should equal count", func() {
			So(t1, ShouldEqual, p1)
			So(p2, ShouldEqual, "zzeeab*cdef*geffgabc")
		})
	})
}
