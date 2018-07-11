package main

import (

	"github.com/netc0/netco"
	"github.com/netc0/netco/def"
	"github.com/netc0/netco/common"
	"github.com/netc0/netco/rpc"
	"log"
)
//
//type ExampleContext struct {
//	gateRPC *rpc.Client
//	gateLock *sync.Mutex
//
//	nodeName string
//	auth string
//}
//
//func runRPCServer(context *ExampleContext) {
//	var v = new(Example)
//	v.context = context;
//	netco.RPCServerStart(":8001", v)
//}
//
//func connectGate(context *ExampleContext) {
//	context.gateLock.Lock()
//	defer context.gateLock.Unlock()
//
//	if context.gateRPC != nil {
//		// 已经连接 connected
//		return
//	}
//
//	cli, err := netco.RPCClientConnect("127.0.0.1:9002")
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	context.gateRPC = cli
//	log.Println("注册 Proxy")
//
//	var info netco.RPCBackendInfo
//	info.RCPRemote = "127.0.0.1:8001"
//
//	info.Name = context.nodeName
//	info.Auth = context.auth
//	info.Routes = append(info.Routes, "Example.Test")
//	info.Routes = append(info.Routes, "Example.Login")
//
//	reply := 0
//	rs := context.gateRPC.Call("GateProxy.RegisterBackend", info, &reply)
//	log.Println("GateProxy reply:", rs)
//}
//
//// 监控网关的连接状态
//func gateMonitor(context *ExampleContext) {
//	ticker := time.NewTicker(time.Second * 3)
//	for range ticker.C {
//		if context.gateRPC == nil {
//			go connectGate(context)
//		} else {
//			go gateHeartBeat(context)
//		}
//	}
//}
//
//func gateHeartBeat(context *ExampleContext) {
//	if context.gateRPC != nil {
//		var info netco.RPCBackendInfo
//		info.Name = context.nodeName
//		info.Auth = context.auth
//		rs := context.gateRPC.Call("GateProxy.BackendHeartBeat", info, nil)
//		if rs != nil { // 断开连接
//			log.Println(rs.Error())
//			context.gateRPC = nil
//		}
//	}
//}

type Game struct {
	def.IService
	netco.App
	rpc.MailHandler

	gateAddress string // 网关地址
	mailBox rpc.IMailBox

}

var (
	logger = common.GetLogger()
)

func (this *Game) OnStart() {
	logger.Debug("游戏开始")
	this.mailBox = rpc.NewMailBox(":9003")
	this.mailBox.SetHandler(this)
	go this.mailBox.Start()
	if err := this.mailBox.Connect(this.gateAddress); err != nil {
		log.Println("连接网关失败", err)
	}

}
func (this *Game) OnDestroy() {
	logger.Debug("游戏结束")
}

func main() {
	//log.SetFlags(log.LstdFlags|log.Lshortfile)
	//var ctx ExampleContext
	//ctx.gateLock = new(sync.Mutex)
	//ctx.nodeName = "exampleGame"
	//ctx.auth = "netc0"
	//log.Println("启动游戏")

	//go runRPCServer(&ctx)
	//go connectGate(&ctx)
	//go gateMonitor(&ctx)
	//
	//for {
	//	time.Sleep(time.Second)
	//}

	var game Game
	game.Derived = &game
	game.gateAddress = "127.0.0.1:9002"
	game.Run()
}

// 导出路由给网关
func (this* Game)exportRoute() []uint32{
	var routes []uint32
	routes = append(routes, rpc.RouteHash("game.login"))
	routes = append(routes, rpc.RouteHash("game.join"))
	return routes
}

// 处理客户端的请求
func (this *Game)processRequest(req def.MailClientData) {
	switch req.Route {
	case rpc.RouteHash("game.login"):
		logger.Debug("执行登录")
		var i def.MailClientData
		i.RequestId = req.RequestId
		i.ClientId = req.ClientId
		i.Data = []byte("你好, 欢迎")
		this.mailBox.SendTo(this.gateAddress, &rpc.Mail{Type:def.Mail_ResponseData, Object:i})
	case rpc.RouteHash("game.join"):
		logger.Debug("执行加入")
	}
}