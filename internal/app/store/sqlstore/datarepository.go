package sqlstore

import (
	"database/sql"

	"github.com/melZula/securebin/internal/app/model"
	"github.com/melZula/securebin/internal/app/store"
)

// DataRepository ...
type DataRepository struct {
	store *Store
}

// Create ...
func (r *DataRepository) Create(u *model.Data) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO datas (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

// FindByEmail ...
func (r *DataRepository) FindByEmail(email string) (*model.Data, error) {
	u := &model.Data{}
	if err := r.store.db.QueryRow("SELECT id, email, encrypted_password FROM datas WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

// Find ...
func (r *DataRepository) Find(id int) (*model.Data, error) {
	u := &model.Data{}
	if err := r.store.db.QueryRow("SELECT id, email, encrypted_password FROM datas WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}
