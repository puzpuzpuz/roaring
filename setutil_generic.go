// +build !arm64 gccgo appengine

package roaring

import "unsafe"

var (
	k1table = [3]uint{1, 1, 0}
	k2table = [3]uint{1, 0, 1}
)

func compareuint16(x, y uint16) int {
	// branchless analog of the following code:
	// if x == y {
	// 	return 0
	// }
	// if x < y {
	// 	return 1
	// }
	// return 2
	bsmaller := (int(x) - int(y)) >> 63
	bgreater := (int(y) - int(x)) >> 63
	return int(bgreater*(-2) + bsmaller*(-1))
}

func compareuint(x, y uint) int {
	bsmaller := (int(x) - int(y)) >> 63
	bgreater := (int(y) - int(x)) >> 63
	return int(bgreater*(-2) + bsmaller*(-1))
}

func uint16SlicePtr(slice []uint16, idx uint) *uint16 {
	p := unsafe.Pointer(&slice[0])
	indexp := (unsafe.Pointer)(uintptr(p) + 2*uintptr(idx))
	return (*uint16)(indexp)
}

func uint16TableItem(arr [3]uint16, idx int) uint16 {
	p := unsafe.Pointer(&arr[0])
	indexp := (unsafe.Pointer)(uintptr(p) + 2*uintptr(idx))
	return *(*uint16)(indexp)
}

func uintTableItem(arr [3]uint, idx int) uint {
	p := unsafe.Pointer(&arr[0])
	indexp := (unsafe.Pointer)(uintptr(p) + 8*uintptr(idx))
	return *(*uint)(indexp)
}

func union2by2_branchless(set1 []uint16, set2 []uint16, buffer []uint16) int {
	if 0 == len(set2) {
		buffer = buffer[:len(set1)]
		copy(buffer, set1[:])
		return len(set1)
	}
	if 0 == len(set1) {
		buffer = buffer[:len(set2)]
		copy(buffer, set2[:])
		return len(set2)
	}
	var (
		s1, s2     uint16
		nooverflow int = 1
		stable     [3]uint16
	)
	pos := uint(0)
	k1 := uint(0)
	k2 := uint(0)
	len1 := uint(len(set1))
	len2 := uint(len(set2))
	buffer = buffer[:len1+len2]
	for nooverflow == 1 {
		s1 = *uint16SlicePtr(set1, k1)
		s2 = *uint16SlicePtr(set2, k2)
		c := compareuint16(s1, s2)

		stable[0] = s1
		stable[1] = s1
		stable[2] = s2
		*uint16SlicePtr(buffer, pos) = uint16TableItem(stable, c)

		pos++
		k1 += uintTableItem(k1table, c)
		k2 += uintTableItem(k2table, c)

		nooverflow = compareuint(k1, len1) & compareuint(k2, len2)
	}
	if k1 >= len1 {
		copy(buffer[pos:], set2[k2:])
		pos += len2 - k2
		return int(pos)
	}
	if k2 >= len2 {
		copy(buffer[pos:], set1[k1:])
		pos += len1 - k1
		return int(pos)
	}
	return int(pos)
}

func union2by2(set1 []uint16, set2 []uint16, buffer []uint16) int {
	pos := 0
	k1 := 0
	k2 := 0
	if 0 == len(set2) {
		buffer = buffer[:len(set1)]
		copy(buffer, set1[:])
		return len(set1)
	}
	if 0 == len(set1) {
		buffer = buffer[:len(set2)]
		copy(buffer, set2[:])
		return len(set2)
	}
	s1 := set1[k1]
	s2 := set2[k2]
	buffer = buffer[:cap(buffer)]
	for {
		if s1 < s2 {
			buffer[pos] = s1
			pos++
			k1++
			if k1 >= len(set1) {
				copy(buffer[pos:], set2[k2:])
				pos += len(set2) - k2
				break
			}
			s1 = set1[k1]
		} else if s1 == s2 {
			buffer[pos] = s1
			pos++
			k1++
			k2++
			if k1 >= len(set1) {
				copy(buffer[pos:], set2[k2:])
				pos += len(set2) - k2
				break
			}
			if k2 >= len(set2) {
				copy(buffer[pos:], set1[k1:])
				pos += len(set1) - k1
				break
			}
			s1 = set1[k1]
			s2 = set2[k2]
		} else { // if (set1[k1]>set2[k2])
			buffer[pos] = s2
			pos++
			k2++
			if k2 >= len(set2) {
				copy(buffer[pos:], set1[k1:])
				pos += len(set1) - k1
				break
			}
			s2 = set2[k2]
		}
	}
	return pos
}
