package persons

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mrKrabsmr/commerce-edu-api/internal/apps"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type controller struct {
	service *service
	logger  *logrus.Logger
}

func newController() *controller {
	return &controller{
		service: newService(),
		logger:  core.GetLogger(),
	}
}

func (c *controller) list(writer http.ResponseWriter, request *http.Request) {
	c.logger.Info("GET /persons/")

	params := request.URL.Query()

	persons, paginate, err := c.service.getFormatList(params)
	if err != nil {
		core.SendResponse(writer, nil, http.StatusInternalServerError)
		c.logger.Error(err)
		return
	}

	if paginate != nil {
		core.SendPaginateResponse(writer, persons, paginate)
	} else {
		core.SendResponse(writer, persons, http.StatusOK)
	}
	c.logger.Info("SUCCESS")
}

func (c *controller) create(writer http.ResponseWriter, request *http.Request) {
	c.logger.Info("POST /persons/")

	var req *PersonRequestDTO

	data, err := io.ReadAll(request.Body)
	if err != nil {
		core.SendResponse(writer, nil, http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	defer request.Body.Close()

	if err = json.Unmarshal(data, &req); err != nil {
		core.SendResponse(writer, nil, http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	if err = req.Validate(); err != nil {
		core.SendResponse(writer, err.Error(), http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	if err = c.service.createObject(req); err != nil {
		core.SendResponse(writer, nil, http.StatusInternalServerError)
		c.logger.Error(err)
		return
	}

	core.SendResponse(writer, "success created", http.StatusCreated)
	c.logger.Info("SUCCESS")
}

func (c *controller) update(writer http.ResponseWriter, request *http.Request) {
	c.logger.Info("PUT or PATCH /persons/{id}/")

	var req *Person

	idStr := chi.URLParam(request, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		core.SendResponse(writer, "incorrect uuid", http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		core.SendResponse(writer, nil, http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	defer request.Body.Close()

	if err = json.Unmarshal(data, &req); err != nil {
		core.SendResponse(writer, nil, http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	done, err := c.service.updateObject(id, req)
	if err != nil {
		core.SendResponse(writer, nil, http.StatusInternalServerError)
		c.logger.Error(err)
		return
	}

	if !done {
		core.SendResponse(writer, "no object found by this id", http.StatusNotFound)
		c.logger.Info("not found")
		return
	}

	core.SendResponse(writer, "success updated", http.StatusOK)
	c.logger.Info("SUCCESS")
}

func (c *controller) delete(writer http.ResponseWriter, request *http.Request) {
	c.logger.Info("DELETE /persons/{id}/")

	idStr := chi.URLParam(request, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		core.SendResponse(writer, nil, http.StatusBadRequest)
		c.logger.Error(err)
		return
	}

	done, err := c.service.deleteObject(id)
	if err != nil {
		core.SendResponse(writer, nil, http.StatusInternalServerError)
		c.logger.Error(err)
		return
	}

	if !done {
		core.SendResponse(writer, "no object found by this id", http.StatusNotFound)
		c.logger.Error("not found")
		return
	}

	core.SendResponse(writer, "success deleted", http.StatusOK)
	c.logger.Info("SUCCESS")
}
