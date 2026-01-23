package router

import (
	"net/http"
	"strings"

	"accounting/internal/handler/http/account"
	"accounting/internal/handler/http/transaction"
	"accounting/internal/handler/http/user"
	"accounting/internal/service"
)

// Router holds all route handlers
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new router with all routes configured
func NewRouter(
	userService *service.UserService,
	accountService *service.AccountService,
	transactionService *service.TransactionService,
) *Router {
	mux := http.NewServeMux()

	// User handlers
	createUserHandler := user.NewCreateUserHandler(userService)
	updateUserHandler := user.NewUpdateUserHandler(userService)
	deleteUserHandler := user.NewDeleteUserHandler(userService)
	getUserHandler := user.NewGetUserHandler(userService)
	getUserByEmailHandler := user.NewGetUserByEmailHandler(userService)

	// Account handlers
	createAccountHandler := account.NewCreateAccountHandler(accountService)
	updateAccountHandler := account.NewUpdateAccountHandler(accountService)
	deleteAccountHandler := account.NewDeleteAccountHandler(accountService)
	getAccountHandler := account.NewGetAccountHandler(accountService)
	listUserAccountsHandler := account.NewListUserAccountsHandler(accountService)

	// Transaction handlers
	createTransactionHandler := transaction.NewCreateTransactionHandler(transactionService)
	updateTransactionHandler := transaction.NewUpdateTransactionHandler(transactionService)
	deleteTransactionHandler := transaction.NewDeleteTransactionHandler(transactionService)
	getTransactionHandler := transaction.NewGetTransactionHandler(transactionService)
	listAccountTransactionsHandler := transaction.NewListAccountTransactionsHandler(transactionService)

	// User routes
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createUserHandler.Handle(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/users/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUserByEmailHandler.Handle(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		// Handle /api/v1/users/{userId}/accounts
		if strings.HasSuffix(r.URL.Path, "/accounts") && r.Method == http.MethodGet {
			listUserAccountsHandler.Handle(w, r)
			return
		}

		// Handle /api/v1/users/{id}
		switch r.Method {
		case http.MethodGet:
			getUserHandler.Handle(w, r)
		case http.MethodPut:
			updateUserHandler.Handle(w, r)
		case http.MethodDelete:
			deleteUserHandler.Handle(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Account routes
	mux.HandleFunc("/api/v1/accounts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createAccountHandler.Handle(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/accounts/", func(w http.ResponseWriter, r *http.Request) {
		// Handle /api/v1/accounts/{accountId}/transactions
		if strings.HasSuffix(r.URL.Path, "/transactions") && r.Method == http.MethodGet {
			listAccountTransactionsHandler.Handle(w, r)
			return
		}

		// Handle /api/v1/accounts/{id}
		switch r.Method {
		case http.MethodGet:
			getAccountHandler.Handle(w, r)
		case http.MethodPut:
			updateAccountHandler.Handle(w, r)
		case http.MethodDelete:
			deleteAccountHandler.Handle(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Transaction routes
	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createTransactionHandler.Handle(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTransactionHandler.Handle(w, r)
		case http.MethodPut:
			updateTransactionHandler.Handle(w, r)
		case http.MethodDelete:
			deleteTransactionHandler.Handle(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return &Router{mux: mux}
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
