package cmd

import (
	"bytes"
	"fmt"
)

func Help() string {
	var buf bytes.Buffer
	format15 := "    %-15s%-s\n"
	format25 := "    %-25s%-s\n"
	format25t17 := "    %-25s%-17s%-s\n"
	buf.WriteString("Commands:\n")
	buf.WriteString(fmt.Sprintf(format25t17, "setting", "基础设置", "使用 'setting' 查看详细用法"))
	buf.WriteString(fmt.Sprintf(format25t17, "node", "节点管理", "使用 'node' 查看详细用法"))
	buf.WriteString(fmt.Sprintf(format25t17, "sub", "订阅管理", "使用 'sub' 查看详细用法"))
	buf.WriteString(fmt.Sprintf(format25t17, "dns", "DNS管理", "  使用 'dns' 查看详细用法"))
	buf.WriteString(fmt.Sprintf(format25t17, "route", "路由管理", "使用 'route' 查看详细用法"))
	buf.WriteString(fmt.Sprintf(format25, "help, -h", "查看帮助信息"))
	buf.WriteString(fmt.Sprintf(format25, "version, -v", "查看版本"))
	buf.WriteString(fmt.Sprintf(format25, "clear", "清屏"))
	buf.WriteString(fmt.Sprintf(format25, "exit", "退出程序"))
	buf.WriteString(fmt.Sprintf(format25, "stop", "停止服务"))
	buf.WriteString(fmt.Sprintf(format25, "run", "启动或重启服务"))
	buf.WriteString("\nUsage: run [索引式 | -t [索引式]]\n\n")
	buf.WriteString(fmt.Sprintf(format15, "run [索引式]", "默认为上一次运行节点，如果选中多个节点，则选择访问YouTube延迟最小的"))
	buf.WriteString(fmt.Sprintf(format15, "run -t [索引式]", "按tcp延迟选择节点，默认'1'，比如 'run -t 1-10' 为选择tcp延迟最小的10个节点"))
	buf.WriteString("\n\n说明：\n")
	buf.WriteString("一、索引式：更简单的批量选择\n")
	buf.WriteString("1.选择前6个：'1,2,3,4,5,6' 或 '1-3,4-6' 或 '1-6' 或 '-6'\n")
	buf.WriteString("2.选择第6个及后面的所有：'6-'\n")
	buf.WriteString("3.选择第6个：'6'\n")
	buf.WriteString("4.选择所有：'all' 或 '-'\n")
	buf.WriteString("注意：超出部分会被忽略，'all' 只能单独使用\n\n")
	buf.WriteString("二、[] 和 {}：帮助说明中的中括号和大括号\n")
	buf.WriteString("1. []: 表示该选项可忽略\n")
	buf.WriteString("2. {}: 表示该选项为必须，不可忽略\n")
	return buf.String()
}

func HelpSetting() string {
	var buf bytes.Buffer
	format30 := "    %-30s%-s\n"
	buf.WriteString("setting {commands} [flags] ...\n")
	buf.WriteString("\n")
	buf.WriteString("Commands:\n")
	buf.WriteString(fmt.Sprintf(format30, "show", "查看基本设置"))
	buf.WriteString(fmt.Sprintf(format30, "alter [flags]", "修改基础设置"))
	buf.WriteString("\n")
	buf.WriteString("alter Flags\n")
	buf.WriteString(fmt.Sprintf(format30, "-p, --port {port}", "设置socks端口"))
	buf.WriteString(fmt.Sprintf(format30, "-h, --http {port}", "设置http端口, 0为关闭http监听"))
	buf.WriteString(fmt.Sprintf(format30, "-u, --udp {y|n}", "是否启用udp"))
	buf.WriteString(fmt.Sprintf(format30, "-s, --sniffing {y|n}", "是否启用流量监听"))
	buf.WriteString(fmt.Sprintf(format30, "-l, --lanconn {y|n}", "是否启用局域网连接"))
	buf.WriteString(fmt.Sprintf(format30, "-m, --mux {y|n}", "是否启用多路复用"))
	buf.WriteString(fmt.Sprintf(format30, "-b, --bypass {y|n}", "是否绕过局域网及大陆"))
	buf.WriteString(fmt.Sprintf(format30, "-r, --route {1|2|3}", "设置路由策略为{AsIs|IPIfNonMatch|IPOnDemand}"))
	return buf.String()
}

func HelpNode() string {
	var buf bytes.Buffer
	format30 := "    %-30s%-s\n"
	format28 := "    %-28s%-s\n"
	format27 := "    %-27s%-s\n"
	buf.WriteString("node {commands} [flags] ...\n")
	buf.WriteString("\n")
	buf.WriteString("Commands:\n")
	buf.WriteString(fmt.Sprintf(format27, "show [索引式|t]", "查看节点信息, 默认'all', 't'表示按延迟降序查看"))
	buf.WriteString(fmt.Sprintf(format28, "info {索引}", "查看单个节点详细信息"))
	buf.WriteString(fmt.Sprintf(format27, "del {索引式}", "删除节点"))
	buf.WriteString(fmt.Sprintf(format27, "tcping {索引式}", "测试节点tcp延迟"))
	buf.WriteString(fmt.Sprintf(format27, "find {关键词}", "查找节点（按别名）"))
	buf.WriteString(fmt.Sprintf(format30, "add [flags]", "添加节点"))
	buf.WriteString(fmt.Sprintf(format27, "export [索引式] [flags]", "导出节点链接, 默认'all'"))
	buf.WriteString("\nadd Flags\n")
	buf.WriteString(fmt.Sprintf(format30, "-l, --link {link}", "从链接导入一条节点"))
	buf.WriteString(fmt.Sprintf(format30, "-f, --file {path}", "从节点链接文件或订阅文件导入节点"))
	buf.WriteString(fmt.Sprintf(format30, "-c, --clipboard", "从剪贴板读取的节点链接或订阅文本导入节点"))
	buf.WriteString("\nexport Flags\n")
	buf.WriteString(fmt.Sprintf(format30, "-c, --clipboard", "导出节点链接到剪贴板"))
	return buf.String()
}

