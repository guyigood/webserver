package admin

import (
	"strconv"
	"fmt"
	"gylib/dbfun"
	"encoding/json"
	"gylib/common/datatype"
	"strings"
	"websever/public/module"
)

func (this *AdminController) AdminModule_indexAction() {
	db := lib.NewQuerybuilder()
	rows := db.Query("show tables")
	for _, val := range rows {
		fd_name := strings.Title("tables_in_" + lib.Db_Struct.Db_name);
		tbname := val[fd_name]
		db.Dbinit()
		list := db.Tbname("db_tb_dict").Where("name='" + tbname + "'").Find()
		if (len(list) == 0) {
			db.Dbinit()
			db.Tbname("db_tb_dict").Insert(map[string]interface{}{"name": tbname, "cn_name": tbname})
		}
	}
	this.Tplname = "views/admin/module/index.html"
	/*tpl := make([]string, 0)
	tpl = append(tpl, "views/public/list_head.html")
	tpl = append(tpl, "views/public/list_foot.html")
	this.MuitplRander(tpl...)*/
	this.Rander()
}

func (this *AdminController) AdminModule_buildAction() {
	if (this.R.Method == "POST") {
		gofile := module.NewBuildGo()
		gofile.Build_all(this.R.FormValue("id"), this.R.FormValue("lb"))
		jsonstr := new(Json_msg)
		jsonstr.Status = 100
		jsonstr.Msg = "模板生成成功"
		b, _ := json.Marshal(&jsonstr)
		this.W.Write(b)
	} else {
		this.Data["id"] = this.R.FormValue("id")
		this.Tplname = "views/admin/module/build.html"
		this.Rander()
	}
}

func (this *AdminController) AdminModule_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("db_tb_dict").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "name asc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	ct := db.Tbname("db_tb_dict").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("db_tb_dict").Where(where).Order(order_str).Limit(limit_str).Select()
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

func (this *AdminController) AdminModule_getallAction() {
	db := lib.NewQuerybuilder()
	where := db.Tbname("db_tb_dict").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	db.Dbinit()
	data := db.Tbname("db_tb_dict").Where(where).Order(order_str).Select()
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

/* 获取当条记录Api */
func (this *AdminController) AdminModule_editAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "text/html")
	id := this.R.FormValue("id")
	edit_flag := "";
	if (id != "") {
		edit_flag = "1"
	}
	data := make(map[string]string)
	if (this.R.FormValue("id") != "") {
		data = db.Tbname("db_tb_dict").Where(fmt.Sprintf("id=%v", id)).Find()
	} else {
		data = db.Tbname("db_tb_dict").Get_new_add()
	}
	this.Data["edit_flag"] = edit_flag
	this.Data["data"] = datatype.MapString2interface(data)
	this.Tplname = "views/admin/module/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminModule_saveAction() {
	this.W.Header().Set("content-type", "application/json")
	db := lib.NewQuerybuilder()
	msg := "保存失败"
	var err error
	if (this.R.FormValue("id") != "") {
		Id := this.R.FormValue("id")
		_, err = db.Tbname("db_tb_dict").Where("id=" + Id).Update(this.Postdata)

	} else {
		_, err = db.Tbname("db_tb_dict").Insert(this.Postdata)

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

func (this *AdminController) AdminModule_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("db_tb_dict").Where(fmt.Sprintf("id in(%v)", this.R.FormValue("id"))).Delete()
	if (err == nil) {
		jsonstr.Status = 100
		jsonstr.Msg = "删除成功"
	}

	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}
