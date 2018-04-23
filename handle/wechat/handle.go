package wechat

import (
	"net/http"
	"websever/controller/wechat"
)

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	mycon := wechat.NewWechatController(w, r)
	mycon.Replay()
}
