package main

import (
	"io"
	"github.com/viphxin/xingo/fnet"
	"time"
	"os"
	"os/signal"
	"fmt"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"xingo_examples/helloword/pb"
	"github.com/viphxin/xingo/iface"
	"github.com/viphxin/xingo/logger"
)

type HelloWorldCPtotoc struct{
	Name string
}

func (this *HelloWorldCPtotoc)OnConnectionMade(fconn iface.Iclient){
	fmt.Println("链接建立")
	req := &pb.HelloReq{
		Name: this.Name,
	}
	this.Send(fconn, 1, req)
}

func (this *HelloWorldCPtotoc)OnConnectionLost(fconn iface.Iclient){
	fmt.Println("链接丢失")
}

func (this *HelloWorldCPtotoc) Unpack(headdata []byte) (head *fnet.PkgData, err error) {
	headbuf := bytes.NewReader(headdata)

	head = &fnet.PkgData{}

	// 读取Len
	if err = binary.Read(headbuf, binary.LittleEndian, &head.Len); err != nil {
		return nil, err
	}

	// 读取MsgId
	if err = binary.Read(headbuf, binary.LittleEndian, &head.MsgId); err != nil {
		return nil, err
	}

	// 封包太大
	//if head.Len > MaxPacketSize {
	//	return nil, packageTooBig
	//}

	return head, nil
}

func (this *HelloWorldCPtotoc) Pack(msgId uint32, data proto.Message) (out []byte, err error) {
	outbuff := bytes.NewBuffer([]byte{})
	// 进行编码
	dataBytes := []byte{}
	if data != nil {
		dataBytes, err = proto.Marshal(data)
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("marshaling error:  %s", err))
	}
	// 写Len
	if err = binary.Write(outbuff, binary.LittleEndian, uint32(len(dataBytes))); err != nil {
		return
	}
	// 写MsgId
	if err = binary.Write(outbuff, binary.LittleEndian, msgId); err != nil {
		return
	}

	//all pkg data
	if err = binary.Write(outbuff, binary.LittleEndian, dataBytes); err != nil {
		return
	}

	out = outbuff.Bytes()
	return

}

func (this *HelloWorldCPtotoc)DoMsg(fconn iface.Iclient, pdata *fnet.PkgData){
	//处理消息
	fmt.Println(fmt.Sprintf("msg id :%d, data len: %d", pdata.MsgId, pdata.Len))
	switch pdata.MsgId {
	case 2:
		ack := &pb.HelloAck{}
		err := proto.Unmarshal(pdata.Data, ack)
		if err == nil {
			logger.Debug(ack.Content)
		}else{
			logger.Error("Unmarshal ack err: ", err)
		}
	case 3:
		nft := &pb.DelayNtf{}
		err := proto.Unmarshal(pdata.Data, nft)
		if err == nil {
			logger.Debug(nft.Ts)
		}else{
			logger.Error("Unmarshal ntf err: ", err)
		}
	default:
		logger.Error("Unkown message!!!!")
	}
}

func (this *HelloWorldCPtotoc)Send(fconn iface.Iclient, msgID uint32, data proto.Message){
	dd, err := this.Pack(msgID, data)
	if err == nil{
		fconn.Send(dd)
	}else{
		fmt.Println(err)
	}

}

func (this *HelloWorldCPtotoc)StartReadThread(fconn iface.Iclient){
	go func() {
		for {
			//read per head data
			headdata := make([]byte, 8)

			if _, err := io.ReadFull(fconn.GetConnection(), headdata); err != nil {
				fmt.Println(err)
				this.OnConnectionLost(fconn)
				return
			}
			pkgHead, err := this.Unpack(headdata)
			if err != nil {
				this.OnConnectionLost(fconn)
				return
			}
			//data
			if pkgHead.Len > 0 {
				pkgHead.Data = make([]byte, pkgHead.Len)
				if _, err := io.ReadFull(fconn.GetConnection(), pkgHead.Data); err != nil {
					this.OnConnectionLost(fconn)
					return
				}
			}
			this.DoMsg(fconn, pkgHead)
		}
	}()
}

func (this *HelloWorldCPtotoc)AddRpcRouter(interface{}){

}
func (this *HelloWorldCPtotoc)GetMsgHandle() iface.Imsghandle{
	return nil
}
func (this *HelloWorldCPtotoc)GetDataPack() iface.Idatapack{
	return nil
}

func (this *HelloWorldCPtotoc)InitWorker(int32){}

func main() {
	for i := 0; i< 100; i ++{
		client := fnet.NewTcpClient("0.0.0.0", 8999, &HelloWorldCPtotoc{fmt.Sprintf("xingo_fans_%d", i)})
		client.Start()
		time.Sleep(1*time.Second)
	}

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	fmt.Println("=======", sig)
}
