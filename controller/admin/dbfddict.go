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

func (this *AdminController) AdminDbfddict_indexAction() {
	this.Tplname = "views/admin/dbfddict/index.html"
	t_id := this.R.FormValue("id")
	this.Data["t_id"] = t_id
	db := lib.NewQuerybuilder()
	tb_list := db.Tbname("db_tb_dict").Where("id=" + t_id).Find()
	db.Dbinit()
	rows := db.Query("SHOW full COLUMNS FROM " + tb_list["name"])
	for k, v := range rows {
		db.Dbinit()
		mx := db.Tbname("db_fd_dict").Where(fmt.Sprintf("t_id=%v and name='%v'", t_id, v["Field"])).Find()
		if (len(mx) == 0) {
			db.Dbinit()
			s_data := make(map[string]interface{})
			if (v["Key"] == "PRI" && v["Extra"] == "auto_increment") {
				s_data["f_is_display"] = 0
				s_data["f_ed_display"] = 0
				s_data["f_ed_display"] = 0
				s_data["f_add_display"] = 0
				s_data["f_readonly"] = 1
			} else {
				s_data["f_is_display"] = 1
				s_data["f_ed_display"] = 1
				s_data["f_ed_display"] = 1
				s_data["f_add_display"] = 1
			}
			s_data["name"] = v["Field"]
			if (v["Comment"] != "") {
				s_data["cn_name"] = v["Comment"]
			} else {
				s_data["cn_name"] = v["Field"]
			}
			s_data["t_id"] = t_id
			s_data["f_list_order"] = k
			s_data["f_ed_order"] = k
			s_data["f_readonly"] = 0
			s_data["f_type"] = "单行文本"
			s_data["f_check"] = "不验证"
			s_data["f_isnull"] = 0
			s_data["list_tb_name"] = 0
			db.Tbname("db_fd_dict").Insert(s_data)
		}
	}
	this.Rander()
}

func (this *AdminController) AdminDbfddict_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("db_fd_dict").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	if (this.R.FormValue("t_id") != "") {
		if (where != "") {
			where += " and t_id=" + this.R.FormValue("t_id")
		} else {
			where = "t_id="+this.R.FormValue("t_id")
		}
	}

	ct := db.Tbname("db_fd_dict").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("db_fd_dict").Where(where).Order(order_str).Limit(limit_str).Select()
	//fmt.Println(this.Postdata)
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

/* 获取当条记录Api */
func (this *AdminController) AdminDbfddict_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"
		data = db.Tbname("db_fd_dict").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
	} else {
		data = db.Tbname("db_fd_dict").Get_new_add()
	}
	this.Data["edit_flag"] = edit_flag
	this.Data["data"] = datatype.MapString2interface(data)
	this.Tplname = "views/admin/dbfddict/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminDbfddict_saveAction() {
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
		_, err = db.Tbname("db_fd_dict").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("db_fd_dict").Insert(this.Postdata)

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

func (this *AdminController) AdminDbfddict_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("db_fd_dict").Where(fmt.Sprintf("id in(%v)", this.R.FormValue("id"))).Delete()
	if (err != nil) {
		jsonstr.Msg = "删除成功"
		jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}
