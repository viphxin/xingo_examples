package main

import (
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo"
	"github.com/viphxin/xingo/utils"
	"xingo_examples/helloword/api"
	"xingo_examples/helloword/cmd"
	"fmt"
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

       //add cmd ----------------start
       utils.GlobalObject.CmdInterpreter.AddCommand(cmd.NewTimerCommand())
       //add cmd ----------------end

	//regest callback
	utils.GlobalObject.OnConnectioned = DoConnectionMade
	utils.GlobalObject.OnClosed = DoConnectionLost
	s.Serve()
}
