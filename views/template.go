package views

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/zongjie233/lenslocked/context"
	"github.com/zongjie233/lenslocked/models"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
)

type public interface {
	Public() string
}

// Must 简化使用模板时的错误处理。如果在解析或执行模板时发生错误，程序将死机并停止执行。
func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFS takes a filesystem and a list of patterns and returns a Template and an error
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {

	// Create a new template with the first pattern as the name
	// path.Base()返回路径中最后一个部分,即最后一个"/"的后边内容
	tpl := template.New(path.Base(patterns[0]))
	// Set the template functions
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("功能还未完成")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("功能还未完成")
			},
			"errors": func() []string {
				return []string{
					"Don't do that!",
					"The email is used",
					"something wrong",
				}
			},
		})

	// Parse the template with the given filesystem and patterns
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsfs template: %v", err)
	}

	// Return the Template and no error
	return Template{
		htmlTpl: tpl,
	}, nil
}

// 已经解析的Gohtml模板
type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "there was an error rendering the page", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
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
			"errors": func() []string {
				return errMsgs
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

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			println(err)
			msgs = append(msgs, "something wrong")
		}
	}
	return msgs
}
