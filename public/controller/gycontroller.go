package controller

import (
	"net/http"
	"strings"
	"fmt"
	"time"
	"strconv"
	"html/template"
)

type Controller struct {
	W        http.ResponseWriter
	R        *http.Request
	Tplname  string
	Data     map[string]interface{}
	Url_list map[string]interface{}
	JsonData map[string]interface{}
	Postdata map[string]interface{}
	Body     []byte
}

func set_funcmap() template.FuncMap {
	tempfunc := make(template.FuncMap)
	tempfunc["str2html"] = Strtohtml
	tempfunc["ischeckbox"] = ISCheckbox
	tempfunc["date2local"] = Date2Local
	tempfunc["date2int"] = Date2Int
	tempfunc["round"] = Round
	return (tempfunc)
}

func (this *Controller) Rander() {
	tempfunc := set_funcmap()
	filename := strings.Split(this.Tplname, "/")
	tmpname := filename[len(filename)-1]
	t := template.New(tmpname)
	//t := template.New(tmpname).Delims("<<",">>")
	t = t.Funcs(tempfunc)
	t, _ = t.ParseFiles(this.Tplname)
	t.Execute(this.W, this.Data)
}

func (this *Controller) MuitplRander(arg ...string) {
	tempfunc := set_funcmap()
	filename := strings.Split(this.Tplname, "/")
	tmpname := filename[len(filename)-1]
	t := template.New(tmpname)
	//t := template.New(tmpname).Delims("<<",">>")
	t = t.Funcs(tempfunc)
	t, _ = t.ParseFiles(this.Tplname)
	t, _ = t.ParseFiles(arg...)

	t.Execute(this.W, this.Data)
}

func Strtohtml(htmlstr interface{}) interface{} {
	return template.HTML(htmlstr.(string))
}

func ISCheckbox(code interface{}, qxz interface{}) interface{} {
	if qxz == nil {
		return template.HTML("")
	}
	qxzmemo := fmt.Sprintf("%v", qxz)
	code_str := fmt.Sprintf("%v", code)
	result := ""
	qxzarr := strings.Split(qxzmemo, ",")
	for _, val := range qxzarr {
		if code_str == val {
			result = "checked"
			break
		}
	}
	return template.HTML(result)
}

func Date2Local(htmlstr interface{}) interface{} {
	result := fmt.Sprintf("%v", htmlstr)
	result = strings.Replace(result, "+0800 CST", "", -1)
	result = strings.Replace(result, "+0000 UTC", "", -1)
	return template.HTML(result)
}

func Date2Int(htmlstr interface{}) interface{} {
	result := time.Unix(htmlstr.(int64), 0).Format("2006-01-02 15:04:05")
	return template.HTML(result)
}

func Round(htmlstr interface{}, s int) interface{} {
	result := fmt.Sprintf("%."+strconv.Itoa(s)+"f", htmlstr)
	return template.HTML(result)
}
