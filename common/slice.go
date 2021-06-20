package common

import "strings"

// 3, 2, 1, 2, 3, 4, 2, 3, 4, 3, 2, 1
// 0, 1, 2, 3     4, 5, 6, 7     8, 9
func Slice(size, step int) (resp [][]int) {
	var idx []int
	for i := 0; i < size; i++ {
		idx = append(idx, i)
		if len(idx) >= step {
			resp = append(resp, idx)
			idx = nil
		}
	}
	if len(idx) != 0 {
		resp = append(resp, idx)
	}
	return
}

func Smallify(in string, max int) string {
	in = strings.ReplaceAll(in, "\n", "$n$")
	if len(in) > max {
		in = in[:max] + "..."
	}
	return in
}
