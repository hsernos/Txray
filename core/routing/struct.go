package routing

import (
	"Txray/core"
	"Txray/log"
	"encoding/json"
	"os"
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

var route *Routing = NewRouting()

func NewRouting() *Routing {
	return &Routing{
		Proxy:  make([]*routing, 0),
		Direct: make([]*routing, 0),
		Block:  make([]*routing, 0),
	}
}

func init() {
	if _, err := os.Stat(core.RoutingFile); os.IsNotExist(err) {
		route.save()
	} else {
		file, _ := os.Open(core.RoutingFile)
		defer file.Close()
		err := json.NewDecoder(file).Decode(route)
		if err != nil {
			log.Error(err)
		}
	}
}

func (r *Routing) save() {
	err := core.WriteJSON(r, core.RoutingFile)
	if err != nil {
		log.Error(err)
	}
}
