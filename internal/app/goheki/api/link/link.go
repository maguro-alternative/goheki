package link

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
	var links []Link
	query := `
		INSERT INTO link (
			entry_id,
			type,
			url,
			nsfw,
			darkness
		) VALUES (
			:entry_id,
			:type,
			:url,
			:nsfw,
			:darkness
		)
	`
	err := json.NewDecoder(r.Body).Decode(&links)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
		return
	}
	for _, link := range links {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, link)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
			return
		}
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
	var links []Link
	query := `
		SELECT
			entry_id,
			type,
			url,
			nsfw,
			darkness
		FROM
			link
		WHERE
			id IN (?)
	`
	queryIDs, ok := r.URL.Query()["id"]
	if !ok {
		query = `
			SELECT
				entry_id,
				type,
				url,
				nsfw,
				darkness
			FROM
				link
		`
		err := h.svc.DB.SelectContext(r.Context(), &links, query)
		if err != nil {
			log.Fatal(fmt.Sprintf("db error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&links)
		return
	} else if len(queryIDs) == 1 {
		query = `
			SELECT
				entry_id,
				type,
				url,
				nsfw,
				darkness
			FROM
				link
			WHERE
				id = $1
		`
		err := h.svc.DB.SelectContext(r.Context(), &links, query, queryIDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("select error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&links)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
	}
	query, args, err := db.In(query, queryIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in query error: %v", err))
	}
	query = db.Rebind(len(queryIDs), query)
	err = h.svc.DB.SelectContext(r.Context(), &links, query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("select error: %v", err))
	}
	err = json.NewEncoder(w).Encode(&links)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
	var links []Link
	query := `
		UPDATE
			link
		SET
			entry_id = :entry_id,
			type = :type,
			url = :url,
			nsfw = :nsfw,
			darkness = :darkness
		WHERE
			id = :id
	`
	err := json.NewDecoder(r.Body).Decode(&links)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	for _, link := range links {
		_, err = h.svc.DB.NamedExecContext(r.Context(), query, link)
		if err != nil {
			log.Fatal(fmt.Sprintf("update error: %v", err))
		}
	}
	err = json.NewEncoder(w).Encode(&links)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
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
			link
		WHERE
			id IN (?)
	`
	err := json.NewDecoder(r.Body).Decode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json decode error: %v body:%v", err, r.Body))
	}
	if len(delIDs.IDs) == 0 {
		return
	} else if len(delIDs.IDs) == 1 {
		query = `
			DELETE FROM
				link
			WHERE
				id = $1
		`
		_, err = h.svc.DB.ExecContext(r.Context(), query, delIDs.IDs[0])
		if err != nil {
			log.Fatal(fmt.Sprintf("delete error: %v", err))
		}
		err = json.NewEncoder(w).Encode(&delIDs)
		if err != nil {
			log.Fatal(fmt.Sprintf("json encode error: %v", err))
		}
		return
	}
	query, args, err := db.In(query, delIDs.IDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("in query error: %v", err))
	}
	query = db.Rebind(len(delIDs.IDs), query)
	_, err = h.svc.DB.ExecContext(r.Context(), query, args...)
	if err != nil {
		log.Fatal(fmt.Sprintf("delete error: %v", err))
	}
	err = json.NewEncoder(w).Encode(&delIDs)
	if err != nil {
		log.Fatal(fmt.Sprintf("json encode error: %v", err))
	}
}
