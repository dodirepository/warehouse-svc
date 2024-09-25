package handlers

import (
	"fmt"
	"net/http"

	lib "github.com/dodirepository/common-lib"
	handlers "github.com/dodirepository/warehouse-svc/internal/domain/handlers"
	domain "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	"github.com/go-playground/validator/v10"
)

type CreateControllerHandlers struct {
	Usecase domain.WarehouseUseCaseInterface
}

func NewCreate(u domain.WarehouseUseCaseInterface) handlers.HttpApi {
	return &CreateControllerHandlers{
		Usecase: u,
	}
}

func (u CreateControllerHandlers) Handlers(w http.ResponseWriter, r *http.Request) {
	payload := domain.CreateWarehouse{}
	err := lib.ParseBody(r, &payload)
	if err != nil {
		lib.Render(domain.ErrorResponse{
			Message: "Failed To Decode Payload",
		}, http.StatusUnprocessableEntity, w)
		return
	}

	validate := validator.New()
	trans := lib.TranslatorValidatorIDN(validate)
	err = validate.Struct(payload)
	errs := lib.TranslateError(err, trans)
	if err != nil {
		lib.Render(domain.ErrorResponse{
			Message: fmt.Sprintf("%v", errs),
		}, http.StatusBadRequest, w)
		return
	}

	resperr := u.Usecase.CreateHouse(payload)
	if resperr != nil {
		lib.Render(resperr, resperr.Status, w)
		return
	}

	lib.Render(http.StatusText(http.StatusOK), http.StatusCreated, w)
}
