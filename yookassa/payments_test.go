package yookassa

import (
	"io"
	"net/http"
	"strings"
	"testing"

	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newPaymentHandlerWithResponse(statusCode int, body string) *PaymentHandler {
	client := NewClient("account_id", "secret_key")
	client.client = http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	return NewPaymentHandler(client)
}

func TestCreatePaymentReturnsErrorWhenConfirmationIsEmptyAndPaymentMethodIDIsEmpty(t *testing.T) {
	handler := newPaymentHandlerWithResponse(http.StatusOK, `{"id":"payment-id","status":"pending"}`)

	_, err := handler.CreatePayment(&yoopayment.Payment{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "empty confirmation url" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreatePaymentDoesNotReturnErrorWhenConfirmationIsEmptyAndPaymentMethodIDIsSet(t *testing.T) {
	handler := newPaymentHandlerWithResponse(http.StatusOK, `{"id":"payment-id","status":"pending"}`)

	result, err := handler.CreatePayment(&yoopayment.Payment{PaymentMethodID: "pm-123"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("expected payment response, got nil")
	}

	if result.ID != "payment-id" {
		t.Fatalf("unexpected payment id: %s", result.ID)
	}
}
