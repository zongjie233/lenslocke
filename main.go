package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/zongjie233/lenslocked/migrations"
	"github.com/zongjie233/lenslocked/models"
	"net/http"

	_ "github.com/go-chi/chi/v5"
	"github.com/zongjie233/lenslocked/controllers"
	"github.com/zongjie233/lenslocked/templates"
	"github.com/zongjie233/lenslocked/views"
)

func main() {
	// 设置数据库
	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg)
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// 执行迁移
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// 设置服务项
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// 设置中间件
	umw := controllers.UserMiddleware{SessionService: &sessionService}
	csrfKey := "abcdefghigklmnopqrstuvwxyzsfhsdf"
	csrfMw := csrf.Protect([]byte(csrfKey),
		//TODO: 生产环境下改成true
		csrf.Secure(false))

	//设置控制器
	usersC := controllers.Users{
		UserService:    &userService, // 传入指针
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))

	// 设置路由器和路由
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	fmt.Println("connected!")

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/users/me", func(r chi.Router) {
		// 使用中间件
		r.Use(umw.RequireUser)

		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})

	// 启动服务
	fmt.Println("starting the server on :3000....")
	http.ListenAndServe(":3000", r)
}
