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
	var hairColorTypesJson HairColorTypesJson
	query := `
		INSERT INTO haircolor_type (
			color
		) VALUES (
			:color
		)
	`
	err := json.NewDecoder(r.Body).Decode(&hairColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = hairColorTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, hairColorType := range hairColorTypesJson.HairColorTypes {
		err = hairColorType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairColorType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&hairColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
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
	var hairColorTypesJson HairColorTypesJson
	query := `
		UPDATE
			haircolor_type
		SET
			color = :color
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&hairColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, hairColorType := range hairColorTypesJson.HairColorTypes {
		err = hairColorType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, hairColorType)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
	}
	err = json.NewEncoder(w).Encode(&hairColorTypesJson)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
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
	if r.Method != http.MethodDelete {
		return
	}
	var delIDs IDs
	query := `
		DELETE FROM
			haircolor_type
		WHERE
			id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				haircolor_type
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
		return
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("db error: %v", err))
		return
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
		return
	}
}
