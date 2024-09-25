package handlers

import (
	"net/http"
	"strconv"

	lib "github.com/dodirepository/common-lib"
	handlers "github.com/dodirepository/warehouse-svc/internal/domain/handlers"
	domain "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
)

type UpdateStatusControllerHandler struct {
	Usecase domain.WarehouseUseCaseInterface
}

func NewUpdateStatus(u domain.WarehouseUseCaseInterface) handlers.HttpApi {
	return &UpdateStatusControllerHandler{
		Usecase: u,
	}
}

func (u UpdateStatusControllerHandler) Handlers(w http.ResponseWriter, r *http.Request) {

	ID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		lib.Render(domain.ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest, w)
		return
	}

	isActive, err := strconv.ParseBool(r.URL.Query().Get("is_active"))
	if err != nil {
		lib.Render(domain.ErrorResponse{
			Message: err.Error(),
		}, http.StatusBadRequest, w)
		return
	}
	resperr := u.Usecase.UpdateStatusWarehouseByID(ID, isActive)
	if resperr != nil {
		lib.Render(resperr, resperr.Status, w)
		return
	}

	lib.Render(http.StatusText(http.StatusOK), http.StatusOK, w)
}
