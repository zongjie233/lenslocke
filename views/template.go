package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/zongjie233/lenslocked/context"
	"github.com/zongjie233/lenslocked/models"
	"html/template"
	"io"
	"io/fs"

	"log"
	"net/http"
)

// Must 简化使用模板时的错误处理。如果在解析或执行模板时发生错误，程序将死机并停止执行。
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFS 返回一个 Template 包含已解析模板的结构
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// 必须在解析模版之前声明
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("功能还未完成")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("功能还未完成")
			},
		})
	// 解析文件系统中的模板文件，并将其解析到tpl模板中。
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsfs template: %v", err)
	}

	return Template{
		htmlTpl: tpl,
	}, nil
}

// 已经解析的Gohtml模板
type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "there was an error rendering the page", http.StatusInternalServerError)
		return
	}
	// 为模板添加一个自定义函数。
	// 注释：这个函数为模板提供一个名为 "csrfField" 的自定义函数，返回从请求中获取的CSRF令牌。
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)
	w.Header().Set("Content-type", "text/html,charset=utf-8")
	var buf bytes.Buffer
	// 执行模板，并将结果写入缓存器
	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("解析模板出错: %v", err)
		http.Error(w, "there was an error parsing the template", http.StatusInternalServerError)
		return
	}
	// 复制到w中
	io.Copy(w, &buf)
}

//func Parse(filepath string) (Template, error) {
//	tpl, err := template.ParseFiles(filepath)
//	if err != nil {
//		return Template{}, fmt.Errorf("parsing template: %v", err)
//	}
//	return Template{
//		htmlTpl: tpl,
//	}, nil
//}