func HelpSub() string {
	var buf bytes.Buffer
	format30 := "    %-30s%-s\n"
	format28 := "    %-28s%-s\n"
	format27 := "    %-27s%-s\n"
	buf.WriteString("sub {commands} [flags] ...\n")
	buf.WriteString("\n")
	buf.WriteString("Commands:\n")
	buf.WriteString(fmt.Sprintf(format27, "show {索引式}", "查看订阅信息"))
	buf.WriteString(fmt.Sprintf(format27, "del {索引式}", "删除订阅"))
	buf.WriteString(fmt.Sprintf(format28, "add {订阅url} [flags]", "添加订阅"))
	buf.WriteString(fmt.Sprintf(format27, "alter {索引式} {flags}", "修改订阅"))
	buf.WriteString(fmt.Sprintf(format27, "update-node [索引式] [flags]", "从订阅更新节点, 索引式会忽略是否启用"))
	buf.WriteString("\nadd Flags\n")
	buf.WriteString(fmt.Sprintf(format28, "-r, --remarks {别名} ", "定义别名"))
	buf.WriteString("\nalter Flags\n")
	buf.WriteString(fmt.Sprintf(format28, "-u, --url {订阅url}", "修改订阅链接"))
	buf.WriteString(fmt.Sprintf(format28, "-r, --remarks {别名}", "定义别名"))
	buf.WriteString(fmt.Sprintf(format30, "--using {y|n}", "是否启用此订阅"))
	buf.WriteString("\nupdate-node Flags\n")
	buf.WriteString(fmt.Sprintf(format30, "-s, --socks5 [port]", "通过本地的socks5代理更新, 默认为设置中的socks5端口"))
	buf.WriteString(fmt.Sprintf(format30, "-h, --http [port]", "通过本地的http代理更新, 默认为设置中的http端口"))
	buf.WriteString(fmt.Sprintf(format30, "-a, --addr {address}", "对上面两个参数的补充, 修改代理地址"))
	return buf.String()
}

func HelpDNS() string {
	var buf bytes.Buffer
	format30 := "    %-30s%-s\n"
	buf.WriteString("dns {commands} [flags] ...\n")
	buf.WriteString("\n")
	buf.WriteString("Commands:\n")
	buf.WriteString(fmt.Sprintf(format30, "show", "查看DNS设置"))
	buf.WriteString(fmt.Sprintf(format30, "alter {flags}", "修改DNS设置"))
	buf.WriteString("\n")
	buf.WriteString("alter Flags\n")
	buf.WriteString(fmt.Sprintf(format30, "-p, --port {port}", "设置dns端口, 0为关闭"))
	buf.WriteString(fmt.Sprintf(format30, "-i, --inland {dns}", "设置一条境内DNS"))
	buf.WriteString(fmt.Sprintf(format30, "-o, --outland {dns}", "设置一条境外DNS"))
	buf.WriteString(fmt.Sprintf(format30, "-b, --backup {dns}", "设置备用DNS，多条以 ',' 分隔"))
	return buf.String()
}

func HelpRoute() string {
	var buf bytes.Buffer
	format30 := "    %-30s%-s\n"
	format27 := "    %-27s%-s\n"
	format26 := "    %-26s%-s\n"
	buf.WriteString("route {commands} [flags] ...\n")
	buf.WriteString("\n")
	buf.WriteString("Commands:\n")
	buf.WriteString(fmt.Sprintf(format26, "add {路由规则} {flags}", "添加路由规则"))
	buf.WriteString(fmt.Sprintf(format30, "show {flags}", "查看路由规则"))
	buf.WriteString(fmt.Sprintf(format30, "del {flags}", "删除路由规则"))
	buf.WriteString("\nadd Flags\n")
	buf.WriteString(fmt.Sprintf(format30, "-b, --block", "指定到禁止名单, 同下面2个中选择一个"))
	buf.WriteString(fmt.Sprintf(format30, "-d, --direct", "指定到直连名单, 同上下2个中选择一个"))
	buf.WriteString(fmt.Sprintf(format30, "-p, --proxy", "指定到代理名单, 同上面2个中选择一个"))
	buf.WriteString(fmt.Sprintf(format30, "-f, --file {path}", "从文件导入"))
	buf.WriteString(fmt.Sprintf(format30, "-c, --clipboard", "从剪贴板导入"))
	buf.WriteString("\nshow Flags\n")
	buf.WriteString(fmt.Sprintf(format27, "-b, --block [索引式]", "指定到禁止名单, 同下面2个中选择一个, 索引式默认'all'"))
	buf.WriteString(fmt.Sprintf(format27, "-d, --direct [索引式]", "指定到直连名单, 同上下2个中选择一个, 索引式默认'all'"))
	buf.WriteString(fmt.Sprintf(format27, "-p, --proxy [索引式]", "指定到代理名单, 同上面2个中选择一个, 索引式默认'all'"))
	buf.WriteString("\ndel Flags\n")
	buf.WriteString(fmt.Sprintf(format27, "-b, --block {索引式}", "指定到禁止名单, 同下面2个中选择一个"))
	buf.WriteString(fmt.Sprintf(format27, "-d, --direct {索引式}", "指定到直连名单, 同上下2个中选择一个"))
	buf.WriteString(fmt.Sprintf(format27, "-p, --proxy {索引式}", "指定到代理名单, 同上面2个中选择一个"))
	return buf.String()
}
