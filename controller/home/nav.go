package home

import (
	"strconv"
	"fmt"
	"gylib/common"
	"gylib/dbfun"
	"encoding/json"
)


func (this *HomeController) HomeNav_indexAction() {
	this.Tplname = "views/nav/index.html"
	this.Rander()
}

func (this *HomeController) HomeNav_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page-1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("nav").Get_where_data(this.Postdata)
	//fmt.Println(where)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	}else{
		order_str += " "+this.R.FormValue("order_type")
	}
	ct := db.Tbname("nav").Where(where).Count()
	data := db.Tbname("nav").Where(where).Order(order_str).Limit(limit_str).Select()
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

/* 获取当条记录Api */
func (this *HomeController) HomeNav_editAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	id_str := this.R.FormValue("id")
	id, err := strconv.Atoi(id_str)
	data := make(map[string]string)
	if (err != nil) {
		return
	}
	if (id > 0) {
		data = db.Tbname("nav").Where(fmt.Sprintf("id=%v", id)).Find()
	} else {
		data = db.Tbname("nav").Get_new_add()
	}
	jsonstr := new(Json_msg)
	jsonstr.Status = 100
	jsonstr.Msg = "获取数据成功"
	jsonstr.Data = data
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}

/* 保存记录Api */
func (this *HomeController) HomeNav_saveAction() {
	db := lib.NewQuerybuilder()
	msg := "保存失败"
	temp_file := common.Upload_File(this.R, "img_file")
	if (len(temp_file) > 0) {
		if (temp_file[len(temp_file)-1] != "") {
			this.Postdata["img"] = temp_file[len(temp_file)-1]
		}
	}
	var err error
	if (this.R.FormValue("id") != "") {
		Id := this.R.FormValue("id")
		_, err = db.Tbname("nav").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("nav").Insert(this.Postdata)

	}
	jsonstr := new(Json_msg)
	if (err != nil) {
		jsonstr.Status = 100
		msg="保存成功"
	} else {
		jsonstr.Status = 101
	}
	jsonstr.Msg = msg
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)

}

func (this *HomeController) HomeNav_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	var Id int
	var err error
	Id, err = strconv.Atoi(this.R.FormValue("id"))
	if err == nil {
		_, err = db.Tbname("nav").Where(fmt.Sprintf("id in (%v)", Id)).Delete()
		jsonstr.Msg = "删除成功"
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}