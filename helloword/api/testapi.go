package api

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/viphxin/xingo/fnet"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
	"github.com/viphxin/xingo/utils"
	"github.com/viphxin/xingo_examples/helloword/pb"
	"time"
)

type TestRouter struct {
	fnet.BaseRouter
}

func (this *TestRouter) sendDelayMsg(fconn iface.Iconnection) {
	utils.GlobalObject.GetSafeTimer().CreateTimer(5000, func(args ...interface{}) {
		con := args[0].(iface.Iconnection)
		ntf := &pb.DelayNtf{
			Ts: time.Now().String(),
		}
		ntfRaw, err := utils.GlobalObject.Protoc.GetDataPack().Pack(3, ntf)
		if err == nil {
			con.Send(ntfRaw)
		}
	}, []interface{}{fconn})
}

/*
HelloReq
*/
func (this *TestRouter) Handle(request iface.IRequest) {
	msg := &pb.HelloReq{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err == nil {
		request.GetConnection().SetProperty("name", msg.Name)
		//send ack
		ack := &pb.HelloAck{
			Content: fmt.Sprintf("Hello %s.You will receive a Ntf after 5 seconds.\n", msg.Name),
		}
		data, err := utils.GlobalObject.Protoc.GetDataPack().Pack(2, ack)
		if err == nil {
			request.GetConnection().Send(data)
			this.sendDelayMsg(request.GetConnection())
		}
	} else {
		logger.Error(err)
		request.GetConnection().LostConnection()
	}
}
