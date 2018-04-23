package admin

import (
	"strconv"
	"fmt"
	"gylib/common"
	"gylib/dbfun"
	"encoding/json"
	"strings"
)

func (this *AdminController) AdminNav_indexAction() {
	this.Tplname = "views/admin/nav/index.html"
	this.Rander()
}

func (this *AdminController) AdminNav_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("nav").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "parent_id"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	ct := db.Tbname("nav").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("nav").Where(where).Order(order_str).Limit(limit_str).Select()
	for k, v := range data {
		if (v["parent_id"] == "0") {
			data[k]["parent_id"] = "顶级菜单"
		} else {
			db.Dbinit()
			list := db.Tbname("nav").Where("id=" + v["parent_id"]).Find()
			data[k]["parent_id"] = list["nav_name"]
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
func (this *AdminController) AdminNav_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"
		data = db.Tbname("nav").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
	} else {
		data = db.Tbname("nav").Get_new_add()
	}
	if (data["parent_id"] == "") {
		data["parent_id"] = "0"
		data["parent_id_name"] = "顶级菜单"
	} else if (data["parent_id"] == "0") {
		data["parent_id"] = "0"
		data["parent_id_name"] = "顶级菜单"
	} else {
		db.Dbinit()
		list := db.Tbname("nav").Where("id=" + data["parent_id"]).Find()
		data["parent_id_name"] = list["name"]
	}
	this.Data["edit_flag"] = edit_flag
	db.Dbinit()
	this.Data["nav"] = db.Tbname("nav").Where("").Select()
	this.Data["data"] = data
	this.Tplname = "views/admin/nav/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminNav_saveAction() {
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
		_, err = db.Tbname("nav").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("nav").Insert(this.Postdata)

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

func (this *AdminController) AdminNav_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("nav").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
	if (err == nil) {
		jsonstr.Msg = "删除成功"
		jsonstr.Status = 100
	}
	//fmt.Println(db.GetLastSql())
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}

func (this *AdminController) AdminNav_menuAction() {
	qxz_list := make(map[string]string)
	if (this.R.FormValue("qxz_id") != "") {
		db := lib.NewQuerybuilder()
		qxz_list = db.Tbname("qxz").Where("id=" + this.R.FormValue("qxz_id")).Find()
	}

	data := this.Get_Build_menu(0, qxz_list)
	b, _ := json.Marshal(&data)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}

func (this *AdminController) check_menu_flag(qxz_list map[string]string, data map[string]string) (bool) {
	if (len(qxz_list) == 0) {
		return false
	}
	qxz := strings.Split(qxz_list["memo"], ",")
	flag := false
	for _, v := range qxz {
		if (v == data["nav_code"]) {
			flag = true
			break
		}
	}
	return flag
}

func (this *AdminController) Get_Build_menu(p_id int, qxz_list map[string]string) ([]map[string]interface{}) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("nav").Where(fmt.Sprintf("parent_id=%v", p_id)).Select()
	list := make([]map[string]interface{}, 0)
	if (len(data) > 0) {
		for _, v := range data {
			tmp := make(map[string]interface{})
			tmp["title"] = v["nav_name"]
			tmp["id"] = v["nav_code"]
			id, _ := strconv.Atoi(v["id"])
			tmp["is_check"] = this.check_menu_flag(qxz_list, v)
			tmp_list := this.Get_Build_menu(id, qxz_list)
			if (len(tmp_list) > 0) {
				tmp["children"] = tmp_list
			}
			list = append(list, tmp)
		}
		return list
	} else {
		return nil
	}
}

func (this *AdminController) AdminShow_nav_parent_menuAction() {
	this.Tplname = "views/admin/nav/pop.html"
	this.Rander()
}

func (this *AdminController) AdminNav_get_parent_menuAction() {
	list := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["name"] = "顶级菜单"
	data["id"] = "0"
	data["spread"] = true;
	data["children"] = this.build_nav_menu("0")
	list = append(list, data)
	b, _ := json.Marshal(&list)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}

func (this *AdminController) build_nav_menu(p_id string) []map[string]interface{} {
	db := lib.NewQuerybuilder()
	data := db.Tbname("nav").Where(fmt.Sprintf("parent_id=%v", p_id)).Select()
	list := make([]map[string]interface{}, 0)
	if (len(data) > 0) {
		for _, v := range data {
			tmp := make(map[string]interface{})
			tmp["name"] = v["nav_name"]
			tmp["id"] = v["id"]
			tmp["spread"] = false
			tmp_list := this.build_nav_menu(v["id"])
			if (len(tmp_list) > 0) {
				tmp["children"] = tmp_list
				list = append(list, tmp)
			} else {
				if (this.R.FormValue("name") != "") {
					if (strings.Contains(v["nav_name"], this.R.FormValue("name"))) {
						list = append(list, tmp)
					}
				} else {
					list = append(list, tmp)
				}
			}
		}
		return list
	} else {
		return nil
	}

}
