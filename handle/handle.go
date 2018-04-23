package handle

import (
	"net/http"
	"strings"
	"encoding/json"
	"reflect"
	"websever/controller/admin"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	switch_action(parts, w, r)
}

func switch_action(parts []string, w http.ResponseWriter, r *http.Request) {
	var method_name string
	var action string

	switch len(parts) {
	case 0:
		method_name = "Admin"
		action = "AdminIndexAction"
	case 1:
		method_name = strings.Title(parts[0])
		action = method_name + "IndexAction"
	default:
		method_name = strings.Title(parts[0])
		for _, val := range parts {
			if (action == "") {
				action = strings.Title(val)
			} else {
				action += strings.Title(val)
			}
		}
		action += "Action"
		//action = strings.Title(parts[0]) + strings.Title(parts[1]) + "Action"
	}

	switch strings.Title(method_name) {
	case "Admin": //登录接口--需要检查key
		controller_mothod,flag := admin.NewAdminController(w, r,action)
		if(!flag){
			http.Redirect(w, r, "/admin/login", 302)
			return
		}
		controller_name := reflect.ValueOf(controller_mothod)
		method := controller_name.MethodByName(action)
		//fmt.Println(action)
		if !method.IsValid() {
			error_return(w, "路由错误")
			return
			//method = controller_name.MethodByName("ApiLoginAction")
		}
		//requestValue := reflect.ValueOf(r)
		//responseValue := reflect.ValueOf(w)
		//method.Call([]reflect.Value{responseValue, requestValue})
		method.Call([]reflect.Value{})
	case "Api": //form接口--需要检查token

	case "Sdk": //json接口解析--需要检查token

	case "Public": //json接口解析--通用接口，无需鉴权

	default:
		error_return(w, "路由错误")
		return
	}

}

func error_return(w http.ResponseWriter, msg string) {
	jsonstr := make(map[string]interface{})
	jsonstr["status"] = 101
	jsonstr["msg"] = msg
	b, _ := json.Marshal(&jsonstr)
	w.Write(b)
}
