package store

import "github.com/melZula/securebin/internal/app/model"

// DataRepository ...
type DataRepository interface {
	Create(*model.Data) error
	Find(int) (*model.Data, error)
	// FindByEmail(string) (*model.Data, error)
}
