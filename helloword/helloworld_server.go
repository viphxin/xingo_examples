package main

import (
	"fmt"
	"github.com/viphxin/xingo"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"
	"github.com/viphxin/xingo_examples/helloword/api"
)

func DoConnectionMade(fconn iface.Iconnection) {
	logger.Debug(fmt.Sprintf("session %d connectioned helloworld server.", fconn.GetSessionId()))
}

func DoConnectionLost(fconn iface.Iconnection) {
	logger.Debug(fmt.Sprintf("session %d disconnectioned helloworld server.", fconn.GetSessionId()))
}

func main() {
	s := xingo.NewXingoTcpServer()

	//add api ---------------start
	TestRouterObj := &api.TestRouter{}
	s.AddRouter("1", TestRouterObj)
	//add api ---------------end

	//regest callback
	utils.GlobalObject.OnConnectioned = DoConnectionMade
	utils.GlobalObject.OnClosed = DoConnectionLost
	s.Serve()
}
