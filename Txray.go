// Txray.go 是 Txray 项目的主入口，负责初始化日志、命令行参数解析和交互式 shell 启动
package main

import (
	// 导入项目核心模块和依赖
	"Txray/cmd"           // shell 命令初始化
	"Txray/core"          // 配置目录等核心功能
	"Txray/core/setting"  // 设置相关
	"Txray/core/setting/key" // 设置项 key
	"Txray/log"           // 日志模块
	"os"                  // 操作系统相关
	"path/filepath"       // 路径处理

	"github.com/abiosoft/ishell" // 交互式 shell
	"github.com/spf13/viper"     // 配置读取
)

const (
	version = "v25.9.8.1" // 程序版本号
	name    = "Txray"   // 程序名称
)

// init 在 main 前自动执行，主要用于日志初始化
func init() {
	// 获取日志文件绝对路径
	absPath := filepath.Join(core.GetConfigDir(), "info.log")
	// 初始化日志，控制台和文件各自的 zapcore
	log.Init(
		log.GetConsoleZapcore(log.INFO),
		log.GetFileZapcore(absPath, log.INFO, 5),
	)
}

// beforeOfRun 在 shell 启动前执行预设命令（如配置的 RunBefore 脚本）
// shell: ishell.Shell 指针，交互式 shell 实例
func beforeOfRun(shell *ishell.Shell) {
	cmd := viper.GetString(key.RunBefore) // 读取配置项
	if cmd != "" {
		for _, line := range setting.NewAlias("", cmd).GetCmd() {
			shell.Process(line...) // 依次执行命令
		}
		shell.Print("\n>>> ") // 打印提示符
	}
}

// main 程序主入口，负责 shell 初始化、参数处理和主循环
func main() {
	shell := ishell.New()      // 创建 shell 实例
	cmd.InitShell(shell)       // 注册所有命令
	shell.Set("version", version) // 设置 shell 变量
	shell.Set("name", name)
	if len(os.Args) > 1 {
		// 如果有命令行参数，直接处理参数命令
		_ = shell.Process(os.Args[1:]...)
	} else {
		// 否则进入交互式 shell
		go beforeOfRun(shell) // 并发执行预设命令
		shell.Printf("%s - Xray Shell Client - %s\n", name, version)
		shell.Run()           // 启动 shell 主循环
	}
}
