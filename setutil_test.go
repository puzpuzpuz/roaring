package roaring

// to run just these tests: go test -run TestSetUtil*

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUtilDifference(t *testing.T) {
	data1 := []uint16{0, 1, 2, 3, 4, 9}
	data2 := []uint16{2, 3, 4, 5, 8, 9, 11}
	result := make([]uint16, 0, len(data1)+len(data2))
	expectedresult := []uint16{0, 1}
	nl := difference(data1, data2, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	expectedresult = []uint16{5, 8, 11}
	nl = difference(data2, data1, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	// empty set2
	data2 = []uint16{}
	expectedresult = []uint16{0, 1, 2, 3, 4, 9}
	nl = difference(data1, data2, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	// empty set 1
	data1 = []uint16{}
	data2 = []uint16{2, 3, 4, 5, 8, 9, 11}
	expectedresult = []uint16{}
	nl = difference(data1, data2, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)
}

func TestCompareuint16(t *testing.T) {
	assert.Equal(t, 0, compareuint16(42, 42))
	assert.Equal(t, 0, compareuint16(42, 1))
	assert.Equal(t, -1, compareuint16(1, 42))
}

func TestUint16SlicePtr(t *testing.T) {
	slice := []uint16{42, 41, 1, 2, 3}
	for i := range slice {
		assert.Equal(t, slice[i], *uint16SlicePtr(slice, uint(i)))
	}
}

func TestSetUtilUnion(t *testing.T) {
	data1 := []uint16{0, 1, 2, 3, 4, 9}
	data2 := []uint16{2, 3, 4, 5, 8, 9, 11}
	result := make([]uint16, 0, len(data1)+len(data2))
	expectedresult := []uint16{0, 1, 2, 3, 4, 5, 8, 9, 11}
	nl := union2by2(data1, data2, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	nl = union2by2(data2, data1, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)
}

func TestSetUtilExclusiveUnion(t *testing.T) {
	data1 := []uint16{0, 1, 2, 3, 4, 9}
	data2 := []uint16{2, 3, 4, 5, 8, 9, 11}
	result := make([]uint16, 0, len(data1)+len(data2))
	expectedresult := []uint16{0, 1, 5, 8, 11}
	nl := exclusiveUnion2by2(data1, data2, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	nl = exclusiveUnion2by2(data2, data1, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)
}

func TestSetUtilIntersection(t *testing.T) {
	data1 := []uint16{0, 1, 2, 3, 4, 9}
	data2 := []uint16{2, 3, 4, 5, 8, 9, 11}
	result := make([]uint16, 0, len(data1)+len(data2))
	expectedresult := []uint16{2, 3, 4, 9}
	nl := intersection2by2(data1, data2, result)
	result = result[:nl]
	result = result[:len(expectedresult)]

	assert.Equal(t, expectedresult, result)

	nl = intersection2by2(data2, data1, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	data1 = []uint16{4}
	data2 = make([]uint16, 10000)

	for i := range data2 {
		data2[i] = uint16(i)
	}

	result = make([]uint16, 0, len(data1)+len(data2))
	expectedresult = data1
	nl = intersection2by2(data1, data2, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)

	nl = intersection2by2(data2, data1, result)
	result = result[:nl]

	assert.Equal(t, expectedresult, result)
}

func TestSetUtilIntersection2(t *testing.T) {
	data1 := []uint16{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}
	data2 := []uint16{0, 3, 6, 9, 12, 15, 18}
	result := make([]uint16, 0, len(data1)+len(data2))
	expectedresult := []uint16{0, 6, 12, 18}
	nl := intersection2by2(data1, data2, result)
	result = result[:nl]
	result = result[:len(expectedresult)]

	assert.Equal(t, expectedresult, result)
}

func TestSetUtilBinarySearch(t *testing.T) {
	data := make([]uint16, 256)
	for i := range data {
		data[i] = uint16(2 * i)
	}
	for i := 0; i < 2*len(data); i++ {
		key := uint16(i)
		loc := binarySearch(data, key)
		if (key & 1) == 0 {
			assert.Equal(t, int(key)/2, loc)
		} else {
			assert.Equal(t, -int(key)/2-2, loc)
		}
	}
}

// go test -bench BenchmarkUnion2by2 -run -
func BenchmarkUnion2by2(b *testing.B) {
	r := rand.New(rand.NewSource(123456))

	// this is important: we pre-generate a large amount of randomized
	// sorted arrays in order to disable the effects branch prediction,
	// making benchmarks against non-branchless implementations
	// more realistic.

	sarrsnum := 1024
	sz1 := 1024
	sarrs := make([][]uint16, sarrsnum)
	for i := 0; i < sarrsnum; i++ {
		sarrs[i] = make([]uint16, sz1)
		for j := 0; j < sz1; j++ {
			sarrs[i][j] = uint16(r.Intn(MaxUint16))
		}
		sort.Sort(uint16Slice(sarrs[i]))
	}

	sz2 := 1024
	s2 := make([]uint16, sz2)

	sz3 := 1024
	s3 := make([]uint16, sz3)

	sz4 := 1024
	s4 := make([]uint16, sz4)

	// We are going to populate our arrays with random data.
	// Importantly, we need to sort. There might be a few
	// duplicates, by random chance, but it should not affect
	// results too much.

	for i := 0; i < sz2; i++ {
		s2[i] = uint16(r.Intn(MaxUint16))
	}
	sort.Sort(uint16Slice(s2))

	for i := 0; i < sz3; i++ {
		s3[i] = uint16(r.Intn(MaxUint16))
	}
	sort.Sort(uint16Slice(s3))

	for i := 0; i < sz4; i++ {
		s4[i] = uint16(r.Intn(MaxUint16))
	}
	sort.Sort(uint16Slice(s4))

	buf := make([]uint16, sz1+sz2+sz3+sz4)

	b.Run("union2by2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < sarrsnum; i++ {
				union2by2(sarrs[i], s2, buf)
				union2by2(sarrs[i], s3, buf)
				union2by2(sarrs[i], s4, buf)
			}
		}
	})

	// the old, non-branchless implementation for performance
	// comparison can be found here:
	// https://github.com/RoaringBitmap/roaring/blob/ff33c3b226c3ac033bf1a0b0f3ed647fc9cd2efa/setutil_generic.go
}
