package transaction

import (
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type ListAccountTransactionsHandler struct {
	service interfaces.TransactionService
}

func NewListAccountTransactionsHandler(service interfaces.TransactionService) *ListAccountTransactionsHandler {
	return &ListAccountTransactionsHandler{service: service}
}

// @Summary List account transactions
// @Description Retrieve all transactions for a specific account
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountID path string true "Account ID"
// @Success 200 {array} TransactionResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Account not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/accounts/{accountID}/transactions [get]
func (h *ListAccountTransactionsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		problem := common.NewMethodNotAllowedProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	accountID := extractID(r.URL.Path, "/api/v1/accounts/")
	if accountID == "" {
		validationErrors := []common.ValidationError{
			{Field: "accountID", Message: "account ID is required"},
		}
		problem := common.NewValidationProblemWithErrors(r.RequestURI, validationErrors)
		common.WriteProblem(w, problem)
		return
	}

	// Validate ID is a valid UUID
	if validationErr := common.ValidateUUID(accountID, "accountID"); validationErr != nil {
		problem := common.NewValidationProblemWithErrors(r.RequestURI, []common.ValidationError{*validationErr})
		common.WriteProblem(w, problem)
		return
	}

	transactions, err := h.service.ListAccountTransactions(r.Context(), accountID)
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

	response := make([]*TransactionResponse, 0, len(transactions))
	for _, txn := range transactions {
		response = append(response, toTransactionResponse(txn))
	}

	common.WriteJSON(w, http.StatusOK, response)
}
