package entry

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
	"time"
)

type Entry struct {
	Name     string    `db:"name" json:"name"`
	Image    string    `db:"image" json:"image"`
	Content  string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
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
	var entrys []Entry
	query := `
		INSERT INTO entry (
			name,
			image,
			content,
			created_at
		) VALUES (
			:name,
			:image,
			:content,
			:created_at
		)
	`
	err := json.NewDecoder(r.Body).Decode(&entrys)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, entry := range entrys {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, entry)
		if err != nil {
			log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
		}
	}
	json.NewEncoder(w).Encode(&entrys)
}

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
	var entrys []Entry
	query := `
		SELECT
			name,
			image,
			content,
			created_at
		FROM
			entry
	`
	//err := json.NewDecoder(r.Body).Decode(&entrys)
	//if err != nil {
	//log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	//}
	//query, args, err := db.In(query, entrys)
	//if err != nil {
	//return
	//}
	err := h.svc.DB.SelectContext(r.Context(), &entrys, query)
	if err != nil {
		log.Fatal(fmt.Sprintf("db.ExecContext error: %v \nqurey:%v", err, query))
	}
	json.NewEncoder(w).Encode(&entrys)
}

type UpdateHandler struct {
	svc *service.IndexService
}

func NewUpdateHandler(svc *service.IndexService) *UpdateHandler {
	return &UpdateHandler{
		svc: svc,
	}
}

func (h *UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	var entry []Entry
	query := `
		UPDATE
			entry
		SET
			name = :name,
			image = :image,
			content = :content,
			created_at = :create_at
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
