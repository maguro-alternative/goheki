package haircolortype

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
	var hairColorTypes []HairColorType
	query := `
		INSERT INTO haircolor_type (
			color
		) VALUES (
			:color
		)
	`
	err := json.NewDecoder(r.Body).Decode(&hairColorTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		return
	}
	for _, hairColorType := range hairColorTypes {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairColorType)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&hairColorTypes)
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
	var hairColorTypes []HairColorType
	query := `
		SELECT
			id,
			color
		FROM
			haircolor_type
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				id,
				color
			FROM
				haircolor_type
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairColorTypes, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		err = json.NewEncoder(w).Encode(&hairColorTypes)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
			return
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				id,
				color
			FROM
				haircolor_type
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &hairColorTypes, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		err = json.NewEncoder(w).Encode(&hairColorTypes)
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
	err = h.svc.DB.SelectContext(r.Context(), &hairColorTypes, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
		return
	}
	err = json.NewEncoder(w).Encode(&hairColorTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
}
