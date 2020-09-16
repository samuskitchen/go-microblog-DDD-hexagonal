package infrastructure

import (
	"github.com/go-chi/chi"
	v1post "microblog/domain/post/application/v1"
	v1user "microblog/domain/user/application/v1"
	"net/http"
)

// Routes returns post router with each endpoint.
func RoutesPost(pr *v1post.PostRouter) http.Handler {
	newRouter := chi.NewRouter()

	newRouter.Get("/user/{userId}",pr.GetByUserHandler)
	newRouter.Get("/", pr.GetAllPost)
	newRouter.Post("/", pr.CreateHandler)
	newRouter.Get("/{id}", pr.GetOneHandler)
	newRouter.Put("/{id}", pr.UpdateHandler)
	newRouter.Delete("/{id}", pr.DeleteHandler)

	return newRouter
}


// Routes returns user router with each endpoint.
func RoutesUser(ur *v1user.UserRouter) http.Handler {
	newRouter := chi.NewRouter()

	newRouter.Get("/", ur.GetAllUser)
	newRouter.Post("/", ur.CreateHandler)
	newRouter.Get("/{id}", ur.GetOneHandler)
	newRouter.Put("/{id}", ur.UpdateHandler)
	newRouter.Delete("/{id}", ur.DeleteHandler)

	return newRouter
}