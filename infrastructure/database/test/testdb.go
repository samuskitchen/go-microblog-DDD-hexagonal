package testdb

import (
	domain2 "microblog/domain/post/domain"
	"microblog/domain/user/domain"
	"time"

	"database/sql"
	"github.com/pkg/errors"
	db "microblog/infrastructure/database"
)

// Open returns a new database connection for the test database.
func Open() *db.Data {
	return db.NewTest()
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sql.DB) error {
	stmt := "TRUNCATE TABLE users RESTART IDENTITY CASCADE;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

// SeedUsers handles seeding the user table in the database for integration tests.
func SeedUsers(dbc *sql.DB) ([]domain.User, error) {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	users := []domain.User{
		{
			FirstName: "Daniel",
			LastName: "De La Pava Suarez",
			Username: "daniel.delapava",
			Email: "daniel.delapava@jikkosoft.com",
			Password: "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			FirstName: "Rebecca",
			LastName: "Romero",
			Username: "rebecca.romero",
			Email: "rebecca.romero@jikkosoft.com",
			Password: "123456",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for i := range users {
		query := `INSERT INTO users (
				first_name, 
				last_name, 
				username, 
				email, 
				picture, 
				password, 
				created_at, 
				updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

		stmt, err := dbc.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "prepare user insertion")
		}

		row := stmt.QueryRow(&users[i].FirstName, &users[i].LastName, &users[i].Username, &users[i].Email, &users[i].Picture, &users[i].Password, &users[i].CreatedAt, &users[i].UpdatedAt)

		if err = row.Scan(&users[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture user id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return users, nil
}

// SeedPosts handles seeding the post table in the database for integration tests.
func SeedPosts(dbc *sql.DB) ([]domain2.Post, error) {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	users, err := SeedUsers(dbc)
	if err != nil {
		return nil, errors.Wrap(err, "error data users")
	}

	posts := []domain2.Post{
		{
			Body:      "Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
			UserID:    users[0].ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Body:      "Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
			UserID:    users[1].ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for i := range posts {
		query := `INSERT INTO posts 
					(body, user_id, created_at, updated_at) 
					VALUES ($1, $2, $3, $4) 
				  RETURNING id;`

		stmt, err := dbc.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "prepare post insertion")
		}

		row := stmt.QueryRow(posts[i].Body, posts[i].UserID, posts[i].CreatedAt, posts[i].UpdatedAt)

		if err = row.Scan(&posts[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture post id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return posts, nil
}
