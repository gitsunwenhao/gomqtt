package gate

import (
	"testing"

	"github.com/coreos/etcd/clientv3"
)

func Test_loadConfig(t *testing.T) {
	type args struct {
		staticConf bool
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadConfig(tt.args.staticConf)
		})
	}
}

func Test_watchEtcd(t *testing.T) {
	type args struct {
		cli *clientv3.Client
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watchEtcd(tt.args.cli)
		})
	}
}

func Test_uploadEtcd(t *testing.T) {
	type args struct {
		cli *clientv3.Client
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uploadEtcd(tt.args.cli)
		})
	}
}

func Test_getHost(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHost(); got != tt.want {
				t.Errorf("getHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
