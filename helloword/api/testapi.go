package api

import (
	"github.com/viphxin/xingo/fnet"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/logger"
	"fmt"
	"xingo_examples/helloword/pb"
	"github.com/viphxin/xingo/utils"
	"github.com/viphxin/xingo/iface"
	"time"
)

type TestRouter struct{}

func sendDelayMsg(fconn iface.Iconnection){
	utils.GlobalObject.GetSafeTimer().CreateTimer(5000, func(args ...interface{}){
		con := args[0].(iface.Iconnection)
		ntf := &pb.DelayNtf{
			Ts: time.Now().String(),
		}
		ntfRaw, err := utils.GlobalObject.Protoc.GetDataPack().Pack(3, ntf)
		if err == nil {
			con.Send(ntfRaw)
		}
	},[]interface{}{fconn})
}

/*
HelloReq
*/
func (this *TestRouter)HelloReq_1(request *fnet.PkgAll){
	msg := &pb.HelloReq{}
	err := proto.Unmarshal(request.Pdata.Data, msg)
	if err == nil {
		request.Fconn.SetProperty("name", msg.Name)
		//send ack
		ack := &pb.HelloAck{
			Content: fmt.Sprintf("Hello %s.You will receive a Ntf after 5 seconds.\n", msg.Name),
		}
		data, err := utils.GlobalObject.Protoc.GetDataPack().Pack(2, ack)
		if err == nil{
			request.Fconn.Send(data)
			sendDelayMsg(request.Fconn)
		}
	} else {
		logger.Error(err)
		request.Fconn.LostConnection()
	}
}
