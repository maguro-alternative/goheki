package personalitytype

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
	var personalityTypes []PersonalityType
	query := `
		INSERT INTO personality_type (
			type
		) VALUES (
			:type
		)
	`
	err := json.NewDecoder(r.Body).Decode(&personalityTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		return
	}
	for _, personalityType := range personalityTypes {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, personalityType)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&personalityTypes)
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
	var personalityTypes []PersonalityType
	query := `
		SELECT
			id,
			type
		FROM
			personality_type
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				id,
				type
			FROM
				personality_type
		`
		err := h.svc.DB.SelectContext(r.Context(), &personalityTypes, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		err = json.NewEncoder(w).Encode(&personalityTypes)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
			return
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				id,
				type
			FROM
				personality_type
			WHERE
				id = ?
		`
		err := h.svc.DB.SelectContext(r.Context(), &personalityTypes, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		err = json.NewEncoder(w).Encode(&personalityTypes)
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
	err = h.svc.DB.SelectContext(r.Context(), &personalityTypes, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
		return
	}
	err = json.NewEncoder(w).Encode(&personalityTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
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
	if r.Method != http.MethodPut {
		return
	}
	var personalityTypes []PersonalityType
	query := `
		UPDATE
			personality_type
		SET
			type = :type
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&personalityTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		return
	}
	for _, personalityType := range personalityTypes {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, personalityType)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&personalityTypes)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
}
