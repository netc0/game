package main

import (
	"github.com/netc0/netco/rpc"
	"github.com/netc0/netco/def"
)

func (this* Game) OnNewMail(mail rpc.Mail) {
	if mail.Type == 0 { return }

	if mail.Type == def.Mail_RequestData {
		// 来自客户端的请求
		var v def.MailClientInfo
		if err := mail.Decode(&v); err != nil {
			logger.Debug(err)
			return
		}

		this.processRequest(v)
	} else if mail.Type == def.Mail_ClientLeaveNotification {
		// 客户端断开
		var v def.MailClientInfo
		if err := mail.Decode(&v); err != nil {
			logger.Debug(err)
			return
		}
		logger.Debug("客户端断开:", v)
	}
}

func (this* Game)OnRoutineConnected(remote string) {
	logger.Debug("路径连接成功:", remote)
	if remote == this.GetGateAddress() {
		// 在网关注册我的信息
		if err := this.mailBox.SendTo(remote, &rpc.Mail{Type:def.Mail_Reg,
			Object:def.MailNodeInfo{Name:this.GetNodeName(), Address:this.GetNodeAddress()}});
			err != nil {
			logger.Debug(err)
			return
		}
		// 在网关注册路由信息
		var info def.MailRoutineInfo
		info.Name = this.GetNodeName()
		info.Routes = this.exportRoute()
		logger.Debug(info)
		if err := this.mailBox.SendTo(remote, &rpc.Mail{Type:def.Mail_AddRoute, Object:info}); err != nil {
			logger.Debug(err)
		}
	}
}

func (this* Game)OnRoutineDisconnect(remote string, err error) {
	logger.Debug("已经断开:", remote, err)
	if remote != this.GetGateAddress() { // 如果是网关则不用删除, 不删除会一直自动重连
		this.mailBox.Remove(remote)
	}
}
