package gate

import "testing"

func TestWsProvider_Start(t *testing.T) {
	tests := []struct {
		name string
		tp   *WsProvider
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &WsProvider{}
			tp.Start()
		})
	}
}

func TestWsProvider_Close(t *testing.T) {
	tests := []struct {
		name    string
		tp      *WsProvider
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := &WsProvider{}
			if err := tp.Close(); (err != nil) != tt.wantErr {
				t.Errorf("WsProvider.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
