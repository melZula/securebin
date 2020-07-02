package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq" // ...
	"github.com/melZula/securebin/internal/app/store"
)

// Store ...
type Store struct {
	db             *sql.DB
	dataRepository *DataRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// Data ...
func (s *Store) Data() store.DataRepository {
	if s.dataRepository != nil {
		return s.dataRepository
	}
	s.dataRepository = &DataRepository{
		store: s,
	}

	return s.dataRepository
}
