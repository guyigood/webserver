package admin

import (
	"strconv"
   	"fmt"
   	"gylib/common"
   	"gylib/dbfun"
   	"encoding/json"
   	"strings"
)


func (this *AdminController) AdminDepart_indexAction() {
	this.Tplname = "views/admin/depart/index.html"
	this.Rander()
}

func (this *AdminController) AdminDepart_treeAction(){
	list:=make([]map[string]interface{},0)
	data:=make(map[string]interface{})
	db:=lib.NewQuerybuilder()
	e_data:=db.Tbname("enterprise").Find()
	data["name"]=e_data["name"]
	data["id"]="0"
	data["spread"]=true;
	data["children"]= this.Get_dept_tree(0)
	list=append(list,data)
	b, _ := json.Marshal(&list)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}


func (this *AdminController) Get_dept_tree(p_id int) ([]map[string]interface{}) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("depart").Where(fmt.Sprintf("forid=%v", p_id)).Select()
	list := make([]map[string]interface{}, 0)
	if (len(data) > 0) {
		for _, v := range data {
			tmp := make(map[string]interface{})
			tmp["name"]=v["name"]
			tmp["id"]=v["id"]
			tmp["spread"]=false
			id, _ := strconv.Atoi(v["id"])
			tmp_list := this.Get_dept_tree(id)
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

func (this *AdminController) AdminDepart_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page-1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("depart").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
    if (order_str == "") {
    	order_str = "id desc"
    }else{
    	order_str += " "+this.R.FormValue("order_type")
    }
    if(where==""){
    	where="forid=0"
	}
	ct := db.Tbname("depart").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("depart").Where(where).Order(order_str).Limit(limit_str).Select()
	db.Dbinit()
	e_data:=db.Tbname("enterprise").Find()
	for k,v:=range data{
		if(v["forid"]=="0"){
			data[k]["forid"]=e_data["name"]
		}else {
			db.Dbinit()
			d_list := db.Tbname("depart").Where("id="+v["forid"]).Find()
			data[k]["forid"]=d_list["name"]
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
func (this *AdminController) AdminDepart_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data:=make(map[string]string)
	edit_flag:=""
	if (this.R.FormValue("id") !=""){
	    edit_flag="1"
		data = db.Tbname("depart").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
	} else {
		data = db.Tbname("depart").Get_new_add()
		data["forid"]=this.R.FormValue("forid")
	}
	this.Data["edit_flag"]=edit_flag
	
	this.Data["data"]=data
	this.Tplname = "views/admin/depart/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminDepart_saveAction() {
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
		_, err = db.Tbname("depart").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("depart").Insert(this.Postdata)

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

func (this *AdminController) AdminDepart_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("depart").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
	if(err==nil){
	    jsonstr.Msg = "删除成功"
	    jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}