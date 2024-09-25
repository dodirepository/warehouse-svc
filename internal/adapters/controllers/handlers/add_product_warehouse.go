package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	lib "github.com/dodirepository/common-lib"
	handlers "github.com/dodirepository/warehouse-svc/internal/domain/handlers"
	domain "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AddProductControllerHandlers struct {
	Usecase domain.WarehouseUseCaseInterface
}

func NewAddProduct(u domain.WarehouseUseCaseInterface) handlers.HttpApi {
	return &AddProductControllerHandlers{
		Usecase: u,
	}
}

func (u AddProductControllerHandlers) Handlers(w http.ResponseWriter, r *http.Request) {

	payload := domain.AddingProductToWarehouse{}
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
	id := mux.Vars(r)["id"]
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		lib.Render(domain.ErrorResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError, w)
		return
	}

	payload.ID = ID
	resperr := u.Usecase.AddProduct(payload)
	if resperr != nil {
		lib.Render(resperr, resperr.Status, w)
		return
	}

	lib.Render(http.StatusText(http.StatusOK), http.StatusCreated, w)
}
