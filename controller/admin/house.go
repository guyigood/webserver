package admin

import (
	"strconv"
	"fmt"
	"gylib/common"
	"gylib/dbfun"
	"encoding/json"
	"strings"
	"websever/public/module"
)

func (this *AdminController) AdminHouse_indexAction() {
	this.Tplname = "views/admin/house/index.html"
	this.Rander()
}

func (this *AdminController) AdminHouse_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("house").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	ct := db.Tbname("house").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("house").Where(where).Order(order_str).Limit(limit_str).Select()
	china := module.NewOrgStruct()
	for k, v := range data {
		data[k]["city"] = china.Get_City_by_map(v)
		data[k]["unit_no"] = v["dong_no"] + data[k]["unit_no"]
		db.Dbinit()
		d_data:=db.Tbname("depart").Where("id="+v["w_id"]).Find()
		data[k]["w_id"]=d_data["name"]
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
func (this *AdminController) AdminHouse_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"
		data = db.Tbname("house").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
	} else {
		data = db.Tbname("house").Get_new_add()
	}
	this.Data["edit_flag"] = edit_flag

	this.Data["data"] = data
	this.Tplname = "views/admin/house/edit.html"
	this.MuitplRander("views/admin/house/house_xx.html")
}

/* 保存记录Api */
func (this *AdminController) AdminHouse_saveAction() {
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
		_, err = db.Tbname("house").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("house").Insert(this.Postdata)

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

func (this *AdminController) AdminHouse_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("house").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
	if (err == nil) {
		jsonstr.Msg = "删除成功"
		jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}
