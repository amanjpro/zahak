package search

import (
	"fmt"
	"testing"

	. "github.com/amanjpro/zahak/engine"
)

func TestMovepickerNext(t *testing.T) {
	mp := &MovePicker{
		nil,
		nil,
		10,
		[]Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
		[]int32{1000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
		0,
		0,
		false,
	}

	for i := 0; i < mp.Length(); i++ {
		actual := mp.Next()
		expected := mp.Length() - i
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("Expected %d But got %d\n", expected, int32(actual)))
		}
	}
}
