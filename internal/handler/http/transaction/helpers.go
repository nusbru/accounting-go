package transaction

import (
	"strings"

	"accounting/internal/domain/entity"
)

func toTransactionResponse(transaction *entity.Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		Amount:      transaction.Amount,
		Currency:    transaction.Currency,
		Description: transaction.Description,
		Category:    transaction.Category,
		Type:        transaction.Type,
		Date:        transaction.Date,
	}
}

func extractID(path, prefix string) string {
	path = strings.TrimPrefix(path, prefix)
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
