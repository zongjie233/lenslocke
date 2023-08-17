package controllers

import (
	"fmt"
	"github.com/zongjie233/lenslocked/context"
	"github.com/zongjie233/lenslocked/models"
	"net/http"
	"net/url"
)

// 保存用户部分中使用的模板
type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassWord  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		u.Templates.New.Execute(w, r, data, err)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	fmt.Fprintf(w, "User created:%v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "processSignIn wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "processSignIn wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {

	user := context.User(r.Context())
	fmt.Fprintf(w, "current user: %s\n", user.Email)

	//token, err := readCookie(r, CookieSession)
	//if err != nil {
	//	fmt.Println(err)
	//	http.Redirect(w, r, "/signin", http.StatusFound)
	//	return
	//}

	//user, err := u.SessionService.User(token)
	//if err != nil {
	//	fmt.Println(err)
	//	http.Redirect(w, r, "/signin", http.StatusFound)
	//	return
	//}
	//fmt.Fprintf(w, "current user: %s\n", user.Email)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		fmt.Println("读取cookie出错了")
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {
		fmt.Println("删除cookie出错了")

		fmt.Println(err)
		http.Error(w, "sth wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)

}

// ForgotPassword 处理页面渲染
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassWord.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	// auth token
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	//update password
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	//

	/*
		重新设置 Cookie 的目的是确保用户在密码重置后，会话仍然保持有效。如果不重新设置 Cookie，用户在密码重置后需要手动重新登录，才能建立
		一个新的有效会话。通过重新设置 Cookie，用户可以继续保持登录状态，无需再次输入用户名和密码。
	*/
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
	//u.Templates.ResetPassWord.Execute(w, r, data)
}

// ProcessForgotPassword 处理表单提交逻辑
func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.lenslocked.com/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}

	u.Templates.CheckYourEmail.Execute(w, r, data)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

// SetUser 处理HTTP请求时从会话中获取用户信息，并将用户信息存储到请求的上下文中
func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取会话cookie,并获取用户信息
		token, err := readCookie(r, CookieSession)
		// 如果失败则继续处理下一个处理器
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		// 将用户信息存储到请求的上下文中。
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// RequireUser 用于在处理HTTP请求时检查用户是否已登录。如果用户未登录，它将重定向到登录页面，否则它会继续执行下一个处理器函数。
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
