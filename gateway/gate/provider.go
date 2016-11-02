package gate

type Provider interface {
	Start()
	Close() error
}

func providersStart() {

}
