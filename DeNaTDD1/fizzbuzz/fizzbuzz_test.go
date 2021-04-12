package fizzbuzz_test

import (
	"fmt"
	"testing"

	"github.com/Densuke-fitness/PracticeRepo3/fizzbuzz"
)

func TestConvert(t *testing.T) {

	tests := []struct {
		n    int
		want string
	}{
		{n: 1, want: "1"},
		{n: 2, want: "2"},
		{n: 3, want: "Fizz"},
		{n: 4, want: "4"},
		{n: 5, want: "Buzz"},
		{n: 6, want: "Fizz"},
		{n: 7, want: "7"},
		{n: 8, want: "8"},
		{n: 9, want: "Fizz"},
		{n: 10, want: "Buzz"},
		{n: 11, want: "11"},
		{n: 12, want: "Fizz"},
		{n: 13, want: "13"},
		{n: 14, want: "14"},
		{n: 15, want: "FizzBuzz"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("number:%v", tt.n)

		t.Run(name, func(t *testing.T) {
			got := fizzbuzz.Convert(tt.n)
			if got != tt.want {
				t.Errorf(`Convert(%v) = %q but want %q`, tt.n, got, tt.want)
			}

		})

	}

}
