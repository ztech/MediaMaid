package main

import (
	_ "mediamaid/bootstrap"
	cmd "mediamaid/command"
)

// 开发非http接口类服务入口
func main() {
	//  设置运行模式为  cli(console)
	cmd.Execute()
}
