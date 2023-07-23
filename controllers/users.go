package controllers

import (
	"fmt"
	"net/http"
)

// 保存用户部分中使用的模板
type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 下面两种方法都可以
	// fmt.Fprint(w, "Email:", r.PostForm.Get("email"))
	// fmt.Fprint(w, "Password:", r.PostForm.Get("password"))

	fmt.Fprint(w, "Email:", r.FormValue("email"))
	fmt.Fprint(w, "Email:", r.FormValue("password"))

}
