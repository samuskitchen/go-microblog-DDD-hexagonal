package persistence

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"microblog/domain/post/domain"
	data "microblog/infrastructure/database"
	"testing"
	"time"
)

// represent the repository
var (
	dbMockPost         *sql.DB
	connMockPost       data.Data
	postRepositoryMock PostRepository
)

// NewMockUser initialize mock connection to database
func NewMockPost() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbMockPost = db
	connMockPost = data.Data{
		DB: dbMockPost,
	}

	postRepositoryMock = PostRepository{
		Data: &connMockPost,
	}

	return mock
}

// Close attaches the provider and close the connection
func CloseMockPost() {
	err := dbMockPost.Close()
	if err != nil {
		log.Println("Error close database test")
	}
}

// dataUSer is data for test
func dataPost() []domain.Post {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	return []domain.Post{
		{
			ID:        uint(1),
			Body:      "Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
			UserID:    uint(1),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uint(1),
			Body:      "Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
			UserID:    uint(1),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func TestPostRepository_Create(t *testing.T) {

}

func TestPostRepository_Delete(t *testing.T) {

}

func TestPostRepository_GetAll(t *testing.T) {

}

func TestPostRepository_GetByUser(t *testing.T) {

}

func TestPostRepository_GetOne(t *testing.T) {

}

func TestPostRepository_Update(t *testing.T) {

}
