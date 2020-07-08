package sqlstore

import (
	"database/sql"
	"net"
	"net/http"
	"time"

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
		"INSERT INTO securebin (img, encrypted_password, lifetime) VALUES (DECODE($1, 'base64'), $2, TO_TIMESTAMP($3) AT TIME ZONE 'UTC') RETURNING id",
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

// LogRequest ...
func (r *DataRepository) LogRequest(id int, req *http.Request) error {
	var i int
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return err
	}
	return r.store.db.QueryRow(
		"INSERT INTO requests (data_id, remote_addr, time) VALUES ($1, $2, TO_TIMESTAMP($3) AT TIME ZONE 'UTC') RETURNING id",
		id,
		ip,
		time.Now().UTC().Unix(),
	).Scan(&i)
}

// GetPrevRequests ...
func (r *DataRepository) GetPrevRequests(id int) ([]int, error) {
	var reqtsInt []int
	rows, err := r.store.db.Query(
		"SELECT extract(epoch from time)::BIGINT FROM requests WHERE data_id=$1;",
		id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	var v int
	for rows.Next() {
		if err := rows.Scan(&v); err == nil {
			reqtsInt = append(reqtsInt, v)
		}
	}
	return reqtsInt, nil
}
