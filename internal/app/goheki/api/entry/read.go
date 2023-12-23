package entry

import (
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type ReadHandler struct {
	svc *service.IndexService
}

func NewReadHandler(svc *service.IndexService) *ReadHandler {
	return &ReadHandler{
		svc: svc,
	}
}

func (h *ReadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var entry []Entry
	query := `
		SELECT
			name,
			image,
			content,
			create_at
		FROM
			entry
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
