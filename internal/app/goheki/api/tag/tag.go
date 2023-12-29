package tag

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	//"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
	"time"
)

type Tag struct {
	ID        *int64    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DeleteIDs struct {
	IDs []int64 `json:"ids"`
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
	var tags []Tag
	query := `
		INSERT INTO tag (
			name,
			created_at
		) VALUES (
			:name,
			:created_at
		)
	`
	err := json.NewDecoder(r.Body).Decode(&tags)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, tag := range tags {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, tag)
		if err != nil {
			log.Fatal(fmt.Sprintf("insert error: %v", err))
		}
	}
}