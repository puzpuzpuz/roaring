// +build !arm64 gccgo appengine

package roaring

func union2by2(set1 []uint16, set2 []uint16, buffer []uint16) int {
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
	var s1, s2 uint16
	pos := uint(0)
	k1 := uint(0)
	k2 := uint(0)
	len1 := uint(len(set1))
	len2 := uint(len(set2))
	buffer = buffer[:len1+len2]
	for k1 < len1 && k2 < len2 {
		s1 = *uint16SlicePtr(set1, k1)
		s2 = *uint16SlicePtr(set2, k2)

		sflag := compareuint16(s1, s2) // -1 if s1 < s2, zero otherwise
		lflag := compareuint16(s2, s1) // -1 if s2 < s1, zero otherwise
		*uint16SlicePtr(buffer, pos) = uint16(-sflag)*s1 + uint16(1+sflag)*s2

		pos++
		k1 += uint(1 + lflag)
		k2 += uint(1 + sflag)
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
