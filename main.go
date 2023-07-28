package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/zongjie233/lenslocked/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zongjie233/lenslocked/controllers"
	"github.com/zongjie233/lenslocked/templates"
	"github.com/zongjie233/lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	fmt.Println("connected!")
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

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "page not found", http.StatusNotFound)
	})
	fmt.Println("starting the server on :3000....")
	csrfKey := "abcdefghigklmnopqrstuvwxyzsfhsdf"
	csrfMw := csrf.Protect([]byte(csrfKey),
		//TODO: 生产环境下改成true
		csrf.Secure(false))
	http.ListenAndServe(":3000", csrfMw(r))
}
