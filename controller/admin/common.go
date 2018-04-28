package admin

import (
	"net/http"
	"encoding/base64"
	"strconv"
	"time"
	"io/ioutil"
	"encoding/json"
	"gylib/common"
	"fmt"
	"os"
	"io"
	"regexp"
	"websever/public/module"
)

func get_ueditor_config() (string) {
	file, err := os.Open("conf/config.json")
	if err != nil {
		return "error"
	}
	defer file.Close()
	fd, err := ioutil.ReadAll(file)
	src := string(fd)
	re, _ := regexp.Compile("\\/\\*[\\S\\s]+?\\*\\/") //参考php的$CONFIG = json_decode(preg_replace("/\/\*[\s\S]+?\*\//", "", file_get_contents("config.json")), true);
	//将php中的正则移植到go中，需要将/ \/\*[\s\S]+?\*\/  /去掉前后的/，然后将\改成2个\\
	//参考//去除所有尖括号内的HTML代码，并换成换行符
	// re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	// src = re.ReplaceAllString(src, "\n")
	//当把<和>换成/*和*\时，斜杠/和*之间加双斜杠\\才行。
	src = re.ReplaceAllString(src, "")
	return src
}

func (this *AdminController) AdminCommon_popAction() {
	title := "名称"
	url_list := "/admin/nav_get_parent_menu/"
	input_box := `<input type='hidden' id='parent_id_name'>
          <input type='hidden' id='parent_id'>`
	set_js := `$('#parent_id').val(node.id);
         $('#parent_id_name').val(node.name);
         $('#label').text('当前选中:'+node.name);`
	close_js := `parent.layui.$('#parent_id').val($('#parent_id').val());
         parent.layui.$('#parent_id_name').val($('#parent_id_name').val());`

	switch this.R.FormValue("name") {
	case "dept":
		title = "部门名称"
		url_list = "/admin/depart_tree/"
		set_js = `
         var node_id=node.id;
         var node_name=node.name;
         if($('#parent_id').val()!=""){
              node_id=$('#parent_id').val()+","+node.id;
              node_name=$('#parent_id_name').val()+","+node.name;
          }
	      $('#parent_id').val(node_id);
          $('#parent_id_name').val(node_name);
		$('#label').text('当前选中:'+node_name);
		`
		close_js = `parent.layui.$('#dept_id').val($('#parent_id').val());
         parent.layui.$('#dept_name').val($('#parent_id_name').val());`
	case "goods":
		title = "分类选择"
		url_list = "/admin/cate_get_parent_menu/"
		set_js = `
         var node_id=node.id;
         var node_name=node.name;
         $('#parent_id').val(node_id);
         $('#parent_id_name').val(node_name);
		 $('#label').text('当前选中:'+node_name);
		`
		close_js = `parent.layui.$('#cat_id').val($('#parent_id').val());
         parent.layui.$('#cat_id_name').val($('#parent_id_name').val());`
	case "china":
		title = "城市选择"
		url_list = "/admin/common_get_china/"
		input_box = `<input type='hidden' id='province'>
           <input type='hidden' id='full_name'>
		  <input type='hidden' id='area'>
          <input type='hidden' id='city'>`
		set_js = `
         var node_id=node.id;
         var node_name=node.full_name;
         $('#province').val(node.province);
         $('#city').val(node.city);
         $('#area').val(node.area);
         $('#full_name').val(node.full_name);
         $('#label').text('当前选中:'+node_name);
		`
		close_js = `parent.layui.$('#province').val($('#province').val());
					parent.layui.$('#city').val($('#city').val());
					parent.layui.$('#area').val($('#area').val());
         parent.layui.$('#full_name').val($('#full_name').val());`
	}

	this.Data["title"] = title
	this.Data["list_url"] = url_list
	this.Data["set_js"] = set_js
	this.Data["hidden_input"] = input_box
	this.Data["close_js"] = close_js
	//fmt.Println(this.Data)
	this.Tplname = "views/admin/common/pop.html"
	this.Rander()
}

