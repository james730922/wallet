package condition

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func Test_structToMap(t *testing.T) {
	type item struct {
		IDPtr   *int   `json:"id_ptr"`
		IDVal   int    `json:"id_val"`
		ListPtr *[]int `json:"list_ptr"`
		ListVal []int  `json:"list_val"`
	}
	caseItem := &item{
		IDPtr:   aws.Int(1111),
		IDVal:   22222,
		ListPtr: &[]int{2, 4, 6, 8},
		ListVal: []int{1, 3, 5, 7},
	}
	var caseEmptyListVal []int

	type args struct {
		item interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "ptr not nil",
			args: args{
				item: caseItem,
			},
			want: map[string]interface{}{
				"id_ptr":   caseItem.IDPtr,
				"id_val":   caseItem.IDVal,
				"list_ptr": caseItem.ListPtr,
				"list_val": caseItem.ListVal,
			},
		},
		{
			name: "ptr is nil",
			args: args{
				&item{
					IDPtr:   nil,
					IDVal:   caseItem.IDVal,
					ListPtr: nil,
					ListVal: caseItem.ListVal,
				},
			},
			want: map[string]interface{}{
				"id_val":   caseItem.IDVal,
				"list_val": caseItem.ListVal,
			},
		},
		{
			name: "list is nil",
			args: args{
				&item{
					IDPtr:   nil,
					IDVal:   caseItem.IDVal,
					ListPtr: nil,
					ListVal: caseEmptyListVal,
				},
			},
			want: map[string]interface{}{
				"id_val":   caseItem.IDVal,
				"list_val": caseEmptyListVal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := structToMap(tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("structToMap() = %v, want %v", got, tt.want)
			} else {
				t.Log(got)
			}
		})
	}
}
