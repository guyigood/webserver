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


func (this *AdminController) AdminGoodscategory_indexAction() {
	this.Tplname = "views/admin/goodscategory/index.html"
	this.Rander()
}

func (this *AdminController) AdminGoodscategory_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page-1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("goods_category").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
    if (order_str == "") {
    	order_str = "id desc"
    }else{
    	order_str += " "+this.R.FormValue("order_type")
    }
	ct := db.Tbname("goods_category").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("goods_category").Where(where).Order(order_str).Limit(limit_str).Select()
	for k,v:=range data{
		db.Dbinit()
		if(v["parent_id"]=="0"){
			data[k]["parent_id_name"]="顶级菜单"
		}
		list:=db.Tbname("goods_category").Where("id="+v["parent_id"]).Find()
		data[k]["parent_id_name"]=list["name"]

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
func (this *AdminController) AdminGoodscategory_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data:=make(map[string]string)
	edit_flag:=""
	if (this.R.FormValue("id") !=""){
	    edit_flag="1"
		data = db.Tbname("goods_category").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
	} else {
		data = db.Tbname("goods_category").Get_new_add()
	}
	if(data["parent_id"]==""){
		data["parent_id"]="0"
		data["parent_id_name"]="顶级菜单"
	}else{
		db.Dbinit()
		list:=db.Tbname("goods_category").Where("id="+data["parent_id"]).Find()
		data["parent_id_name"]=list["name"]
	}
	this.Data["edit_flag"]=edit_flag
	
	this.Data["data"]=datatype.MapString2interface(data)
	this.Tplname = "views/admin/goodscategory/edit.html"
	this.Rander()
}

func (this *AdminController) AdminShow_parent_menuAction(){
  this.Tplname="views/admin/goodscategory/pop.html"
  this.Rander()
}

func(this *AdminController)AdminCate_get_parent_menuAction(){
	list:=make([]map[string]interface{},0)
	data:=make(map[string]interface{})
	data["name"]="顶级分类"
	data["id"]="0"
	data["spread"]=true;
	data["children"]= this.build_cate_menu("0")
	list=append(list,data)
	b, _ := json.Marshal(&list)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}

func(this *AdminController)build_cate_menu(p_id string)[]map[string]interface{}{
	db:=lib.NewQuerybuilder()
	data:=db.Tbname("goods_category").Where(fmt.Sprintf("parent_id=%v",p_id)).Select()
	list:=make([]map[string]interface{},0)
	if(len(data)>0) {
		for _,v:=range data{
			tmp:=make(map[string]interface{})
			tmp["name"]=v["name"]
			tmp["id"]=v["id"]
			tmp["spread"]=false
			tmp_list:=this.build_cate_menu(v["id"])
			if(len(tmp_list)>0){
				tmp["children"]=tmp_list
				list = append(list, tmp)
			}else {
				if(this.R.FormValue("name")!=""){
				    if(strings.Contains(v["name"],this.R.FormValue("name"))){
				    	list = append(list, tmp)
					}
				}else {
					list = append(list, tmp)
				}
			}
		}
		return list
	}else{
		return nil
	}

}


/*递归查询上级*/
func (this *AdminController) AdminGoodscategory_getpathAction(id string) (string, int) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("goods_category").Where("id=" + id).Find()
	if (len(data) > 0) {
		return data["parent_id_path"], datatype.Type2int(data["level"])
	}
	return "0", 0
}

/* 保存记录Api */
func (this *AdminController) AdminGoodscategory_saveAction() {
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
	c_path, c_level := this.AdminGoodscategory_getpathAction(datatype.Type2str(this.Postdata["parent_id"]))
	if (c_level == 0) {
		c_level = 1
	} else {
		c_level++
	}

	var err error
	if (this.R.FormValue("id") != "") {
		Id := this.R.FormValue("id")
		this.Postdata["parent_id_path"] = c_path + "_" + this.R.FormValue("id")
		this.Postdata["level"] = c_level
		_, err = db.Tbname("goods_category").Where("id=" + Id).Update(this.Postdata)
	} else {
		sqlre, _:= db.Tbname("goods_category").Insert(this.Postdata)
		cat_id, _ := sqlre.LastInsertId()
		c_path += fmt.Sprintf("_%v", cat_id)
		db.Dbinit()
		db.Tbname("goods_category").Where(fmt.Sprintf("id=%v", cat_id)).Update(map[string]interface{}{"parent_id_path": c_path, "level": c_level})

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

func (this *AdminController) AdminGoodscategory_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	ct:=db.Tbname("goods").Where("cat_id="+this.R.FormValue("id")).Count()
	db.Dbinit()
	if(ct>0) {
		jsonstr.Msg = "已在使用中，不能删除"
	}else {
		_, err := db.Tbname("goods_category").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
		if (err == nil) {
			jsonstr.Msg = "删除成功"
			jsonstr.Status = 100
		}
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}