package admin

import (
	"gylib/dbfun"
	"fmt"
	"gylib/common/session"
	"strconv"
	"encoding/json"
)

func (this *AdminController) AdminIndexAction() {
	this.W.Header().Set("content-type", "text/html;charset=utf-8")
	this.Tplname="views/admin/index/index.html"
	this.Rander()
}

func (this *AdminController) AdminLoginAction(){
	this.W.Header().Set("content-type", "text/html;charset=utf-8")
	this.Tplname="views/admin/index/login.html"
	this.Rander()
}

func (this *AdminController) AdminGet_nav_menuAction() {
	data := Build_menu(0)
	b, _ := json.Marshal(&data)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}

func Build_menu(p_id int)([]map[string]interface{}){
	db:=lib.NewQuerybuilder()
	data:=db.Tbname("nav").Where(fmt.Sprintf("parent_id=%v",p_id)).Select()
	list:=make([]map[string]interface{},0)
	if(len(data)>0) {
		for _,v:=range data{
			tmp:=make(map[string]interface{})
			tmp["title"]=v["nav_name"]
			tmp["icon"]=v["nav_image"]
			tmp["spread"]= false
			tmp["href"]=v["nav_module"]
			id,_:=strconv.Atoi(v["id"])
			tmp_list:=Build_menu(id)
			if(len(tmp_list)>0){
				tmp["children"]=tmp_list
			}
			list=append(list,tmp)
		}
		return list
	}else{
		return nil
	}
}

func (this *AdminController) AdminLogin_submitAction(){
	db:=lib.NewQuerybuilder()
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "登录失败"
	code:=this.R.FormValue("code")
	pass:=lib.GetMd5String(this.R.FormValue("pass"))
	data:=db.Tbname("login").Where(fmt.Sprintf("code='%v' and pass='%v' and is_del=0",code,pass)).Find()
	if(len(data)>0){
      sessions.Gsession.Set("u_id",data["id"])
	  jsonstr.Status = 100
	  jsonstr.Msg = "登录成功"
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}