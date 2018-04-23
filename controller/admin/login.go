package admin

import (
	"strconv"
	"fmt"
	"gylib/common"
	"gylib/dbfun"
	"encoding/json"
	"gylib/common/datatype"
	"strings"
)

func (this *AdminController) AdminLogin_indexAction() {
	this.Tplname = "views/admin/login/index.html"
	db := lib.NewQuerybuilder()
	db.Dbinit()
	this.Data["qxz"] = db.Tbname("qxz").Where("").Select()
	this.Rander()
}

func (this *AdminController) AdminLogin_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("login").Get_where_data(this.Postdata)
	if(where==""){
		where="is_del=0"
	}else{
		where+=" and is_del=0"
	}
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	ct := db.Tbname("login").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("login").Where(where).Order(order_str).Limit(limit_str).Select()
	for k,v:=range data{
		db.Dbinit()
		qxz:=db.Tbname("qxz").Where("id="+v["qxz"]).Find()
		if(qxz!=nil){
			data[k]["qxz"]=qxz["name"]
		}
		db.Dbinit()
		r_qxz:=db.Tbname("routeqxz").Where("id="+v["r_id"]).Find()
		if(r_qxz!=nil){
			data[k]["r_id"]=qxz["name"]
		}
	}
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

/* 获取当条记录Api */
func (this *AdminController) AdminLogin_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"
		data = db.Tbname("login").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
		db.Dbinit()
		data["dept_name"]=""
		dept:=db.Tbname("depart").Where("id in("+data["dept_id"]+")").Select()
		for _,v:=range dept{
			if(data["dept_name"]==""){
				data["dept_name"]=v["name"]
			}else{
				data["dept_name"]+=","+v["name"]
			}
		}
	} else {
		data = db.Tbname("login").Get_new_add()
		data["dept_name"]=""
	}
	this.Data["edit_flag"] = edit_flag
	db.Dbinit()
	this.Data["qxz"] = db.Tbname("qxz").Select()
	db.Dbinit()
	this.Data["routeqxz"] = db.Tbname("routeqxz").Select()
	this.Data["data"] = datatype.MapString2interface(data)
	this.Tplname = "views/admin/login/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminLogin_saveAction() {
	db := lib.NewQuerybuilder()
	msg := "保存失败"
	form_type := this.R.Header.Get("Content-Type")
	if strings.Contains(form_type, "multipart/form-data;") {
		temp_file := common.Upload_File(this.R, "img_file")
		if (len(temp_file) > 0) {
			if (temp_file[len(temp_file)-1] != "") {
				this.Postdata["img"] = temp_file[len(temp_file)-1]
			}
		}
	}
	var err error
	if (this.R.FormValue("id") != "") {
		Id := this.R.FormValue("id")
		_, err = db.Tbname("login").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("login").Insert(this.Postdata)

	}
	jsonstr := new(Json_msg)
	if (err == nil) {
		jsonstr.Status = 100
		msg = "保存成功"
	} else {
		jsonstr.Status = 101
	}
	jsonstr.Msg = msg
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)

}

func (this *AdminController) AdminLogin_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("login").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Update(map[string]interface{}{"is_del":1})
	if (err == nil) {
		jsonstr.Msg = "删除成功"
		jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}
