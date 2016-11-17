package gate

import (
	"reflect"
	"testing"
)

func Test_saveCI(t *testing.T) {
	type args struct {
		ci *connInfo
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveCI(tt.args.ci)
		})
	}
}

func Test_getCI(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want *connInfo
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCI(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_delCI(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delCI(tt.args.id)
		})
	}
}
