package gate

/* Websocket Provider */

type TcpProvider struct {
}

func (tp *TcpProvider) Start() {
	Logger.Debug("tcp provider startted") 
}

func (tp *TcpProvider) Close() error {
	return nil
}
