package tools

import (
	"reflect"
	"testing"
)

func TestSliceBatch(t *testing.T) {
	type args struct {
		eachCount int
		slice     interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "have remainder",
			args: args{
				eachCount: 5,
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			},
			want: []interface{}{[]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 10}, []int{11, 12, 13}},
		},
		{
			name: "divide",
			args: args{
				eachCount: 5,
				slice:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: []interface{}{[]int{1, 2, 3, 4, 5}, []int{6, 7, 8, 9, 10}},
		},
		{
			name: "nil",
			args: args{
				eachCount: 5,
				slice:     nil,
			},
			want: []interface{}{},
		},
		{
			name: "empty slice",
			args: args{
				eachCount: 5,
				slice:     []int{},
			},
			want: []interface{}{},
		},
		{
			name: "not slice",
			args: args{
				eachCount: 5,
				slice:     0,
			},
			want: []interface{}{},
		},
		{
			name: "ptr",
			args: args{
				eachCount: 5,
				slice:     &struct{}{},
			},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceBatch(tt.args.eachCount, tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
