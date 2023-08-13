package list

import (
	"reflect"
	"sort"
	"testing"
)

func TestInList(t *testing.T) {
	type args[T comparable] struct {
		needle T
		list   []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want bool
	}
	tests := []testCase[int64]{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args[int64]{
				needle: 1,
				list:   []int64{1, 2, 3, 4},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InList(tt.args.needle, tt.args.list); got != tt.want {
				t.Errorf("InList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	type args[T comparable] struct {
		l1 []T
		l2 []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int64]{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args[int64]{
				l1: []int64{1, 4, 7, 8, 0, 10},
				l2: []int64{2, 4, 5, 6, 8, 9, 10, 4},
			},
			want: []int64{4, 8, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersection(tt.args.l1, tt.args.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnionList(t *testing.T) {
	type args[T comparable] struct {
		l1 []T
		l2 []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args[int]{
				l1: []int{1, 3, 5, 7, 9},
				l2: []int{2, 4, 6, 8, 10},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UnionList(tt.args.l1, tt.args.l2)
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnionList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueList(t *testing.T) {
	type args[T comparable] struct {
		s []T
	}
	type testCase[T comparable] struct {
		name    string
		args    args[T]
		wantRes []T
	}
	tests := []testCase[int64]{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args[int64]{
				s: []int64{1, 2, 3, 4, 5, 2, 4, 23, 124, 5, 1},
			},
			wantRes: []int64{1, 2, 3, 4, 5, 23, 124},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := UniqueList(tt.args.s); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("UniqueList() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDifferenceList(t *testing.T) {
	type args[T comparable] struct {
		l1 []T
		l2 []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int64]{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args[int64]{
				l1: []int64{1, 3, 4, 5},
				l2: []int64{2, 1, 3, 4},
			},
			want: []int64{5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DifferenceList(tt.args.l1, tt.args.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DifferenceList() = %v, want %v", got, tt.want)
			}
		})
	}
}
