package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type UpdateTransactionHandler struct {
	service interfaces.TransactionService
}

func NewUpdateTransactionHandler(service interfaces.TransactionService) *UpdateTransactionHandler {
	return &UpdateTransactionHandler{service: service}
}

// @Summary Update a transaction
// @Description Update an existing transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param body body UpdateTransactionRequest true "Transaction update request"
// @Success 200 {object} TransactionResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Transaction not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/transactions/{id} [put]
func (h *UpdateTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	var req UpdateTransactionRequest
	if r.Body == nil {
		problem := common.NewBadRequestProblem("invalid request body", r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem := common.NewBadRequestProblem("invalid request body", r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	// Validate optional input fields - only validate fields that are provided
	var validationErrs []*common.ValidationError
	if req.Amount != 0 {
		if err := common.ValidatePositive(req.Amount, "amount"); err != nil {
			validationErrs = append(validationErrs, err)
		}
	}
	if req.Currency != "" {
		if err := common.ValidateCurrency(req.Currency, "currency"); err != nil {
			validationErrs = append(validationErrs, err)
		}
	}
	if req.Type != "" {
		if err := common.ValidateEnum(string(req.Type), []string{"INCOME", "EXPENSE", "TRANSFER"}, "type"); err != nil {
			validationErrs = append(validationErrs, err)
		}
	}

	validationErrors := common.CollectErrors(validationErrs...)
	if len(validationErrors) > 0 {
		problem := common.NewValidationProblemWithErrors(r.RequestURI, validationErrors)
		common.WriteProblem(w, problem)
		return
	}

	var date time.Time
	if req.Date != nil {
		date = *req.Date
	}

	transaction, err := h.service.UpdateTransaction(
		r.Context(),
		id,
		req.Amount,
		req.Currency,
		req.Description,
		req.Category,
		req.Type,
		date,
	)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			problem := common.NewNotFoundProblem(err.Error(), r.RequestURI)
			common.WriteProblem(w, problem)
			return
		}
		var invalidErr *domainerrors.ErrInvalidInput
		if errors.As(err, &invalidErr) {
			problem := common.NewValidationProblem(err.Error(), r.RequestURI)
			common.WriteProblem(w, problem)
			return
		}
		problem := common.NewInternalErrorProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	common.WriteJSON(w, http.StatusOK, toTransactionResponse(transaction))
}
