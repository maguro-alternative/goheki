package eyecolortype

import (
	"fmt"
	"log"
	"io"

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
	var eyeColorTypesJson EyeColorTypesJson
	query := `
		INSERT INTO eyecolor_type (
			color
		) VALUES (
			:color
		)
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &eyeColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = eyeColorTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, eyeColorType := range eyeColorTypesJson.EyeColorTypes {
		err = eyeColorType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, eyeColorType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&eyeColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var eyeColorTypesJson EyeColorTypesJson
	query := `
		SELECT
			id,
			color
		FROM
			eyecolor_type
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
				eyecolor_type
		`
		err := h.svc.DB.SelectContext(r.Context(), &eyeColorTypesJson.EyeColorTypes, query)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&eyeColorTypesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				id,
				color
			FROM
				eyecolor_type
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &eyeColorTypesJson.EyeColorTypes, query, queryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(&eyeColorTypesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &eyeColorTypesJson.EyeColorTypes, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&eyeColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var eyeColorTypesJson EyeColorTypesJson
	query := `
		UPDATE
			eyecolor_type
		SET
			color = :color
		WHERE
			id = :id
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &eyeColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = eyeColorTypesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validation error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, eyeColorType := range eyeColorTypesJson.EyeColorTypes {
		err = eyeColorType.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validation error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, eyeColorType)
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(w).Encode(&eyeColorTypesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			eyecolor_type
		WHERE
			id IN (?)
	`
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(fmt.Sprintf("read error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = json.Unmarshal(jsonBytes, &delIDs)
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
				eyecolor_type
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("db error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("db error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
