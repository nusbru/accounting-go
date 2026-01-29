package transaction

import (
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type GetTransactionHandler struct {
	service interfaces.TransactionService
}

func NewGetTransactionHandler(service interfaces.TransactionService) *GetTransactionHandler {
	return &GetTransactionHandler{service: service}
}

// @Summary Get a transaction
// @Description Retrieve a transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} TransactionResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Transaction not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/transactions/{id} [get]
func (h *GetTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		problem := common.NewMethodNotAllowedProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	id := extractID(r.URL.Path, "/api/v1/transactions/")
	if id == "" {
		validationErrors := []common.ValidationError{
			{Field: "id", Message: "transaction ID is required"},
		}
		problem := common.NewValidationProblemWithErrors(r.RequestURI, validationErrors)
		common.WriteProblem(w, problem)
		return
	}

	// Validate ID is a valid UUID
	if validationErr := common.ValidateUUID(id, "id"); validationErr != nil {
		problem := common.NewValidationProblemWithErrors(r.RequestURI, []common.ValidationError{*validationErr})
		common.WriteProblem(w, problem)
		return
	}

	transaction, err := h.service.GetTransaction(r.Context(), id)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			problem := common.NewNotFoundProblem(err.Error(), r.RequestURI)
			common.WriteProblem(w, problem)
			return
		}
		problem := common.NewInternalErrorProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}
	if transaction == nil {
		problem := common.NewNotFoundProblem("transaction not found", r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	common.WriteJSON(w, http.StatusOK, toTransactionResponse(transaction))
}
