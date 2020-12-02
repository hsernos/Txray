package cmd

import (
	"github.com/olekukonko/tablewriter"
	"io"
)

func ShowSetting(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"socks5端口", "http端口", "udp转发", "启用流量监听", "多路复用", "允许局域网连接", "绕过局域网和大陆", "路由策略"})
	center := tablewriter.ALIGN_CENTER
	table.SetAlignment(center)

	//fgGreen := tablewriter.Color(tablewriter.FgGreenColor)
	//table.SetHeaderColor(fgGreen, fgGreen, fgGreen, fgGreen, fgGreen, fgGreen, fgGreen, fgGreen)

	table.AppendBulk(datas)
	table.Render()
}

func ShowDNS(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"索引", "DNS"})
	center := tablewriter.ALIGN_CENTER
	left := tablewriter.ALIGN_LEFT
	table.SetColumnAlignment([]int{center, left})

	//fgGreen := tablewriter.Color(tablewriter.FgGreenColor)
	//table.SetHeaderColor(fgGreen, fgGreen)

	table.AppendBulk(datas)
	table.Render()
}

func ShowSimpleNode(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"索引", "别名", "地址", "端口", "测试结果"})

	center := tablewriter.ALIGN_CENTER
	left := tablewriter.ALIGN_LEFT
	table.SetColumnAlignment([]int{center, left, center, center, center})

	//fgGreen := tablewriter.Color(tablewriter.FgGreenColor)
	//table.SetHeaderColor(fgGreen, fgGreen, fgGreen, fgGreen, fgGreen)

	table.SetColWidth(70)
	table.AppendBulk(datas)
	table.Render()
}

func ShowSub(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"索引", "别名", "url", "是否启用"})

	center := tablewriter.ALIGN_CENTER
	table.SetColumnAlignment([]int{center, center, center, center})

	//fgGreen := tablewriter.Color(tablewriter.FgGreenColor)
	//table.SetHeaderColor(fgGreen, fgGreen, fgGreen, fgGreen)

	table.AppendBulk(datas)
	table.Render()
}

func ShowRouter(writer io.Writer, datas ...[]string) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"索引", "规则"})

	center := tablewriter.ALIGN_CENTER
	table.SetColumnAlignment([]int{center, center})

	//fgGreen := tablewriter.Color(tablewriter.FgGreenColor)
	//table.SetHeaderColor(fgGreen, fgGreen)

	table.AppendBulk(datas)
	table.Render()
}
