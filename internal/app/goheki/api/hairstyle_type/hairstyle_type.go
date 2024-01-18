package hairstyletype

import (
	"fmt"
	"log"

	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/pkg/db"

	"encoding/json"
	"net/http"
)

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
	var hairStyleTypes []HairStyleType
	query := `
		INSERT INTO hairstyle_type (
			style
		) VALUES (
			:style
		)
	`
	err := json.NewDecoder(r.Body).Decode(&hairStyleTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		return
	}
	for _, hairStyleType := range hairStyleTypes {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairStyleType)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&hairStyleTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
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
	if r.Method != http.MethodGet {
		return
	}
	var hairStyleTypes []HairStyleType
	query := `
		SELECT
			id,
			style
		FROM
			hairstyle_type
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				id,
				style
			FROM
				hairstyle_type
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairStyleTypes, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		err = json.NewEncoder(w).Encode(&hairStyleTypes)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
			return
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				id,
				style
			FROM
				hairstyle_type
			WHERE
				id = $1
		`
		err := h.svc.DB.GetContext(r.Context(), &hairStyleTypes, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		err = json.NewEncoder(w).Encode(&hairStyleTypes)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
			return
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
		return
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &hairStyleTypes, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
		return
	}
	err = json.NewEncoder(w).Encode(&hairStyleTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
}
