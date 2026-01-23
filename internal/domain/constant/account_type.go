package constant

type AccountType string

const (
	AccountTypeChecking   AccountType = "CHECKING"
	AccountTypeSavings    AccountType = "SAVINGS"
	AccountTypeCreditCard AccountType = "CREDIT_CARD"
	AccountTypeCash       AccountType = "CASH"
	AccountTypeInvestment AccountType = "INVESTMENT"
)
