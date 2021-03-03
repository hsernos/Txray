package node

import (
	"Txray/core/config"
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
)

type node struct {
	data       protocols.Protocol
	Link       string `json:"link"`
	Subid      string `json:"subid"`
	TestResult string `json:"testResult"`
}

type subscribe struct {
	Remarks string `json:"remarks"`
	URL     string `json:"url"`
	Using   bool   `json:"using"`
	ID      string `json:"id"`
}

type NodesAndSub struct {
	Index int          `json:"index"`
	Subs  []*subscribe `json:"subscribe"`
	Nodes []*node      `json:"nodes"`
}

var data *NodesAndSub = initNodeData()

func initNodeData() *NodesAndSub {
	nas := new(NodesAndSub)
	file := config.NodesAndSub
	if tools.IsFile(file) {
		err := tools.ReadJSON(file, nas)
		if err != nil {
			log.Error(err)
			return nil
		}
		for _, n := range nas.Nodes {
			n.data = protocols.ParseLink(n.Link)
		}
	} else {
		nas.Index = 1
	}
	return nas
}

func (r *NodesAndSub) save() {
	file := config.NodesAndSub
	err := tools.WriteJSON(*r, file)
	if err != nil {
		log.Error(err)
	}
}
