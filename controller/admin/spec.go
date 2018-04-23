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

func (this *AdminController) AdminSpec_indexAction() {
	this.Tplname = "views/admin/spec/index.html"
	db := lib.NewQuerybuilder()
	db.Dbinit()
	this.Data["goods_type"] = db.Tbname("goods_type").Where("").Select()
	this.Rander()
}

func (this *AdminController) AdminSpec_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("spec").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	ct := db.Tbname("spec").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("spec").Where(where).Order(order_str).Limit(limit_str).Select()
	for k,v:=range data{
		db.Dbinit()
		list:=db.Tbname("goods_type").Where("id="+v["type_id"]).Find()
		data[k]["type_id"]=list["name"]
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
func (this *AdminController) AdminSpec_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"
		data = db.Tbname("spec").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
		db.Dbinit()
		list:=db.Tbname("spec_item").Where("spec_id="+data["id"]).Select()
		data["memo"]=""
		for _,v:=range list {
			data["memo"] +=v["item"]+"\n"
		}

	} else {
		data = db.Tbname("spec").Get_new_add()
		data["memo"]=""
	}

	this.Data["edit_flag"] = edit_flag
	db.Dbinit()
	this.Data["goods_type"] = db.Tbname("goods_type").Where("").Select()
	this.Data["data"] = datatype.MapString2interface(data)
	this.Tplname = "views/admin/spec/edit.html"
	this.Rander()
}

/* 保存记录Api */
func (this *AdminController) AdminSpec_saveAction() {
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
	id:=""
	var err error
	if (this.R.FormValue("id") != "") {
		id = this.R.FormValue("id")
		_, err = db.Tbname("spec").Where("id=" + id).Update(this.Postdata)

	} else {
		sqlre, _:= db.Tbname("spec").Insert(this.Postdata)
		id64,_:=sqlre.LastInsertId()
		id=strconv.Itoa(datatype.Int64toint(id64))
	}


   /*获取items数据，去掉为空的数据*/
	spec_item:=strings.Split(this.R.FormValue("memo"),"\n")
	post_items:=make([]string,0)
	for _,val:=range spec_item {
		val = strings.TrimSpace(val)
		if (val != "") {
			post_items = append(post_items, val)
		}
	}

	/*检查不存在的数据,插入新纪录*/
	db.Dbinit()
	db_items := db.Tbname("spec_item").Where("spec_id=" + id).Select()

	for _,val:=range post_items{
		flag:=true
		for _,d_val:=range db_items{
		   if(val==d_val["item"]){
		   	 flag=false
		   	 break
		   }
		}
		if(flag){
			db.Dbinit()
			db.Tbname("spec_item").Insert(map[string]interface{}{"spec_id":id,"item":val})
		}
	}
	/*删除数据库中间不存在的数据*/
	for _,val:=range db_items{
		flag:=true
		for _,d_val:=range post_items{
			if(val["item"]==d_val){
				flag=false
				break
			}
		}
		if(flag){
			db.Dbinit()
			db.Tbname("spec_goods_price").Where("`key` = '"+val["id"]+"'").Delete(); // 删除规格项价格表
			db.Dbinit()
			db.Tbname("spec_item").Where("id="+val["id"]).Delete()
		}
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

func (this *AdminController) AdminSpec_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	ct:=db.Tbname("goods_spec").Where("spec_type="+this.R.FormValue("id")).Count()
	db.Dbinit()
	if(ct>0) {
		jsonstr.Msg = "已在使用中，不能删除"
	}else {
		_, err := db.Tbname("spec").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
		if (err == nil) {
			jsonstr.Msg = "删除成功"
			jsonstr.Status = 100
		}
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}
