package entry

import (
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"

	"encoding/json"
	"net/http"
	"time"
)

type Entry struct {
	Name     string    `db:"name"`
	Image    string    `db:"image"`
	Content  string    `db:"content"`
	CreateAt time.Time `db:"create_at"`
}

type CreateHandler struct {
	svc *service.IndexService
}

func NewCreateHandler(svc *service.IndexService) *CreateHandler {
	return &CreateHandler{
		svc: svc,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var entry Entry
	query := `
		INSERT INTO entry (
			name,
			image,
			content,
			create_at
		) VALUES (
			:Name,
			:Image,
			:Content,
			:CreateAt
		)
	`
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		return
	}
	_, err = h.svc.DB.ExecContext(r.Context(), query, entry)
	if err != nil {
		return
	}
}
