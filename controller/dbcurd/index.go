package dbcurd

import (
	"gylib/dbfun"
	"strconv"
	"encoding/json"
	"fmt"
)

func (this *DbcurdController) IndexAction() {
	this.Tplname = "views/admin/" + this.Masterdb + "/index.html"
	this.Rander()
}

func (this *DbcurdController) Get_page_listAction() {
	db := lib.NewQuerybuilder()
	where := db.Tbname(this.Masterdb).Get_where_data(this.Postdata)
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		if(this.T_data["order_str"]!="") {
			order_str = this.T_data["order_str"]+" "+this.T_data["order_type"]
		}
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	if (this.R.FormValue("where") != "") {
		if (where == "") {
			where = this.R.FormValue("where")
		} else {
			where += "and (" + this.R.FormValue("where") + ")"
		}
	}
	db.Dbinit()
	ct := db.Tbname(this.Masterdb).Where(where).Count()
	db.Dbinit()
	data := db.Tbname(this.Masterdb).Where(where).Order(order_str).Limit(limit_str).Select()
	for k,v:=range data{
		data[k]=this.Set_page_list(v)
	}
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

func (this *DbcurdController) Get_listAction() {
	db := lib.NewQuerybuilder()
	where := db.Tbname(this.Masterdb).Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		if(this.T_data["order_str"]!="") {
			order_str = this.T_data["order_str"]+" "+this.T_data["order_type"]
		}
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	if (this.R.FormValue("where") != "") {
		if (where == "") {
			where = this.R.FormValue("where")
		} else {
			where += "and (" + this.R.FormValue("where") + ")"
		}
	}
	db.Dbinit()
	data := db.Tbname(this.Masterdb).Where(where).Order(order_str).Select()
	db.Dbinit()
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

func (this *DbcurdController) Get_oneAction() {
	db := lib.NewQuerybuilder()
	where := db.Tbname(this.Masterdb).Get_where_data(this.Postdata)
	db.Dbinit()
	data := db.Tbname(this.Masterdb).Where(where).Find()
	list := make(map[string]interface{})
	list["data"] = data
	list["status"] = 100
	list["code"] = 0
	list["msg"] = ""
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

func (this *DbcurdController) EditAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"

		data = db.Tbname(this.Masterdb).Where(db.Get_key_eq_value(this.R.FormValue("id"))).Find()
	} else {
		data = db.Tbname(this.Masterdb).Get_new_add()
	}
	this.Data["edit_flag"] = edit_flag

	this.Data["data"] = data
	this.Tplname = "views/admin/" + this.Masterdb + "/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *DbcurdController) SaveAction() {
	db := lib.NewQuerybuilder()
	msg := "保存失败"
	var err error
	if (this.R.FormValue("id") != "") {
		_, err = db.Tbname(this.Masterdb).Where(db.Get_key_eq_value(this.R.FormValue("id"))).Update(this.Postdata)
	} else {
		_, err = db.Tbname(this.Masterdb).Insert(this.Postdata)

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

func (this *DbcurdController) DelAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname(this.Masterdb).Where(db.Get_key_in_value(this.R.FormValue("id"))).Delete()
	if (err == nil) {
		jsonstr.Msg = "删除成功"
		jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}


func (this *DbcurdController) Set_page_list(data map[string]string)map[string]string{
	db:=lib.NewQuerybuilder()
	switch this.Masterdb {
	case "":
		db.Tbname("").Where("").Find()
	}
	return data
}