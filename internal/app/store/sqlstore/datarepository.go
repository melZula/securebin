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
		"INSERT INTO securebin (img, encrypted_password, lifetime) VALUES (DECODE($1, 'base64'), $2, TO_TIMESTAMP($3)) RETURNING id",
		u.Img,
		u.EncryptedPassword,
		u.Lifetime,
	).Scan(&u.ID)
}

// Find ...
func (r *DataRepository) Find(id int) (*model.Data, error) {
	u := &model.Data{}
	if err := r.store.db.QueryRow("SELECT id, TRANSLATE(ENCODE(img, 'base64'), E'\n', ''), encrypted_password, extract(epoch from lifetime)::BIGINT FROM securebin WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Img,
		&u.EncryptedPassword,
		&u.Lifetime,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}
