package admin

import (
	"strconv"
	"fmt"
	"gylib/common"
	"gylib/dbfun"
	"encoding/json"
	"strings"
	"gylib/common/datatype"
)

func (this *AdminController) AdminGoods_indexAction() {
	this.Tplname = "views/admin/goods/index.html"
	db := lib.NewQuerybuilder()
	db.Dbinit()
	this.Data["goodscategory"] = db.Tbname("goods_category").Select()
	this.Rander()
}

func (this *AdminController) AdminGoods_listAction() {
	db := lib.NewQuerybuilder()
	get_page := this.R.FormValue("page")
	page := 1
	if ( get_page != "") {
		page, _ = strconv.Atoi(get_page)
	}
	limit, _ := strconv.Atoi(this.R.FormValue("limit"))
	Page_no := (page - 1) * limit
	limit_str := fmt.Sprintf("%v,%v", Page_no, limit)
	where := db.Tbname("goods").Get_where_data(this.Postdata)
	order_str := this.R.FormValue("order_str")
	if (order_str == "") {
		order_str = "id desc"
	} else {
		order_str += " " + this.R.FormValue("order_type")
	}
	ct := db.Tbname("goods").Where(where).Count()
	db.Dbinit()
	data := db.Tbname("goods").Where(where).Order(order_str).Limit(limit_str).Select()
	list := make(map[string]interface{})
	list["data"] = data
	list["code"] = 0
	list["msg"] = ""
	list["count"] = ct
	b, _ := json.Marshal(&list)
	this.W.Write(b)
}

/* 获取当条记录Api */
func (this *AdminController) AdminGoods_editAction() {
	db := lib.NewQuerybuilder()
	//this.W.Header().Set("content-type", "text/html")
	data := make(map[string]string)
	edit_flag := ""
	if (this.R.FormValue("id") != "") {
		edit_flag = "1"
		data = db.Tbname("goods").Where(fmt.Sprintf("id=%v", this.R.FormValue("id"))).Find()
		db.Dbinit()
		cate := db.Tbname("goods_category").Where("id=" + data["cat_id"]).Find()
		data["cat_id_name"] = cate["name"]
	} else {
		data = db.Tbname("goods").Get_new_add()
		data["cat_id_name"] = ""
	}
	this.Data["edit_flag"] = edit_flag
	db.Dbinit()
	this.Data["goodstype"] = db.Tbname("goods_type").Where("").Select()
	db.Dbinit()
	this.Data["spec"] = db.Tbname("spec").Where("").Select()
	this.Data["data"] = data
	this.Tplname = "views/admin/goods/edit.html"
	this.MuitplRander("views/admin/goods/base.html", "views/admin/goods/goods_attr.html", "views/admin/goods/goods_img.html", "views/admin/goods/goods_type.html")
}

/* 保存记录Api */
func (this *AdminController) AdminGoods_saveAction() {
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
		_, err = db.Tbname("goods").Where("id=" + Id).Update(this.Postdata)
	} else {
		_, err = db.Tbname("goods").Insert(this.Postdata)

	}
	jsonstr := new(Json_msg)
	if (err == nil) {
		jsonstr.Status = 100
		msg = "保存成功"
	} else {
		fmt.Println(err)
		fmt.Println(db.GetLastSql())
		jsonstr.Status = 101
	}

	jsonstr.Msg = msg
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)

}

func (this *AdminController) AdminGoods_delAction() {
	db := lib.NewQuerybuilder()
	this.W.Header().Set("content-type", "application/json")
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "删除失败"
	_, err := db.Tbname("goods").Where(fmt.Sprintf("id in (%v)", this.R.FormValue("id"))).Delete()
	if (err == nil) {
		jsonstr.Msg = "删除成功"
		jsonstr.Status = 100
	}
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}

func (this *AdminController) AdminGet_goods_specAction() {
	db := lib.NewQuerybuilder()
	data := db.Tbname("spec").Where("type_id=" + this.R.FormValue("type_id")).Select()
	goods_id := this.R.FormValue("goods_id")
	db.Dbinit()
	items := db.Query(fmt.Sprintf("select GROUP_CONCAT(`key` SEPARATOR '_') AS items_id from %vspec_goods_price where goods_id=%v ", lib.Db_perfix, goods_id))
	items_id := strings.Split(items[0]["items_id"], "_")
	list := make([]map[string]interface{}, 0)
	for _, v := range data {
		temp := make(map[string]interface{})
		temp = datatype.MapString2interface(v)
		temp["checked"] = false
		if (this.R.FormValue("goods_id") != "0") {
			temp["checked"] = common.Check_array_in(v["id"], items_id)
		}
		db.Dbinit()
		temp["data"] = db.Tbname("spec_item").Where("spec_id=" + v["id"]).Select()
		list = append(list, temp)
	}
	jsonstr := new(Json_msg)
	jsonstr.Status = 100
	jsonstr.Msg = ""
	jsonstr.Data = list
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)

}

func (this *AdminController) AdminGet_goods_spec_itemAction() {
	db := lib.NewQuerybuilder()
	specItem := db.Tbname("spec_item").Where("id in (" + this.R.FormValue("spec_arr")).Select()
	db.Dbinit()
	//goods_spec := db.Tbname("spec_goods_price").Where("goods_id=" + this.R.FormValue("goods_id")).Select()
	str1 := "<table class='layui-table' id='spec_input_tab'>"
	str1 += "<tr>"
	for _,v:=range specItem{
		db.Dbinit()
		spec := db.Tbname("spec").Where("id="+v["spec_id"]).Find()
		str1 += " <td><b>"+spec["name"]+"</b></td>";
	}
	str1 += `<td><b>价格</b></td>
	<td><b>库存</b></td>
	<td><b>SKU</b></td>
	</tr>`
	for _,v:=range specItem{

	}


}