func (this *AdminController) AdminCommon_get_chinaAction() {
	//redis_ct := redispack.Get_redis_pool()
	//client := redis_ct.Get()
	//hasok, _ := client.Do("GET", "china")
	data := make([]map[string]interface{}, 0)
	china := module.NewOrgStruct()
	data = china.Get_china_data(this.R.FormValue("name"))
	//fmt.Println(data)
	/*fmt.Println(hasok)
	if (hasok ==nil) {
		china := module.NewOrgStruct()
		data = china.Get_china_data( this.R.FormValue("name"))
		fmt.Println("data", data)
		tmp, _ := json.Marshal(&data)
		client.Do("SETEX", "china", 3600, tmp)
	} else {
		if(hasok!=nil) {
			json.Unmarshal(hasok.([]byte), &data)
		}

	}*/
	raw, _ := json.Marshal(&data)
	this.W.Write(raw)
}

func (this *AdminController) AdminCommon_get_db_listAction(){

}

func (this *AdminController) AdminCommon_get_china_treeAction() {
	//redis_ct := redispack.Get_redis_pool()
	//client := redis_ct.Get()
	//hasok, _ := client.Do("GET", "china")
	data := make([]map[string]interface{}, 0)
	china := module.NewOrgStruct()
	data = china.Get_china_tree()
	//fmt.Println(data)
	/*fmt.Println(hasok)
	if (hasok ==nil) {
		china := module.NewOrgStruct()
		data = china.Get_china_data( this.R.FormValue("name"))
		fmt.Println("data", data)
		tmp, _ := json.Marshal(&data)
		client.Do("SETEX", "china", 3600, tmp)
	} else {
		if(hasok!=nil) {
			json.Unmarshal(hasok.([]byte), &data)
		}

	}*/
	raw, _ := json.Marshal(&data)
	this.W.Write(raw)
}


func (this *AdminController) AdminCommon_uploadAction() {
	op := this.R.FormValue("action")
	switch op {
	case "config":
		src := get_ueditor_config()
		tt := []byte(src)
		var r interface{}
		json.Unmarshal(tt, &r) //这个byte要解码
		conf, _ := json.Marshal(&r)
		this.W.Write(conf)
	case "uploadimage", "uploadfile", "uploadvideo":
		baidu_json := save_baidu_file(this.R)
		b, _ := json.Marshal(&baidu_json)
		this.W.Write(b)
	case "uploadscrawl":
		baidu_json := save_baidu_base64(this.R)
		b, _ := json.Marshal(&baidu_json)
		this.W.Write(b)
	case "listimage":

	case "catchimage":

	}

}

func save_baidu_base64(r *http.Request) map[string]interface{} {
	fhs := r.FormValue("upfile")
	ddd, _ := base64.StdEncoding.DecodeString(fhs)
	newname := strconv.FormatInt(time.Now().Unix(), 10) // + "_" + filename
	filename := common.Get_Upload_filename(newname, "")
	ioutil.WriteFile(filename+".jpg", ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	json_str := map[string]interface{}{
		"state":    "SUCCESS",
		"url":      filename[1:],
		"title":    newname + ".jpg",
		"original": newname + ".jpg",
	}
	return json_str

}

func save_baidu_file(r *http.Request) map[string]interface{} {
	baidu_json := make(map[string]interface{})
	fhs := r.MultipartForm.File["upfile"][0]
	file, err := fhs.Open()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	filename := common.Get_Upload_filename(fhs.Filename, "")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()
	io.Copy(f, file)
	baidu_json["original"] = "图片"
	baidu_json["size"] = 11220
	baidu_json["state"] = "SUCCESS"
	baidu_json["title"] = "图片"
	baidu_json["type"] = ".jpg"
	baidu_json["url"] = filename[1:]
	return baidu_json
}
