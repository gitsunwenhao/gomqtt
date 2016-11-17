package gate

import (
	"testing"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
)

func Test_subscribe(t *testing.T) {
	type args struct {
		ci *connInfo
		p  *proto.SubscribePacket
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
			if err := subscribe(tt.args.ci, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_unsubscribe(t *testing.T) {
	type args struct {
		ci *connInfo
		p  *proto.UnsubscribePacket
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
			if err := unsubscribe(tt.args.ci, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("unsubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subToStream(t *testing.T) {
	type args struct {
		t   []byte
		qos byte
	}
	tests := []struct {
		name    string
		args    args
		want    byte
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := subToStream(tt.args.t, tt.args.qos)
			if (err != nil) != tt.wantErr {
				t.Errorf("subToStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("subToStream() = %v, want %v", got, tt.want)
			}
		})
	}
}
