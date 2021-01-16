package format

import (
	"github.com/olekukonko/tablewriter"
	"io"
)

func ShowSetting(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"socks端口", "http端口", "udp转发", "启用流量监听", "多路复用", "允许局域网连接", "绕过局域网和大陆", "路由策略"})
	center := tablewriter.ALIGN_CENTER
	table.SetAlignment(center)
	table.AppendBulk(datas)
	table.Render()
}

func ShowDNS(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"DNS端口", "境外DNS", "境内DNS", "备用DNS"})
	center := tablewriter.ALIGN_CENTER
	table.SetAlignment(center)
	table.AppendBulk(datas)
	table.Render()
}

func ShowSimpleNode(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	center := tablewriter.ALIGN_CENTER
	left := tablewriter.ALIGN_LEFT
	if len(datas) > 0 && len(datas[0]) == 7 {
		table.SetHeader([]string{"索引", "tcp", "协议", "别名", "地址", "端口", "测试结果"})
		table.SetColumnAlignment([]int{center, center, center, left, center, center, center})
	} else {
		table.SetHeader([]string{"索引", "协议", "别名", "地址", "端口", "测试结果"})
		table.SetColumnAlignment([]int{center, center, left, center, center, center})
	}
	table.SetColWidth(70)
	table.AppendBulk(datas)
	table.Render()
}

func ShowSub(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"索引", "别名", "url", "是否启用"})

	center := tablewriter.ALIGN_CENTER
	table.SetColumnAlignment([]int{center, center, center, center})

	table.AppendBulk(datas)
	table.Render()
}

func ShowRouter(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"索引", "类型", "规则"})
	center := tablewriter.ALIGN_CENTER
	table.SetAlignment(center)
	table.AppendBulk(datas)
	table.Render()
}
