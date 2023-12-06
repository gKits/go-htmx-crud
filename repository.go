package crud

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, user User) (*User, error)
	All(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, id uint, updated User) (*User, error)
	Delete(ctx context.Context, id uint) error
	AvgAge(ctx context.Context) (float32, error)
	AvgHeight(ctx context.Context) (float32, error)
	Search(context.Context, string) ([]User, error)
}

func NewRepository(driver string, conn string) (Repository, error) {
	var repo Repository
	var err error

	switch driver {
	case "postgres":
		repo, err = NewPostgresRepository(conn)
		break
	case "sqlite3":
		repo, err = NewSQLiteRepository(conn)
		break
	default:
		return nil, fmt.Errorf("driver '%s' not supported", driver)
	}
	if err != nil {
		return nil, err
	}

	if err := repo.Migrate(context.Background()); err != nil {
		return nil, err
	}

	return repo, err
}

/**
SQLite Repository
**/

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(conn string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Migrate(ctx context.Context) error {
	q := `
	CREATE TABLE IF NOT EXISTS user(
		id integer NOT NULL PRIMARY KEY,
		name string NOT NULL,
		age integer NOT NULL,
		height float NOT NULL
	);
	`

	_, err := r.db.ExecContext(ctx, q)
	return err
}

func (r *SQLiteRepository) Create(ctx context.Context, user User) (*User, error) {
	q := "INSERT INTO user(name, age, height) values($1, $2, $3) RETURNING id;"

	var id uint
	if err := r.db.QueryRowContext(ctx, q, user.Name, user.Age, user.Height).Scan(&id); err != nil {
		return nil, err
	}

	user.ID = id

	return &user, nil
}

func (r *SQLiteRepository) All(ctx context.Context) ([]User, error) {
	q := "SELECT * FROM user;"

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Height); err != nil {
			return nil, err
		}
		all = append(all, user)
	}

	return all, nil
}

func (r *SQLiteRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	q := "SELECT * FROM user WHERE ID = $1;"

	row := r.db.QueryRowContext(ctx, q, id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Height); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *SQLiteRepository) Update(ctx context.Context, id uint, user User) (*User, error) {
	q := "UPDATE user SET name = $1, age = $2, height = $3 WHERE id = $4 RETURNING id;"

	row := r.db.QueryRowContext(ctx, q, user.Name, user.Age, user.Height, id)

	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, id uint) error {
	q := "DELETE FROM user WHERE id = ?;"

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(res)

	// affected, err := res.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *SQLiteRepository) AvgAge(ctx context.Context) (float32, error) {
	q := "SELECT AVG(age) FROM user;"

	row := r.db.QueryRowContext(ctx, q)

	var avgAge float32
	if err := row.Scan(&avgAge); err != nil {
		return -1, err
	}

	return avgAge, nil
}

func (r *SQLiteRepository) AvgHeight(ctx context.Context) (float32, error) {
	q := "SELECT AVG(height) FROM user;"

	row := r.db.QueryRowContext(ctx, q)

	var avgHeight float32
	if err := row.Scan(&avgHeight); err != nil {
		return -1, err
	}

	return avgHeight, nil
}

func (r *SQLiteRepository) Search(ctx context.Context, search string) ([]User, error) {
	q := "SELECT * FROM user WHERE name LIKE $1;"
	rows, err := r.db.QueryContext(ctx, q, "%"+search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Height); err != nil {
			return nil, err
		}
		all = append(all, user)
	}

	return all, nil
}

/**
PostgreSQL Repository
**/

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(conn string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Migrate(ctx context.Context) error {
	q := `
	CREATE TABLE IF NOT EXISTS user(
		id SERIAL NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		age INT NOT NULL,
		height DOUBLE NOT NULL
	);`

	_, err := r.db.ExecContext(ctx, q)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, user User) (*User, error) {
	q := "INSERT INTO user(name, age, height) values($1, $2, $3) RETURNING id;"

	var id uint
	if err := r.db.QueryRowContext(ctx, q).Scan(&id); err != nil {
		return nil, err
	}

	user.ID = id

	return &user, nil
}

func (r *PostgresRepository) All(ctx context.Context) ([]User, error) {
	q := "SELECT * FROM user;"

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Height); err != nil {
			return nil, err
		}
		all = append(all, user)
	}

	return all, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	q := "SELECT * FROM user WHERE ID = $1;"

	row := r.db.QueryRowContext(ctx, q, id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Height); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint, user User) (*User, error) {
	q := "UPDATE user SET name = $1, age = $2, height = $3 WHERE id = $4;"

	res, err := r.db.ExecContext(ctx, q, user.Name, user.Age, user.Height, id)
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return nil, nil
	}

	return &user, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
	q := "DELETE FROM user WHERE id = $1;"

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return nil
	}

	return nil
}

func (r *PostgresRepository) AvgAge(ctx context.Context) (float32, error) {
	q := "SELECT AVG(age) FROM user;"

	row := r.db.QueryRowContext(ctx, q)

	var avgAge float32
	if err := row.Scan(&avgAge); err != nil {
		return -1, err
	}

	return avgAge, nil
}

func (r *PostgresRepository) AvgHeight(ctx context.Context) (float32, error) {
	q := "SELECT AVG(height) FROM user;"

	row := r.db.QueryRowContext(ctx, q)

	var avgHeight float32
	if err := row.Scan(&avgHeight); err != nil {
		return -1, err
	}

	return avgHeight, nil
}

func (r *PostgresRepository) Search(ctx context.Context, search string) ([]User, error) {
	q := "SELECT * FROM user WHERE name LIKE $1;"
	rows, err := r.db.QueryContext(ctx, q, search+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Height); err != nil {
			return nil, err
		}
		all = append(all, user)
	}

	return all, nil
}
