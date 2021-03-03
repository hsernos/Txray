package routing

import (
	"Txray/core/config"
	"Txray/log"
	"Txray/tools"
)

type routing struct {
	Data string `json:"data"`
	Mode Mode   `json:"mode"`
}

type Routing struct {
	Proxy  []*routing `json:"proxy"`
	Direct []*routing `json:"direct"`
	Block  []*routing `json:"block"`
}

var route *Routing = initRouting()

func initRouting() *Routing {
	r := new(Routing)
	file := config.Routing
	if tools.IsFile(file) {
		err := tools.ReadJSON(file, r)
		if err != nil {
			log.Error(err)
			return nil
		}
	}
	return r
}

func (r *Routing) save() {
	file := config.Routing
	err := tools.WriteJSON(*r, file)
	if err != nil {
		log.Error(err)
	}
}
