package gate

import "testing"

func Test_userValidate(t *testing.T) {
	type args struct {
		u []byte
		p []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userValidate(tt.args.u, tt.args.p); got != tt.want {
				t.Errorf("userValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
