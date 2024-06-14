package data

import (
	"context"
	"database/sql"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

func New(db *sql.DB) Models {
	return Models{
		User: UserModel{db},
	}
}

type Models struct {
	User UserModel
}

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where email = $1`

	var user User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at from users where id = $1`

	var user User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// func (m *UserModel) Update() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	stmt := `update users set
// 		email = $1,
// 		first_name = $2,
// 		last_name = $3,
// 		user_active = $4,
// 		updated_at = $5
// 		where id = $6
// 	`

// 	_, err := db.ExecContext(ctx, stmt,
// 		u.Email,
// 		u.FirstName,
// 		u.LastName,
// 		u.Active,
// 		time.Now(),
// 		u.ID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (m *UserModel) Delete() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	stmt := `delete from users where id = $1`

// 	_, err := m.DB.ExecContext(ctx, stmt, u.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (m *UserModel) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `insert into users (email, first_name, last_name, password, user_active, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err = m.DB.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// func (m *UserModel) ResetPassword(password string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
// 	if err != nil {
// 		return err
// 	}

// 	stmt := `update users set password = $1 where id = $2`
// 	_, err = m.DB.ExecContext(ctx, stmt, hashedPassword, u.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (m *UserModel) PasswordMatches(plainText string) (bool, error) {
// 	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
// 			// invalid password
// 			return false, nil
// 		default:
// 			return false, err
// 		}
// 	}

// 	return true, nil
// }
