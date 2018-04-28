package dbcurd

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
	"gylib/common/session"
	"gylib/common/webclient"
	"gylib/dbfun"
)

type DbcurdController struct {
	Site_name string
	Api_url   string
	Get_url   string
	Masterdb  string
	T_data    map[string]string
	F_data    []map[string]string
	controller.Controller
}

type Json_msg struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func NewDbcurdController(w http.ResponseWriter, r *http.Request, masterdb string) (*DbcurdController, bool) {
	this := new(DbcurdController)
	this.W = w
	this.R = r
	this.Data = make(map[string]interface{})
	this.Url_list = make(map[string]interface{})
	this.JsonData = make(map[string]interface{})
	this.Postdata = make(map[string]interface{})
	this.T_data=make(map[string]string)
	this.F_data=make([]map[string]string,0)
	this.T_data,this.F_data=this.Get_redis_tbname(masterdb)
	this.Get_url = r.URL.RawQuery
	data := common.Getini("conf/app.ini", "server", map[string]string{"site_name": "", "api_url": ""})
	this.Site_name = data["site_name"]
	this.Api_url = data["api_url"]
	this.Masterdb=masterdb
	sessions.Start_session(this.W, this.R)
	sessions.Gsession.Set("updatetime", time.Now().Unix())
	this.Struct_init()
	flag := this.check_login_sess()
	this.Set_token_ttl()
	return this, flag
}

func (this *DbcurdController) check_login_sess() (bool) {
	flag:=false
	u_id := sessions.Gsession.Get("u_id");
	if u_id != nil {
		flag = true
	}
	return flag
}

func (this *DbcurdController) Http_Post(data map[string]interface{}) []byte {
	url_add := this.Api_url
	if (this.Get_url != "") {
		url_add = this.Api_url + "/?" + this.Get_url
	}
	token := sessions.Gsession.Get("token")
	if (token != nil) {
		data["token"] = token
	} else {
		data["sign"] = this.Set_map_Signature(data)
	}
	result := webclient.Https_post(url_add, data)
	return []byte(result)

}

func (this *DbcurdController) Struct_init() {

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

func (this *DbcurdController) error_return(msg string) {
	jsonstr := make(map[string]interface{})
	jsonstr["status"] = 101
	jsonstr["msg"] = msg
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}

//获取sign
func (this *DbcurdController) Set_map_Signature(data map[string]interface{}) string {
	redis_ct := redispack.Get_redis_pool()
	client := redis_ct.Get()
	defer client.Close()
	hasok, _ := client.Do("EXISTS", "appid")
	var appid, appkey string
	if (hasok == 0) {
		data := common.Getini("conf/app.ini", "token", map[string]string{"appid": "gysdk", "appkey": ""})
		appid = data["appid"]
		appkey = data["appkey"]
		client.Do("SET", "appid", data["appid"])
		client.Do("SET", "appkey", data["appkey"])
	}
	sin_str := common.Signature_MD5(appid, appkey, datatype.Map2str(data))
	return sin_str

}

//appid有效性验证
func (this *DbcurdController) Check_Signatur() bool {
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

func (this *DbcurdController) Set_token_ttl() {
	token := sessions.Gsession.Get("token")
	if (token == nil) {
		return
	}
	client := redispack.Get_redis_pool().Get()
	defer client.Close()
	hasok, _ := client.Do("EXISTS", token)
	if hasok == 0 {
		return
	}
	raw, _ := client.Do("Get", token)
	if raw != nil {
		client.Do("SETEX", token, 3600, raw)
	}

}

func (this *DbcurdController) Check_token() bool {
	token := this.R.FormValue("token")
	sess_token := datatype.Type2str(sessions.Gsession.Get("token"))
	if (token != sess_token) {
		return false
	}
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

func (this *DbcurdController) Multi_upload(uploadfile string) {
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

func (this *DbcurdController) Single_upload(uploadfile string) {
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

func (this *DbcurdController) Get_token() {
	uu_id := common.Get_UUID()
	Redis_Pool := redispack.Get_redis_pool()
	client := Redis_Pool.Get()
	defer client.Close()
	jsonstr := new(Json_msg)
	jsonstr.Data = uu_id
	data, _ := json.Marshal(&jsonstr.Data)
	client.Do("SETEX", uu_id, 3600, string(data))
	sessions.Gsession.Set("token", uu_id)
	jsonstr.Status = 100
	jsonstr.Msg = "Token获取成功"
	b, _ := json.Marshal(&jsonstr)
	this.W.Write(b)
}

/*{
  "code": 0 //0表示成功，其它失败
  ,"msg": "" //提示信息 //一般上传失败后返回
  ,"data": {
    "src": "图片路径"
    ,"title": "图片名称" //可选
  }
}*/

func (this *DbcurdController) AdminLay_upload_fileAction() {
	jsonstr := make(map[string]interface{})
	jsonstr["msg"] = "上传失败，请检查"
	jsonstr["code"] = 1
	jsonstr["data"] = nil
	//fmt.Println(this.R)
	temp_file := common.Upload_web_file(this.R)
	if (len(temp_file) > 0) {
		if (temp_file[len(temp_file)-1] != "") {
			data := make(map[string]interface{})
			data["src"] = this.Site_name + temp_file[len(temp_file)-1]
			jsonstr["msg"] = "上传成功"
			jsonstr["code"] = 0
			jsonstr["data"] = data
		}
	}
	this.W.Header().Set("content-type", "application/json")
	b, _ := json.Marshal(&jsonstr)
	//fmt.Println(string(b))
	//fmt.Println(jsonstr)
	this.W.Write(b)
}

func (this *DbcurdController) Get_redis_tbname(tbname string)(map[string]string,[]map[string]string){
	redis_ct := redispack.Get_redis_pool()
	client := redis_ct.Get()
	defer client.Close()
	r_tb, _ := client.Do("HGET","tb_dict" , tbname)
	if(r_tb==nil){
		return this.Get_mysql_dict(tbname)
	}
	f_db, _ := client.Do("HGET", "fd_dict",tbname)
	tb_data:=make(map[string]string)
	fd_data:=make([]map[string]string,0)
	json.Unmarshal(r_tb.([]byte), &tb_data)
	json.Unmarshal(f_db.([]byte), &fd_data)
	return tb_data, fd_data

}

func (this *DbcurdController) Get_mysql_dict(tbname string) (map[string]string,[]map[string]string) {
	db := lib.NewQuerybuilder()
	data := db.Tbname("db_tb_dict").Where(fmt.Sprintf("name='%v'", lib.Db_perfix+tbname)).Find();
	db.Dbinit()
	fd_data := db.Tbname("db_fd_dict").Where(fmt.Sprintf("t_id=%v", data["id"])).Select()
	redis_ct := redispack.Get_redis_pool()
	client := redis_ct.Get()
	defer client.Close()
	db_data,_:=json.Marshal(&data)
	f_data,_:=json.Marshal(&fd_data)
	client.Do("HSET",  "tb_dict",tbname,db_data)
	client.Do("HSET",  "fd_dict",tbname,f_data)
	return data,fd_data
}