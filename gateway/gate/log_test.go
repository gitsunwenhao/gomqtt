package gate

import "testing"

func TestInitLogger(t *testing.T) {
	type args struct {
		lp      string
		lv      string
		isDebug bool
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitLogger(tt.args.lp, tt.args.lv, tt.args.isDebug)
		})
	}
}
