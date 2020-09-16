package persistence

import (
	"context"
	"microblog/domain/user/domain"
	"time"

	conn "microblog/infrastructure/database"
)

// UserRepository manages the operations with the database that correspond to the user model.
type UserRepository struct {
	Data *conn.Data
}

// GetAll returns all users.
func (ur *UserRepository) GetAllUser(ctx context.Context) ([]domain.User, error) {
	rows, err := ur.Data.DB.QueryContext(ctx, selectAllUser)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var userRow domain.User
		_ = rows.Scan(&userRow.ID, &userRow.FirstName, &userRow.LastName, &userRow.Username, &userRow.Email, &userRow.Picture, &userRow.CreatedAt, &userRow.UpdatedAt)
		users = append(users, userRow)
	}

	return users, nil
}

// GetOne returns one user by id.
func (ur *UserRepository) GetOne(ctx context.Context, id uint) (domain.User, error) {
	row := ur.Data.DB.QueryRowContext(ctx, selectUserById, id)

	var userScan domain.User
	err := row.Scan(&userScan.ID, &userScan.FirstName, &userScan.LastName, &userScan.Username, &userScan.Email, &userScan.Picture, &userScan.CreatedAt, &userScan.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return userScan, nil
}

// GetByUsername returns one user by username.
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	row := ur.Data.DB.QueryRowContext(ctx, selectUSerByUsername, username)

	var userScan domain.User
	err := row.Scan(&userScan.ID, &userScan.FirstName, &userScan.LastName, &userScan.Username,
		&userScan.Email, &userScan.Picture, &userScan.PasswordHash, &userScan.CreatedAt, &userScan.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return userScan, nil
}

// Create adds a new user.
func (ur *UserRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	if user.Picture == "" {
		user.Picture = "https://placekitten.com/g/300/300"
	}

	stmt, err := ur.Data.DB.PrepareContext(ctx, insertUser)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, user.FirstName, user.LastName, user.Username, user.Email,
		user.Picture, user.PasswordHash, now, now,
	)

	err = row.Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a user by id.
func (ur *UserRepository) Update(ctx context.Context, id uint, u domain.User) error {
	now := time.Now().Truncate(time.Second).Truncate(time.Millisecond).Truncate(time.Microsecond)

	stmt, err := ur.Data.DB.PrepareContext(ctx, updateUser)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.FirstName, u.LastName, u.Email, u.Picture, now, id)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a user by id.
func (ur *UserRepository) Delete(ctx context.Context, id uint) error {

	stmt, err := ur.Data.DB.PrepareContext(ctx, deleteUser)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
