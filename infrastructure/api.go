package infrastructure

import (
	v1post "microblog/domain/post/application/v1"
	persistencePost "microblog/domain/post/infraestructure/persistence"
	v1user "microblog/domain/user/application/v1"
	persistenceUser "microblog/domain/user/infraestructure/persistence"
	"net/http"

	"github.com/go-chi/chi"
	data "microblog/infrastructure/database"
)

// New returns the API V1 Handler with configuration.
func New(conn *data.Data) http.Handler {
	r := chi.NewRouter()

	ur := &v1user.UserRouter{
		Repository: &persistenceUser.UserRepository{
			Data: conn,
		},
	}
	r.Mount("/users", RoutesUser(ur))

	pr := &v1post.PostRouter{
		Repository: &persistencePost.PostRepository{
			Data: conn,
		},
	}
	r.Mount("/posts", RoutesPost(pr))

	return r
}
