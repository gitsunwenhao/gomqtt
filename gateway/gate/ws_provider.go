package gate

/* Websocket Provider */
type WsProvider struct {
}

func (tp *WsProvider) Start() {
    Logger.Debug("websocket provider startted") 
}

func (tp *WsProvider) Close() error {
	return nil
}
