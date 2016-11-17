package gate

import (
	"testing"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
)

func Test_publish(t *testing.T) {
	type args struct {
		ci *connInfo
		p  *proto.PublishPacket
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
			if err := publish(tt.args.ci, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
