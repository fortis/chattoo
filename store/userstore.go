package store

import (
	"database/sql"

	_ "github.com/lib/pq"
	"chattoo/user"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(connStr string) (*UserStore, error) {
	db, err := sql.Open("postgres", connStr)
	return &UserStore{db}, err
}

func (store *UserStore) Insert(user user.User) error {
	_, err := store.db.Query("INSERT INTO users (username, password) VALUES ($1,$2)", user.Username, user.Password)
	return err
}

func (store *UserStore) IsExists(username string) bool {
	var count int
	row := store.db.QueryRow("SELECT COUNT(*) from users WHERE username = $1", username)
	row.Scan(&count)
	return count > 0
}

func (store *UserStore) FindOneByName(username string, dest *user.User) error {
	err := store.db.QueryRow("SELECT id, username, password from users where username = $1", username).Scan(&dest.Id, &dest.Username, &dest.Password)
	if err != nil {
		return err
	}

	return nil
}

func (store *UserStore) FindAll() ([]user.User, error) {
	var users []user.User
	rows, err := store.db.Query("SELECT id, username from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (store *UserStore) UpdateUsername(id int64, username string) {
	store.db.Query("UPDATE users SET username = $1 WHERE id = $2", username, id)
}
