package article

import (
	"github.com/maguro-alternative/goheki/internal/app/goheki/service"
	"github.com/maguro-alternative/goheki/internal/app/goheki/model"

	"log"
	"encoding/json"
	"net/http"
)

type Article struct {
	Id    int
	Title string
	Body  string
}

type IndexHandler struct {
	svc *service.IndexService
}

type ShowHandler struct {
	svc *service.IndexService
}

type CreateHandler struct {
	svc *service.IndexService
}

type EditHandler struct {
	svc *service.IndexService
}

type DeleteHandler struct {
	svc *service.IndexService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewIndexHandler(svc *service.IndexService) *IndexHandler {
	return &IndexHandler{
		svc: svc,
	}
}

func NewShowHandler(svc *service.IndexService) *ShowHandler {
	return &ShowHandler{
		svc: svc,
	}
}

func NewCreateHandler(svc *service.IndexService) *CreateHandler {
	return &CreateHandler{
		svc: svc,
	}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(&model.IndexResponse{
		Message: "OK",
	})
	if err != nil {
		log.Println(err)
	}
}