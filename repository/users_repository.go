package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-open-api/models"
	"go-open-api/utils"
	"log"
)

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Delete(ctx context.Context, id string) error
}

type usersRepository struct {
	conn *sql.DB
}

func NewUsersRepository(conn *sql.DB) UsersRepository {
	return &usersRepository{conn: conn}
}

func (r *usersRepository) Create(ctx context.Context, user *models.User) error {
	query := `
insert into users
(
	id,
	email,
	password,
	created_at,
	updated_at
)
values
(
	?,
	?,
	?,
	?,
	?
)
`

	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer utils.HandleClose(stmt)

	_, err = stmt.ExecContext(
		ctx,
		user.ID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	log.Printf("User created. ID=%v\n", user.ID)

	return nil
}

func (r *usersRepository) Update(ctx context.Context, user *models.User) error {
	query := `
update users set
	email=?,
	password=?,
	updated_at=?
where
	id=?
`

	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer utils.HandleClose(stmt)

	_, err = stmt.ExecContext(
		ctx,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	log.Printf("User updated. ID=%v\n", user.ID)

	return nil
}

func (r *usersRepository) Get(ctx context.Context, id string) (*models.User, error) {
	query := `select * from users where id=?`

	row := r.conn.QueryRowContext(ctx, query, id)

	user := new(models.User)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *usersRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	query := `select * from users`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer utils.HandleClose(rows)

	var users []*models.User

	for rows.Next() {
		user := new(models.User)

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, err
}

func (r *usersRepository) Delete(ctx context.Context, id string) error {
	query := `delete from users where id=?`

	stmt, err := r.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer utils.HandleClose(stmt)

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("cannot delete: %s", id)
	}

	log.Printf("User deleted. ID=%s\n", id)

	return nil
}
