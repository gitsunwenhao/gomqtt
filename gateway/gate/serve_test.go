package gate

import (
	"net"
	"testing"
)

func Test_serve(t *testing.T) {
	type args struct {
		c net.Conn
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serve(tt.args.c)
		})
	}
}

func Test_initConnection(t *testing.T) {
	type args struct {
		ci *connInfo
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
			if err := initConnection(tt.args.ci); (err != nil) != tt.wantErr {
				t.Errorf("initConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
