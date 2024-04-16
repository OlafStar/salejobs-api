package api

import (
	"encoding/json"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"log"
	"net/http"
)

var PriceId = "price_1P5z1OLphFr90yTIQimY2ww7"

func checkout(email string) (*stripe.CheckoutSession, error) {
	// Creating a customer is optional, you might want to create it to attach to the session for tracking
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	newCustomer, err := customer.New(customerParams)
	if err != nil {
		return nil, err
	}

	log.Printf("New customer created: %s", newCustomer.ID)

	params := &stripe.CheckoutSessionParams{
		Customer:           &newCustomer.ID,
		SuccessURL:         stripe.String("https://www.example.com/success"),
		CancelURL:          stripe.String("https://www.example.com/cancel"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(PriceId),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	return session.New(params)
}

type EmailInput struct {
	Email string `json:"email"`
}

type SessionOutput struct {
	Id string `json:"id"`
}

func CheckoutCreator(w http.ResponseWriter, req *http.Request) error {
	input := &EmailInput{}
	err := json.NewDecoder(req.Body).Decode(input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return err
	}

	stripeSession, err := checkout(input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	err = json.NewEncoder(w).Encode(&SessionOutput{Id: stripeSession.ID})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}

	return nil
}
