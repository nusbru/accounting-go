package transaction

import (
	"net/http"

	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type DeleteTransactionHandler struct {
	service interfaces.TransactionService
}

func NewDeleteTransactionHandler(service interfaces.TransactionService) *DeleteTransactionHandler {
	return &DeleteTransactionHandler{service: service}
}

// @Summary Delete a transaction
// @Description Delete an existing transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 204 "Transaction deleted successfully"
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Transaction not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/transactions/{id} [delete]
func (h *DeleteTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	if err := h.service.DeleteTransaction(r.Context(), id); err != nil {
		problem := common.NewInternalErrorProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
