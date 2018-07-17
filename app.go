package main

import (

	"github.com/netc0/netco"
	"github.com/netc0/netco/def"
	"github.com/netc0/netco/common"
	"github.com/netc0/netco/rpc"
	"log"
)

type Game struct {
	def.IService
	netco.App
	rpc.MailHandler

	//gateAddress string // 网关地址
	mailBox rpc.IMailBox

}

var (
	logger = common.GetLogger()
)

func (this *Game) OnStart() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logger.Debug("游戏开始")
	go this.mailBox.Start()
	if err := this.mailBox.ConnectGate(this.GetGateAddress()); err != nil {
		log.Println("连接网关失败", err)
	}

}
func (this *Game) OnDestroy() {
	logger.Debug("游戏结束")
}

func NewGame() *Game{
	game := new(Game)
	game.Derived = game
	game.SetGateAddress("127.0.0.1:9002")
	game.SetNodeAddress("127.0.0.1:1899")
	game.SetNodeName("example-game")

	game.mailBox = rpc.NewMailBox(game.GetNodeAddress())
	game.mailBox.SetHandler(game)

	log.Printf("new game: %p, %p", game, game.OnRoutineConnected)
	return game
}

func main() {
	var game = NewGame()
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
func (this *Game)processRequest(req def.MailClientInfo) {
	switch req.Route {
	case rpc.RouteHash("game.login"):
		logger.Debug("执行登录")
		// 登录逻辑
		var i def.MailClientInfo
		i.RequestId = req.RequestId
		i.ClientId = req.ClientId
		i.SourceAddress = this.GetNodeAddress()
		i.SourceName = this.GetNodeName()

		i.Data = []byte("你好, 欢迎")
		this.mailBox.SendToGate(&rpc.Mail{Type:def.Mail_ResponseData, Object:i})
		// 通知网关, 如果这个会话断开则通知我
		i.RequestId = 0
		i.ClientId = req.ClientId
		i.Data = nil
		this.mailBox.SendToGate(&rpc.Mail{Type:def.Mail_ClientLeaveNotifyMe, Object:i})

		// 推送一条数据
		var pdata  def.MailClientInfo
		pdata.Route = rpc.RouteHash("game.push")
		pdata.Data  = []byte("推送一条数据");
		pdata.ClientId = req.ClientId

		this.mailBox.SendToGate(&rpc.Mail{Type:def.Mail_PushData, Object:pdata})

	case rpc.RouteHash("game.join"):
		logger.Debug("执行加入", req)
	}
}