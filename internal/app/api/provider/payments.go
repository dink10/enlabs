package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gookit/validate"

	"github.com/dink10/enlabs/internal/pkg/logger"
	"github.com/dink10/enlabs/internal/pkg/payments"
	"github.com/dink10/enlabs/internal/pkg/server"
)

// PaymentsProvider provides endpoints to interact with PAYMENT.Service.
type PaymentsProvider struct {
	service *payments.Service
	logger  *logger.ProviderLogger
}

// NewPaymentProvider returns a new instance of PaymentsProvider.
func NewPaymentProvider(service *payments.Service) PaymentsProvider {
	return PaymentsProvider{
		service: service,
		logger:  logger.NewProviderLogger("payments"),
	}
}

// Router returns PaymentsProvider router.
func (p *PaymentsProvider) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(p.userMiddleware)
		r.Post("/", p.create)
		r.Get("/balance", p.balance)
	})

	return r
}

// Request body format for inserting/updating a payment.
type paymentRequest struct {
	State         string `json:"state" validate:"stateValidator"`
	TransactionId string `json:"transactionId" validate:"required"`
	Amount        string `json:"amount" validate:"required"`
}

// StateValidator source validator in the source struct.
func (pr paymentRequest) StateValidator(val string) bool {
	switch val {
	case "win", "lost":
		return true
	default:
		return false
	}
}

type paymentsResponse struct {
	*server.Response
	Message string `json:"message"`
}

func (p *PaymentsProvider) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// it could be middleware with user recognition, now - just hardcoded
		ctx := context.WithValue(r.Context(), "account_id", 1)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// @Summary Payment processing
// @Description Process payment in database
// @ID payment-create
// @Tags Account
// @Accept json
// @Produce json
// @Param payment body provider.paymentRequest true "Payment data to create"
// @Header 200 {string} Source-Type "payment"
// @Param Source-Type header string true "With the bearer started"
// @Content-Type application/json
// @Success 200 {object} provider.paymentsResponse "Proceeded payment"
// @Failure 400 {object} server.ErrorResponse "Invalid Request"
// @Failure 404 {object} server.ErrorResponse "Account Not Found"
// @Failure 500 {object} server.ErrorResponse "Service Error"
// @Router /v1/payments [post]
func (p *PaymentsProvider) create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(server.HeaderContentType) != server.JsonContentType {
		p.logger.Logger(r).Errorf(
			"incorrect Content-Type, required %s, got %s",
			server.JsonContentType,
			r.Header.Get(server.HeaderContentType),
		)
		server.RenderResponse(
			w, r, server.NewErrorResponse(http.StatusBadRequest, fmt.Errorf(
				"incorrect Content-Type, required %s, got %s",
				server.JsonContentType,
				r.Header.Get(server.HeaderContentType)),
			),
		)
		return
	}

	var paymentRequest paymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		p.logger.Logger(r).Errorf("failed to decode body: %v", err)
		server.RenderResponse(w, r,
			server.NewErrorResponse(http.StatusBadRequest, fmt.Errorf("failed to decode request payload")),
		)
		return
	}

	v := validate.Struct(paymentRequest)
	if !v.Validate() {
		p.logger.Logger(r).Errorf("invalid body data: %v", v.Errors)
		server.RenderResponse(w, r, server.NewErrorResponse(http.StatusBadRequest, v.Errors))
		return
	}

	amount, err := strconv.ParseFloat(paymentRequest.Amount, 64)
	if err != nil || amount < 0.0 {
		p.logger.Logger(r).Errorf("incorrect amount value: %v", err)
		server.RenderResponse(w, r,
			server.NewErrorResponse(http.StatusBadRequest, fmt.Errorf("incorrect amount value")),
		)
		return
	}

	sourceType, err := p.service.SourceTypeID(r.Context(), r.Header.Get(server.HeaderSourceType))
	if err != nil {
		p.logger.Logger(r).Errorf("%v", err)
		server.RenderResponse(w, r, server.NewErrorResponse(http.StatusBadRequest, err))
		return
	}

	accountID, ok := r.Context().Value("account_id").(int)
	if !ok {
		p.logger.Logger(r).Error("wrong account_id")
		server.RenderResponse(w, r, server.NewErrorResponse(http.StatusBadRequest, fmt.Errorf("wrong account_id")))
		return
	}

	payment := payments.Payment{
		AccountID:     accountID,
		TransactionID: paymentRequest.TransactionId,
		State:         paymentRequest.State,
		Amount:        amount,
		SourceType:    sourceType,
		Processed:     true,
	}

	if err = p.service.ProceedPayment(r.Context(), payment); err != nil {
		p.logger.Logger(r).Error(err)
		server.RenderResponse(w, r, server.NewErrorResponse(http.StatusBadRequest, err))
		return
	}

	server.RenderResponse(w, r, &paymentsResponse{
		Response: server.NewResponse(http.StatusOK),
		Message:  "payment was successfully proceed",
	})
}

// @Summary Account Balance
// @Description Balance
// @ID show-balance
// @Tags Balance
// @Produce json
// @Content-Type application/json
// @Success 200 {object} provider.paymentsResponse "Balance"
// @Failure 400 {object} server.ErrorResponse "Invalid Request"
// @Failure 404 {object} server.ErrorResponse "Not Found"
// @Failure 500 {object} server.ErrorResponse "Service Error"
// @Router /v1/payments/balance [get]
func (p *PaymentsProvider) balance(w http.ResponseWriter, r *http.Request) {
	balance, err := p.service.Balance(r.Context())
	if err != nil {
		p.logger.Logger(r).Error(err)
		server.RenderResponse(w, r, server.NewErrorResponse(http.StatusInternalServerError, err))
		return
	}

	server.RenderResponse(w, r, &paymentsResponse{
		Response: server.NewResponse(http.StatusOK),
		Message:  fmt.Sprintf("%f", balance.Balance),
	})
}
