package filtertest

import (
	"os"
	"time"
)

type implementsAll struct{}

func (implementsAll) Read([]byte) (int, error) { return 0, nil }
func (implementsAll) String() string           { return "" }

func detectType() {
	{
		type withNamedTime struct {
			x int
			y time.Time
		}
		var foo struct {
			x time.Time
		}
		var bar withNamedTime
		typeTest(withNamedTime{}, "contains time.Time") // want `YES`
		typeTest(foo, "contains time.Time")             // want `YES`
		typeTest(bar, "contains time.Time")             // want `YES`
	}

	{
		type timeFirst struct {
			y time.Time
			x int
		}
		var foo struct {
			x time.Time
		}
		var bar timeFirst
		typeTest(timeFirst{}, "starts with time.Time") // want `YES`
		typeTest(foo, "starts with time.Time")         // want `YES`
		typeTest(bar, "starts with time.Time")         // want `YES`
	}

	{
		type intPair struct {
			y int
			x int
		}
		var foo struct {
			x float32
			y float32
		}
		var bar intPair
		typeTest(struct { // want `YES`
			_ string
			_ string
		}{}, "non-underlying type test; T + T")
		typeTest(intPair{}, "non-underlying type test; T + T") // type is Named, not struct
		typeTest(foo, "non-underlying type test; T + T")       // want `YES`
		typeTest(bar, "non-underlying type test; T + T")       // type is Named, not struct
	}

	var i1, i2 int
	var ii []int
	var s1, s2 string
	var ss []string
	typeTest(s1 + s2) // want `concat`
	typeTest(i1 + i2) // want `addition`
	typeTest(s1 > s2) // want `s1 !is\(int\)`
	typeTest(i1 > i2) // want `i1 !is\(string\) && pure`
	typeTest(random() > i2)
	typeTest(ss, ss) // want `ss is\(\[\]string\)`
	typeTest(ii, ii)
	typeTest("2 type filters", i1)
	typeTest("2 type filters", s1)
	typeTest("2 type filters", ii) // want `ii !is\(string\) && !is\(int\)`

	typeTest(implementsAll{}, "implements io.Reader") // want `YES`
	typeTest(i1, "implements io.Reader")
	typeTest(ss, "implements io.Reader")
	typeTest(implementsAll{}, "implements foolib.Stringer") // want `YES`
	typeTest(i1, "implements foolib.Stringer")
	typeTest(ss, "implements foolib.Stringer")

	typeTest([100]byte{}, "size>=100") // want `YES`
	typeTest([105]byte{}, "size>=100") // want `YES`
	typeTest([10]byte{}, "size>=100")
	typeTest([100]byte{}, "size<=100") // want `YES`
	typeTest([105]byte{}, "size<=100")
	typeTest([10]byte{}, "size<=100") // want `YES`
	typeTest([100]byte{}, "size>100")
	typeTest([105]byte{}, "size>100") // want `YES`
	typeTest([10]byte{}, "size>100")
	typeTest([100]byte{}, "size<100")
	typeTest([105]byte{}, "size<100")
	typeTest([10]byte{}, "size<100")   // want `YES`
	typeTest([100]byte{}, "size==100") // want `YES`
	typeTest([105]byte{}, "size==100")
	typeTest([10]byte{}, "size==100")
	typeTest([100]byte{}, "size!=100")
	typeTest([105]byte{}, "size!=100") // want `YES`
	typeTest([10]byte{}, "size!=100")  // want `YES`

	var time1, time2 time.Time
	var err error
	typeTest(time1 == time2, "time==time") // want `YES`
	typeTest(err == nil, "time==time")
	typeTest(nil == err, "time==time")
	typeTest(time1 != time2, "time!=time") // want `YES`
	typeTest(err != nil, "time!=time")
	typeTest(nil != err, "time!=time")

	intFunc := func() int { return 10 }
	intToIntFunc := func(x int) int { return x }
	typeTest(intFunc(), "func() int")                 // want `YES`
	typeTest(func() int { return 0 }(), "func() int") // want `YES`
	typeTest(func() string { return "" }(), "func() int")
	typeTest(intToIntFunc(1), "func() int")

	typeTest(intToIntFunc(2), "func(int) int") // want `YES`
	typeTest(intToIntFunc, "func(int) int")
	typeTest(intFunc, "func(int) int")

	var v implementsAll
	typeTest(v.String(), "func() string") // want `YES`
	typeTest(implementsAll.String(v), "func() string")
	typeTest(implementsAll.String, "func() string")

}

func detectPure(x int) {
	pureTest(random()) // want `!pure`
	pureTest(x * x)    // want `pure`
}

func detectText(foo, bar int) {
	textTest(foo, "text=foo") // want `YES`
	textTest(bar, "text=foo")

	textTest("foo", "text='foo'") // want `YES`
	textTest("bar", "text='foo'")

	textTest("bar", "text!='foo'") // want `YES`
	textTest("foo", "text!='foo'")

	textTest(32, "matches d+") // want `YES`
	textTest(0x32, "matches d+")
	textTest("foo", "matches d+")

	textTest(1, "doesn't match [A-Z]") // want `YES`
	textTest("ABC", "doesn't match [A-Z]")
}

func detectParensFilter() {
	var err error
	parensFilterTest(err, "type is error") // want `YES`
}

func fileFilters1() {
	// No matches as this file doesn't import "path/filepath".
	importsTest(os.PathSeparator, "path/filepath")
	importsTest(os.PathListSeparator, "path/filepath")
}
