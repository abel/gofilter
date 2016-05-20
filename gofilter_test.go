package gofilter

import (
	"fmt"
	"runtime/debug"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type Point struct {
	X int
	Y int
}

func TestPrintStack(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			bin := string(debug.Stack())
			fmt.Println("TestPrintStack")
			fmt.Println(bin)
		}
	}()
	value := 111
	zero := 0
	value = value / zero
}

func TestCopy(t *testing.T) {

	defer fmt.Println("defer1")
	defer fmt.Println("defer2")
	points := []Point{}
	for i := 0; i < 5; i++ {
		defer func() { fmt.Println("i:%v", i) }()
	}

	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			defer fmt.Println("append i:%v", i)
		}
		points = append(points, Point{X: i, Y: i * i})
	}

	{
		defer fmt.Println("defer3")
		defer fmt.Println("defer4")
		newP := Point{X: 6, Y: 6 * 6}
		i := 0
		points = append(points, newP)
		copy(points[i+1:len(points)], points[i:len(points)-1])
		points[i] = newP
		fmt.Println("TestCopy1:%v", points)
	}
	defer fmt.Println("defer5")
	fmt.Println("TestCopy2:%v", points)

}

func TestStringSet(t *testing.T) {
	set := StringSet{}
	set.Add([]byte("abc"))
	set.Add([]byte("efg"))

	b1 := set.Contains([]byte("abc"))
	b2 := set.Contains([]byte("xd"))

	Convey("gofilter test", t, func() {
		Convey("TestStringSet", func() {
			So(b1, ShouldEqual, true)
			So(b2, ShouldEqual, false)
		})
	})
}

func TestTrieFilter(t *testing.T) {
	filter := TrieFilter{}
	filter.SetFilter(true)

	filter.LoadMaskFile("maskWord.txt")
	filter.AddKey("abc")
	filter.AddKey("efg")

	i, r := filter.Replace([]byte("zzeeabcdefgeffgabc"), byte('*'))
	allr := filter.FindAll([]byte("zzeeabcdefgeffgabc"))
	Convey("gofilter test", t, func() {
		Convey("TestTrieFilter", func() {
			So(i, ShouldEqual, 3)
			So(string(r), ShouldEqual, "zzee*d*effg*")
			So(len(allr), ShouldEqual, 3)
		})
	})
}

func TestTrieFilterLoad(t *testing.T) {
	filter := TrieFilter{}
	filter.SetFilter(true)

	filter.LoadMaskFile("maskWord.txt")

	i, r := filter.Replace([]byte("zzeeab毛泽东cdef占领中环geffgabc"), byte('*'))
	Convey("gofilter test", t, func() {
		Convey("TestTrieFilterLoad", func() {
			So(i, ShouldEqual, 2)
			So(string(r), ShouldEqual, "zzeeab*cdef*geffgabc")
		})
	})
}

func TestTrieFilterFile(t *testing.T) {
	LoadMaskWordFile("maskWord.txt")
	LoadMaskNameFile("maskSpecial.txt")

	t1 := "Subject： TestSelectRandomAttribute"
	p1 := ReplaceBadWord(t1)

	t2 := "zzeeab毛泽东cdef占领中环geffgabc"
	p2 := ReplaceBadWord(t2)

	Convey("gofilter test", t, func() {
		Convey("TestTrieFilterFile", func() {
			So(t1, ShouldEqual, p1)
			So(p2, ShouldEqual, "zzeeab*cdef*geffgabc")
		})
	})
}
