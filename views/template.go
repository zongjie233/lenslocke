package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %v", err)
	}
	return Template{
		htmlTpl: tpl,
	}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Excute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-type", "text/html,charset=utf-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("解析模板出错: %v", err)
		http.Error(w, "there was an error parsing the template", http.StatusInternalServerError)
		return // 不在运行下面的代码
	}

}
