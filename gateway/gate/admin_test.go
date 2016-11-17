package gate

import (
	"testing"

	"github.com/labstack/echo"
)

func Test_adminStart(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adminStart()
		})
	}
}

func Test_reload(t *testing.T) {
	type args struct {
		c echo.Context
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
			if err := reload(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("reload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
