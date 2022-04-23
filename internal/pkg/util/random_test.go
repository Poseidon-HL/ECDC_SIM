package util

import (
	"testing"
)

func TestSample(t *testing.T) {
	type args struct {
		sample []int
		num    int
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "testSampleList", args: args{
			sample: []int{1, 2, 3, 4, 5},
			num:    3,
		}},
		{name: "testSampleList", args: args{
			sample: []int{1, 2, 3, 4, 5},
			num:    5,
		}},
		{name: "testSampleList", args: args{
			sample: []int{1, 2, 3, 4, 5},
			num:    1,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, item := range Sample(tt.args.sample, tt.args.num) {
				t.Log(item)
			}
		})
	}
}
