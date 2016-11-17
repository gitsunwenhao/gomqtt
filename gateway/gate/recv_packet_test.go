package gate

import (
	"reflect"
	"testing"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
)

func Test_recvPacket(t *testing.T) {
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
			recvPacket(tt.args.ci)
		})
	}
}

func Test_read(t *testing.T) {
	type args struct {
		ci *connInfo
	}
	tests := []struct {
		name  string
		args  args
		want  proto.Packet
		want1 bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := read(tt.args.ci)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("read() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("read() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
