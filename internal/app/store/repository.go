package store

import (
	"net/http"

	"github.com/melZula/securebin/internal/app/model"
)

// DataRepository ...
type DataRepository interface {
	Create(*model.Data) error
	Find(int) (*model.Data, error)
	LogRequest(int, *http.Request) error
	GetPrevRequests(int) ([]int, error)
}
