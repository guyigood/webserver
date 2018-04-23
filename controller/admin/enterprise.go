package admin

import (
	"strconv"
   	"fmt"
   	"gylib/common"
   	"gylib/dbfun"
   	"encoding/json"
   	"strings"
)


func (this *AdminController) AdminEnterprise_indexAction() {
	this.Tplname = "views/admin/enterprise/index.html"

	this.Rander()
}

func (this *AdminController) AdminEnterprise_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page-1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("enterprise").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
    if (order_str == "") {
    	order_str = "id desc"
    }else{
    	order_str += " "+this.R.FormValue("order_type")
    }
	ct := db.Tbname("enterprise").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("enterprise").Where(where).Order(order_str).Limit(limit_str).Select()
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

/* 获取当条记录Api */
func (this *AdminController) AdminEnterprise_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data:=make(map[string]string)
	edit_flag:="1"
	    data = db.Tbname("enterprise").Find()
	this.Data["edit_flag"]=edit_flag
	
	this.Data["data"]=data
	this.Tplname = "views/admin/enterprise/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminEnterprise_saveAction() {
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
		_, err = db.Tbname("enterprise").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("enterprise").Insert(this.Postdata)

	}
	jsonstr := new(Json_msg)
	if (err == nil) {
		jsonstr.Status = 100
		msg="保存成功"
	} else {
		jsonstr.Status = 101
	}
	jsonstr.Msg = msg
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)

}

func (this *AdminController) AdminEnterprise_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("enterprise").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
	if(err==nil){
	    jsonstr.Msg = "删除成功"
	    jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}