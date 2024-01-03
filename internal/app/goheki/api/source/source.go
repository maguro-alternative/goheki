package source

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	//"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

type Source struct {
	ID      *int64 `db:"id" json:"id"`
	EntryID int64  `db:"entry_id" json:"entry_id"`
	Name    string `db:"name" json:"name"`
	Url     string `db:"url" json:"url"`
	Type    string `db:"type" json:"type"`
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
	var sources []Source
	query := `
		INSERT INTO source (
			entry_id,
			name,
			url,
			type
		) VALUES (
			:entry_id,
			:name,
			:url,
			:type
		)
	`
	err := json.NewDecoder(r.Body).Decode(&sources)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, source := range sources {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, source)
		if err != nil {
			log.Fatal(fmt.Sprintf("insert error: %v", err))
		}
	}
	json.NewEncoder(w).Encode(sources)
}

