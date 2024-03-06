package lib

import (
	"testing"
)

func TestFormattedNumber(t *testing.T) {
	type args struct {
		num int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"positive", args{123}, "12.3"},
		{"negative", args{-123}, "-12.3"},
		{"zero", args{0}, "0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormattedNumber(tt.args.num); got != tt.want {
				t.Errorf("FormattedNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormattedAvg(t *testing.T) {
	type args struct {
		sum   int64
		count int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"positive", args{123, 3}, "4.1"},
		{"negative", args{-123, 3}, "-4.1"},
		{"small_neg", args{-1, 10}, "0.0"},
		{"small_pos", args{1, 10}, "0.1"},
		{"small_neg_with_non_zero_res", args{-16, 10}, "-0.1"},
		{"small_pos_with_greater_then_1_res", args{116, 10}, "1.2"},
		{"zero", args{0, 3}, "0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormattedAvg(tt.args.sum, tt.args.count); got != tt.want {
				t.Errorf("FormattedAvg() = %v, want %v", got, tt.want)
			}
		})
	}
}
