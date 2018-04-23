package home

import (
	"net/http"
	"websever/public/controller"
	"time"
	"strconv"
	"strings"
	"gylib/common/datatype"
	"io/ioutil"
	"fmt"
	"math/rand"
	"gylib/common"
	"gylib/common/redispack"
	"encoding/json"
)

type HomeController struct {
	Site_name string
	controller.Controller
}

type Json_msg struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func NewHomeController(w http.ResponseWriter, r *http.Request) (*HomeController) {
	this := new(HomeController)
	this.W=w
	this.R=r
	this.Data = make(map[string]interface{})
	this.Url_list = make(map[string]interface{})
	this.JsonData = make(map[string]interface{})
	this.Postdata = make(map[string]interface{})
	data := common.Getini("conf/app.ini", "server", map[string]string{"site_name": ""})
	this.Site_name = data["site_name"]
	this.Struct_init()
	return this
}



func (this *HomeController) Struct_init() {

	if this.R.Method == "GET" {
		this.R.ParseForm()
	} else if (this.R.Method == "POST") {
		form_type := this.R.Header.Get("Content-Type")
		if strings.Contains(form_type, "multipart/form-data;") {
			this.R.ParseMultipartForm(32 << 20)
			for k, v := range this.R.PostForm {
				this.Postdata[k] = datatype.Type2str(v)
			}
		} else {
			this.R.ParseForm()
			for k, v := range this.R.PostForm {
				this.Postdata[k] = datatype.Type2str(v)
			}
		}
	} else {
		var err error
		this.Body, err = ioutil.ReadAll(this.R.Body)
		if (err != nil) {
			fmt.Println(err)
			return
		}
		if (this.R.Method == "JSON") {
			json.Unmarshal(this.Body, this.JsonData)
		}
	}

}

func (this *HomeController) error_return(msg string) {
	jsonstr := make(map[string]interface{})
	jsonstr["status"] = 101
	jsonstr["msg"] = msg
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}

//appid有效性验证
func (this *HomeController) Check_Signatur() bool {
	flag := false
	redis_ct := redispack.Get_redis_pool()
	client := redis_ct.Get()
	defer client.Close()
	hasok, _ := client.Do("EXISTS", "appid")
	var appid, appkey, sing string
	if (hasok == 0) {
		data := common.Getini("conf/app.ini", "token", map[string]string{"appid": "gysdk", "appkey": ""})
		appid = data["appid"]
		appkey = data["appkey"]
		client.Do("SET", "appid", data["appid"])
		client.Do("SET", "appkey", data["appkey"])
	} else {
		tmp, _ := client.Do("GET", "appid")
		appid = datatype.Type2str(tmp)
		tmp, _ = client.Do("GET", "appkey")
		appkey = common.Type2str(tmp)
	}
	get_data := make(map[string]string)
	for k, v := range this.R.Form {
		if (k == "sign") {
			sing = datatype.Type2str(v)
			continue
		}
		if (k == "appid") {
			continue
		}
		get_data[k] = datatype.Type2str(v)
	}
	sin_str := common.Signature_MD5(appid, appkey, get_data)
	if (sin_str == sing) {
		flag = true
	}
	return flag
}

func (this *HomeController) Check_token() bool {
	token := this.R.FormValue("token")
	if token == "" {
		return false
	} else {
		client := redispack.Get_redis_pool().Get()
		defer client.Close()
		hasok, _ := client.Do("EXISTS", token)
		if hasok == 0 {
			return false
		}
		raw, _ := client.Do("Get", token)
		if raw != nil {
			client.Do("SETEX", token, 3600, raw)
		} else {
			return false
		}

	}
	return true
}

func (this *HomeController) Multi_upload(uploadfile string) {
	jsonstr := new(Json_msg)
	jsonstr.Status = 101
	jsonstr.Msg = "上传失败"
	file_list := common.Upload_File(this.R, uploadfile)
	if (len(file_list) > 0) {
		temp_arr := make([]map[string]interface{}, 0)
		for _, val := range file_list {
			if (val != "") {
				temp := make(map[string]interface{}, 0)
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				temp["title"] = strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(r.Intn(9999))
				temp["img"] = this.Site_name + val
				temp_arr = append(temp_arr, temp)
			}
		}
		jsonstr.Status = 100
		jsonstr.Data = temp_arr
	}
	this.W.Header().Set("content-type", "application/json")
	b, _ := json.Marshal(&jsonstr)
	//fmt.Println(jsonstr)
	this.W.Write(b)
}

func (this *HomeController) Single_upload(uploadfile string) {
	jsonstr := new(Json_msg)
	jsonstr.Msg = "上传失败，请检查"
	jsonstr.Status = 101
	temp_file := common.Upload_File(this.R, uploadfile)
	if (len(temp_file) > 0) {
		if (temp_file[len(temp_file)-1] != "") {
			jsonstr.Data = this.Site_name + temp_file[len(temp_file)-1]
			jsonstr.Msg = "上传成功"
			jsonstr.Status = 100
		}
	}
	this.W.Header().Set("content-type", "application/json")
	b, _ := json.Marshal(&jsonstr)
	//fmt.Println(string(b))
	//fmt.Println(jsonstr)
	this.W.Write(b)
}

func (this *HomeController) Get_token() {
	uu_id := common.Get_UUID()
	Redis_Pool := redispack.Get_redis_pool()
	client := Redis_Pool.Get()
	defer client.Close()
	jsonstr := new(Json_msg)
	jsonstr.Data = uu_id
	data, _ := json.Marshal(&jsonstr.Data)
	client.Do("SETEX", uu_id, 3600, string(data))
	jsonstr.Status = 100
	jsonstr.Msg = "Token获取成功"
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}
