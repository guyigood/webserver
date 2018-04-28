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
	specItem := db.Tbname("spec_item").Where("id in (" + this.R.FormValue("spec_arr") + ")").Select()
	db.Dbinit()
	//goods_spec := db.Tbname("spec_goods_price").Where("goods_id=" + this.R.FormValue("goods_id")).Select()
	str1 := "<table class='layui-table' id='spec_input_tab'>"
	str1 += "<tr>"
	spec_name := ""
	for _, v := range specItem {
		db.Dbinit()

		spec := db.Tbname("spec").Where("id=" + v["spec_id"]).Find()
		if (spec_name != spec["name"]) {
			str1 += " <td><b>" + spec["name"] + "</b></td>";
			spec_name = spec["name"]
		}
	}
	str1 += `<td><b>价格</b></td>
	<td><b>库存</b></td>
	<td><b>SKU</b></td>
	</tr>`
	db.Dbinit()
	spec_group := db.Query("select spec_id from " + lib.Db_perfix + "spec_item Where id in (" + this.R.FormValue("spec_arr") + ") group by spec_id")
	cargo_sql := "select * from "
	for k, v := range spec_group {
		if (k == 0) {
			cargo_sql += fmt.Sprintf("(select id,item from %vspec_item where spec_id=%v and id in ("+this.R.FormValue("spec_arr")+")) as a%v ", lib.Db_perfix, v["spec_id"], k)
		} else {
			cargo_sql += fmt.Sprintf(" CROSS join (select id as id%v,item as item%v from %vspec_item where spec_id=%v and id in ("+this.R.FormValue("spec_arr")+")) as a%v ", k, k, lib.Db_perfix, v["spec_id"], k)
		}

	}
	spec_ct := len(spec_group) - 1
	db.Dbinit()
	//fmt.Println(cargo_sql)
	cargo_data := db.Query(cargo_sql)
	key := ""
	key_name := ""
	///fmt.Println(cargo_sql)
	//fmt.Println(cargo_data)
	for _, v := range cargo_data {
		key = v["id"]
		key_name = v["item"]
		str1 += "<tr>"
		if (spec_ct > 0) {
			for i := 0; i < spec_ct; i++ {
				j := i + 1;
				key += "_" + v[fmt.Sprintf("id%v", j)]
				key_name += "_" + v[fmt.Sprintf("item%v", j)]
				str1 += "<td>" + v[fmt.Sprintf("item%v", j)] + "</td>"
			}
		}
		str1 += "<td>" + key_name + "</td>"
		db.Dbinit()

		s_price_data := db.Tbname("spec_goods_price").Where(fmt.Sprintf("key='%v' and goods_id=%v", key, this.R.FormValue("goods_id"))).Find()
		s_price := "0"
		s_price_sku := ""
		s_price_store := "0"
		if (s_price_data != nil) {
			s_price = s_price_data["price"]
			s_price_sku = s_price_data["sku"]
			s_price_store = s_price_data["store_count"]
		}

		str1 += fmt.Sprintf("<input type=hidden name=items_key value='%v'>", key)
		str1 += fmt.Sprintf("<td><input name='items_price[]' value='%v' class='layui-input'  /></td>", s_price)
		str1 += fmt.Sprintf("<td><input name='item_store[]' value='%v' class='layui-input'/></td>", s_price_store)
		str1 += fmt.Sprintf("<td><input name='item_sku[]' value='%v' class='layui-input'/>", s_price_sku)
		str1 += fmt.Sprintf("<input type='hidden' name='item_name[]' value='%v' /></td>", key_name)
		str1 += "</tr>"
	}
	this.W.Write([]byte(str1))

}
