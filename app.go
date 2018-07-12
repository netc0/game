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

	gateAddress string // 网关地址
	mailBox rpc.IMailBox

}

var (
	logger = common.GetLogger()
)

func (this *Game) OnStart() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logger.Debug("游戏开始")
	go this.mailBox.Start()
	if err := this.mailBox.Connect(this.gateAddress); err != nil {
		log.Println("连接网关失败", err)
	}

}
func (this *Game) OnDestroy() {
	logger.Debug("游戏结束")
}

func NewGame() *Game{
	game := new(Game)
	game.mailBox = rpc.NewMailBox(":9003")
	game.mailBox.SetHandler(game)
	game.Derived = game
	game.gateAddress = "127.0.0.1:9002"

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
func (this *Game)processRequest(req def.MailClientData) {
	switch req.Route {
	case rpc.RouteHash("game.login"):
		logger.Debug("执行登录")
		// 登录逻辑
		var i def.MailClientData
		i.RequestId = req.RequestId
		i.ClientId = req.ClientId
		i.SourceAddress = "127.0.0.1:9003"
		i.SourceName = "game"
		i.Data = []byte("你好, 欢迎")
		this.mailBox.SendTo(this.gateAddress, &rpc.Mail{Type:def.Mail_ResponseData, Object:i})
		// 通知网关, 如果这个会话断开则通知我
		i.RequestId = 0
		i.ClientId = req.ClientId
		i.Data = nil
		this.mailBox.SendTo(this.gateAddress, &rpc.Mail{Type:def.Mail_ClientLeaveNotifyMe, Object:i})

	case rpc.RouteHash("game.join"):
		logger.Debug("执行加入")
	}
}