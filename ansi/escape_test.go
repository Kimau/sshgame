package ansi

import (
	"testing"
)

func TestAnsFileToStr(t *testing.T) {

}

func TestClamp(t *testing.T) {
	for _, tt := range []struct {
		in  [3]int // input
		out int    // expected result
	}{
		{[3]int{1, 2, 3}, 2},
		{[3]int{1, 2, 2}, 2},
		{[3]int{10, 2, 2}, 2},
		{[3]int{100, 2, 5}, 5},
		{[3]int{4, 2, 5}, 4},
		{[3]int{-100, 2, 5}, 2},
		{[3]int{2, 2, 5}, 2},
	} {
		actual := Clamp(tt.in[0], tt.in[1], tt.in[2])
		if actual != tt.out {
			t.Errorf("Fib(%d): out %d, actual %d", tt.in, tt.out, actual)
		}
	}
}
