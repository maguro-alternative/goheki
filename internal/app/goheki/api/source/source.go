package source

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
	var sourcesJson SourcesJson
	query := `
		INSERT INTO source (
			name,
			url,
			type
		) VALUES (
			:name,
			:url,
			:type
		)
	`
	err := json.NewDecoder(r.Body).Decode(&sourcesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = sourcesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validate error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, source := range sourcesJson.Sources {
		err = source.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validate error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, source)
		if err != nil {
			log.Printf(fmt.Sprintf("insert error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&sourcesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var sourcesJson SourcesJson
	query := `
		SELECT
			id,
			name,
			url,
			type
		FROM
			source
		WHERE
			id IN (?)
	`
	entryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				id,
				name,
				url,
				type
			FROM
				source
		`
		err := h.svc.DB.SelectContext(r.Context(), &sourcesJson.Sources, query)
		if err != nil {
			log.Printf(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if len(sourcesJson.Sources) == 0 {
			log.Printf("sources is empty")
			http.Error(w, "sources is empty", http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&sourcesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if len(entryIDs) == 1 {
		query = `
			SELECT
				id,
				name,
				url,
				type
			FROM
				source
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &sourcesJson.Sources, query, entryIDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("select error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if len(sourcesJson.Sources) == 0 {
			log.Printf("sources is empty")
			http.Error(w, "sources is empty", http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&sourcesJson)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, entryIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("in error: %v", err), entryIDs)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(entryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &sourcesJson.Sources, query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("select error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&sourcesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	var sourcesJson SourcesJson
	query := `
		UPDATE
			source
		SET
			name = :name,
			url = :url,
			type = :type
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&sourcesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = sourcesJson.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validate error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, source := range sourcesJson.Sources {
		err = source.Validate()
		if err != nil {
			log.Printf(fmt.Sprintf("validate error: %v", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, source)
		if err != nil {
			log.Printf(fmt.Sprintf("update error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	err = json.NewEncoder(w).Encode(&sourcesJson)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			source
		WHERE
			id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = delIDs.Validate()
	if err != nil {
		log.Printf(fmt.Sprintf("validate error: %v", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				source
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Printf(fmt.Sprintf("delete error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Printf(fmt.Sprintf("json encode error: %v", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Printf(fmt.Sprintf("in error: %v", err), delIDs.IDs)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Printf(fmt.Sprintf("delete error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Printf(fmt.Sprintf("json encode error: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
