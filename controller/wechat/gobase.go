package wechat

import (
	"net/http"
	"websever/public/controller"
	"gylib/common"
	"gylib/weixinsdk/weixinmp"
	"gylib/dbfun"
)

type WechatController struct {
	Site_name string
	Api_url   string
	Get_url   string
	controller.Controller
}

type Json_msg struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func NewWechatController(w http.ResponseWriter, r *http.Request) (*WechatController) {
	this := new(WechatController)
	this.W = w
	this.R = r
	this.Data = make(map[string]interface{})
	this.Url_list = make(map[string]interface{})
	this.JsonData = make(map[string]interface{})
	this.Postdata = make(map[string]interface{})
	this.Get_url = r.URL.RawQuery
	data := common.Getini("conf/app.ini", "server", map[string]string{"site_name": "", "api_url": ""})
	this.Site_name = data["site_name"]
	this.Api_url = data["api_url"]
	return this
}

func (this *WechatController) Replay() {
	db:=lib.NewQuerybuilder()
	data:=db.Tbname("tp_wxuser").Find()
	mp := weixinmp.New(data["token"], data["appid"], data["appsecret"])
	// 检查请求是否有效
	// 仅主动发送消息时不用检查
	if !mp.Request.IsValid(this.W, this.R) {
		return
	}
	// 判断消息类型
	if mp.Request.MsgType == weixinmp.MsgTypeText {
		// 回复消息
		mp.ReplyTextMsg(this.W, "Hello, 世界")
	}
}
