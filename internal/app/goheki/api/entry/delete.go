package entry

import (
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type DeleteHandler struct {
	svc *service.IndexService
}

func NewDeleteHandler(svc *service.IndexService) *DeleteHandler {
	return &DeleteHandler{
		svc: svc,
	}
}

func (h *DeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var entry []Entry
	query := `
		DELETE FROM
			entry
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		return
	}
	query, args, err := db.In(query, entry)
	if err != nil {
		return
	}
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		return
	}
	json.NewEncoder(w).Encode(&entry)
}