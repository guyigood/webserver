package main

import (
	"os"
	"net/http"
	"fmt"
	"runtime"
	"websever/handle"
	"gylib/common"
	"websever/handle/route"
	"websever/handle/wechat"
	"websever/handle/layroute"
)

func main(){
	var port string = ""
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	var MULTICORE int = runtime.NumCPU() //number of core
	runtime.GOMAXPROCS(MULTICORE)        //running in multicore
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//http.Handle("/chai.jpg",http.FileServer(http.Dir("static")))
	http.HandleFunc("/admin/",handle.AdminHandler)
	http.HandleFunc("/wechat/",wechat.ApiHandler)
	http.HandleFunc("/jxc/",route.ApiHandler)
	http.HandleFunc("/lay/",layroute.LayApiHandler)
	data := common.Getini("conf/app.ini", "server", map[string]string{"port": ""})
	if port == "" {
		port = data["port"]
	}
	fmt.Println(data)
	http.ListenAndServe(":"+port, nil)
}
