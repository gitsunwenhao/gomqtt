package gate

import (
	"testing"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
)

func Test_processPacket(t *testing.T) {
	type args struct {
		ci *connInfo
		pt proto.Packet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := processPacket(tt.args.ci, tt.args.pt); (err != nil) != tt.wantErr {
				t.Errorf("processPacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pingReq(t *testing.T) {
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
			pingReq(tt.args.ci)
		})
	}
}
