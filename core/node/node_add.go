package node

import (
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
	"fmt"
)

// 从订阅文本更新节点
func AddNodeBySubText(subtext string) {
	links := protocols.Sub2links(subtext)
	AddNodeByLinks(links...)
}

// 根据链接添加节点
func AddNodeByLinks(links ...string) {
	defer data.save()
	var i = 0
	for _, link := range links {
		nodeData := protocols.ParseLink(link)
		if nodeData != nil {
			n := new(node)
			n.data = nodeData
			n.Link = link
			data.Nodes = append(data.Nodes, n)
			i++
		} else {
			log.Warnf("添加失败: [%s]", link)
		}
	}
	log.Infof("成功更新 [%d] 个节点, 失败 [%d] 个", i, len(links)-i)
}

// 从本地订阅文件批量添加节点
func AddNodeBySubFile(absPath string) {
	if tools.IsFile(absPath) {
		subText := tools.ReadFile(absPath)[0]
		log.Info("解析文件中...")
		links := protocols.Sub2links(subText)
		log.Info("解析文件完成，解析VMess链接如下: ")
		log.Info("=======================================================================")
		for index, x := range links {
			log.Info(fmt.Sprintf("%3d. ", index+1), x)
		}
		log.Info("=======================================================================")
		AddNodeByLinks(links...)
	} else {
		log.Warn("该文件不存在")
	}
}

// 从本地链接文件批量添加节点
func AddNodeByFile(absPath string) {
	if tools.IsFile(absPath) {
		log.Info("读取文件中...")
		links := tools.ReadFile(absPath)
		log.Info("读取文件完成，解析链接如下")
		for index, x := range links {
			log.Info(fmt.Sprintf("%3d. ", index+1), x)
		}
		AddNodeByLinks(links...)
	} else {
		log.Warn("该文件不存在")
	}
}

// 添加一个VMess节点
func AddVMessNode(remarks, addr, port, id, aid, network, types, host, path, tls string) {
	if !tools.IsNetPort(port) {
		log.Errorf("%s: 不是一个端口", port)
		return
	}
	if !tools.IsInt(aid) {
		log.Errorf("%s: 不是一个数字", aid)
		return
	}
	defer data.save()
	// 初始化VMess节点数据
	vmessData := new(protocols.VMess)
	vmessData.V = "2"
	vmessData.Ps = remarks
	vmessData.Add = addr
	vmessData.Port = port
	vmessData.Id = id
	vmessData.Aid = aid
	vmessData.Net = network
	vmessData.Type = types
	vmessData.Host = host
	vmessData.Path = path
	vmessData.Tls = tls
	// 生成节点
	n := new(node)
	n.data = vmessData
	n.Link = vmessData.GetLink()
	// 添加到集合中
	data.Nodes = append(data.Nodes, n)
}

// 添加一个VLess节点
func AddVLESSNode() {
	//TODO
}

// 添加一个Trojan节点
func AddTrojanNode(remarks, addr, port, password string) {
	if !tools.IsNetPort(port) {
		log.Errorf("%s: 不是一个端口", port)
		return
	}
	defer data.save()
	// 初始化VMess节点数据
	trojanData := new(protocols.Trojan)
	trojanData.Remarks = remarks
	trojanData.Address = addr
	trojanData.Port = port
	trojanData.Password = password
	// 生成节点
	n := new(node)
	n.data = trojanData
	n.Link = trojanData.GetLink()
	// 添加到集合中
	data.Nodes = append(data.Nodes, n)
}

// 添加一个ShadowSocks节点
func AddShadowSocksNode(remarks, addr, port, password, method string) {
	if !tools.IsNetPort(port) {
		log.Errorf("%s: 不是一个端口", port)
		return
	}
	defer data.save()
	ss := new(protocols.ShadowSocks)
	ss.Remarks = remarks
	ss.Address = addr
	ss.Port = port
	ss.Password = password
	ss.Method = method
	// 生成节点
	n := new(node)
	n.data = ss
	n.Link = ss.GetLink()
	// 添加到集合中
	data.Nodes = append(data.Nodes, n)
}
