package module

import (
	"gylib/dbfun"
	"strings"
)

type Orgstruct struct {
}

func NewOrgStruct() (*Orgstruct) {
	this := new(Orgstruct)
	return this
}

func (this *Orgstruct) Get_City_by_id(id string) (string) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("china").Where("id=" + id).Find()
	return data["name"]
}

func (this *Orgstruct) Get_City_by_map(list map[string]string) (string) {
	name := this.Get_City_by_id(list["province"])
	name += this.Get_City_by_id(list["city"])
	name += this.Get_City_by_id(list["area"])
	return name
}

func (this *Orgstruct) Get_china_data(search_str string) ([]map[string]interface{}) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("china").Where("pid=0 and open=11").Select()
	list := make([]map[string]interface{}, 0)
	for _, v := range data {
		tmp := make(map[string]interface{})
		tmp["province"] = ""
		tmp["city"] = ""
		tmp["area"] = ""
		tmp["full_name"] = ""
		tmp["province"] = v["id"]
		tmp["full_name"] = v["name"]
		tmp["name"] = v["name"]
		tmp["id"] = "0"
		tmp["spread"] = true;
		t_data:=this.Get_china_list(v["id"], search_str)
		tmp["children"] = t_data
		if(search_str!=""){
			if(len(t_data)>0){
				list = append(list, tmp)
			}
		}else {
			list = append(list, tmp)
		}

	}
	return list
}

func (this *Orgstruct) Get_china_list(p_id string, search_str string) ([]map[string]interface{}) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("china").Where("open=11 and pid=" + p_id).Select()
	list := make([]map[string]interface{}, 0)
	if (len(data) > 0) {
		for _, v := range data {
			tmp := make(map[string]interface{})
			tmp["name"] = v["name"]
			tmp["id"] = v["id"]
			tmp["province"] = ""
			tmp["city"] = ""
			tmp["area"] = ""
			tmp["full_name"] = ""
			if (v["level"] == "1") {

				tmp["province"] = v["id"]
				tmp["full_name"] = v["name"]
			} else if (v["level"] == "2") {
				tmp["city"] = v["id"]
				ca1 := this.Get_china_pid(v["pid"])
				tmp["province"] = ca1["id"]
				tmp["full_name"] = ca1["name"] + v["name"]
			} else {
				tmp["area"] = v["id"]
				ca1 := this.Get_china_pid(v["pid"])
				tmp["city"] = ca1["id"]
				ca2 := this.Get_china_pid(v["pid"])
				tmp["province"] = ca2["id"]
				tmp["full_name"] = ca2["name"] + ca1["name"] + v["name"]
			}
			tmp["spread"] = false
			tmp_list := this.Get_china_list(v["id"], search_str)
			if (len(tmp_list) > 0) {
				tmp["children"] = tmp_list
				list = append(list, tmp)
			} else {
				if (search_str != "") {
					if (strings.Contains(v["name"], search_str)) {
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

func (this *Orgstruct) Get_china_pid(id string) (map[string]string) {
	db := lib.NewQuerybuilder()
	return db.Tbname("china").Where("id=" + id).Find()

}

func (this *Orgstruct) Get_china_tree() ([]map[string]interface{}) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("china").Where("pid=0").Select()
	list := make([]map[string]interface{}, 0)
	for _, v := range data {
		tmp := make(map[string]interface{})
		tmp["value"] = v["id"]
		tmp["title"] = v["name"]
		t_data:=this.Get_china_list_tree(v["id"])
		tmp["data"] = t_data
		list = append(list, tmp)

	}
	return list
}
func (this *Orgstruct) Get_china_list_tree(p_id string) ([]map[string]interface{}) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("china").Where("pid=" + p_id).Select()
	list := make([]map[string]interface{}, 0)
	if (len(data) > 0) {
		for _, v := range data {
			tmp := make(map[string]interface{})
			tmp["title"] = v["name"]
			tmp["value"] = v["id"]
			if(v["open"]=="11") {
				tmp["checked"] = true
			}
			tmp_data:=make([]map[string]interface{},0)
			tmp_data = this.Get_china_list_tree(v["id"])
			if(len(tmp_data)>0) {
				tmp["data"] = tmp_data
			}else{
				tmp["data"]=make([]map[string]interface{},0)
			}
			list = append(list, tmp)
		}
		return list
	} else {
		return nil
	}
}
