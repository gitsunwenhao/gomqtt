package gate

import "fmt"

type Gate struct {
}

func New() *Gate {
	return &Gate{}
}

func (g *Gate) Start(isStatic bool) {
	// init configurations
	fmt.Println(isStatic)

	loadConfig(isStatic)

	// init providers
	providersStart()

	// init addmin service
	go adminStart()

	// start the monitors
	monitorsStart()
}
