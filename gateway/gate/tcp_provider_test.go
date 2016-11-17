package gate

import (
	"net"
	"reflect"
	"testing"
)

func TestTcpProvider_Start(t *testing.T) {
	tests := []struct {
		name string
		tp   *TcpProvider
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TcpProvider{}
			tp.Start()
		})
	}
}

func TestTcpProvider_Close(t *testing.T) {
	tests := []struct {
		name    string
		tp      *TcpProvider
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &TcpProvider{}
			if err := tp.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TcpProvider.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accept(t *testing.T) {
	type args struct {
		ln net.Listener
	}
	tests := []struct {
		name    string
		args    args
		want    net.Conn
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := accept(tt.args.ln)
			if (err != nil) != tt.wantErr {
				t.Errorf("accept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("accept() = %v, want %v", got, tt.want)
			}
		})
	}
}
