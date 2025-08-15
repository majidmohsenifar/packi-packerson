package pack

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCorrect(t *testing.T) {
	tests := []struct {
		ordered int
		want    map[int]int
	}{
		{1, map[int]int{250: 1}},
		{250, map[int]int{250: 1}},
		{500, map[int]int{500: 1}},
		{750, map[int]int{500: 1, 250: 1}},
		{1000, map[int]int{1000: 1}},
		{12001, map[int]int{5000: 2, 2000: 1, 250: 1}},
		{4999, map[int]int{2000: 2, 1000: 1}},
		{5001, map[int]int{5000: 1, 250: 1}},
		{9999, map[int]int{5000: 1, 2000: 2, 1000: 1}},
		{10000, map[int]int{5000: 2}},
		{12500, map[int]int{5000: 2, 2000: 1, 500: 1}},
		{15000, map[int]int{5000: 3}},
		{2250, map[int]int{2000: 1, 250: 1}},
		{3750, map[int]int{2000: 1, 1000: 1, 500: 1, 250: 1}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("Size %d", tt.ordered), func(t *testing.T) {
			got := Correct(tt.ordered)

			if len(got) != len(tt.want) {
				t.Errorf("cuurect(%d) = %v, want %v", tt.ordered, got, tt.want)
				return
			}
			for pack, count := range tt.want {
				if gotCount, exists := got[pack]; !exists || gotCount != count {
					t.Errorf("cuurect(%d)[%d] = %d, want %d", tt.ordered, pack, gotCount, count)
				}
			}
		})
	}
}

func TestInCorrect(t *testing.T) {
	tests := []struct {
		ordered int
		want    []map[int]int
	}{
		{
			ordered: 1,
			want: []map[int]int{
				{5000: 1},
				{2000: 1},
				{1000: 1},
				{500: 1},
			},
		},
		{
			ordered: 251,
			want: []map[int]int{
				{5000: 1},
				{2000: 1},
				{1000: 1},
				{250: 2},
			},
		},
		{
			ordered: 501,
			want: []map[int]int{
				{5000: 1},
				{2000: 1},
				{1000: 1},
				{500: 2},
				{250: 3},
			},
		},
		{
			ordered: 12001,
			want: []map[int]int{
				{5000: 3},
				{2000: 7},
				{1000: 13},
				{500: 25},
				{250: 49},
			},
		},
		{
			ordered: 249,
			want: []map[int]int{
				{5000: 1},
				{2000: 1},
				{1000: 1},
				{500: 1},
			},
		},
		{
			ordered: 1250,
			want: []map[int]int{
				{5000: 1},
				{2000: 1},
				{1000: 2},
				{500: 3},
				{250: 6},
			},
		},
		{
			ordered: 1750,
			want: []map[int]int{
				{5000: 1},
				{2000: 1},
				{1000: 2},
				{500: 4},
				{250: 8},
			},
		},
		{
			ordered: 3000,
			want: []map[int]int{
				{5000: 1},
				{2000: 2},
				{1000: 4},
				{500: 7},
				{250: 13},
			},
		},
		{
			ordered: 4999,
			want: []map[int]int{
				{5000: 1},
				{2000: 3},
				{1000: 5},
				{500: 10},
				{250: 20},
			},
		},
		{
			ordered: 5001,
			want: []map[int]int{
				{5000: 2},
				{2000: 3},
				{1000: 6},
				{500: 11},
				{250: 21},
			},
		},
		{
			ordered: 9999,
			want: []map[int]int{
				{5000: 2},
				{2000: 5},
				{1000: 10},
				{500: 20},
				{250: 40},
			},
		},
		{
			ordered: 10000,
			want: []map[int]int{
				{5000: 3},
				{2000: 6},
				{1000: 11},
				{500: 21},
				{250: 41},
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Size %d", tt.ordered), func(t *testing.T) {
			got := InCorrect(tt.ordered)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InCorrect(%d) = %v, want %v", tt.ordered, got, tt.want)
			}
		})
	}
}

func BenchmarkCorrect(b *testing.B) {
	tests := []int{
		1, 250, 500, 750, 1000, 12001, 4999, 5001, 9999, 10000, 12500, 15000, 2250, 3750,
	}

	for _, size := range tests {
		b.Run(fmt.Sprintf("Size %d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Correct(size)
			}
		})
	}
}

func BenchmarkInCorrect(b *testing.B) {
	tests := []int{
		1, 251, 501, 12001, 249, 1250, 1750, 3000, 4999, 5001, 9999, 10000,
	}

	for _, size := range tests {
		b.Run(fmt.Sprintf("Size %d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				InCorrect(size)
			}
		})
	}
}
