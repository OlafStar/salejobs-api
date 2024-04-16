package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/v72/webhook"
)

const StripeWebhookSecret = "whsec_VhR6nGuj7tHKsINmFJd3gekJ5ZoYsTTI"

func StripeWebhookHandler(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

	payload, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), StripeWebhookSecret)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		http.Error(w, "Error verifying webhook signature", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing payment intent: %v", err)
			http.Error(w, "Error parsing payment intent", http.StatusInternalServerError)
			return
		}
		log.Printf("PaymentIntent was successful! PaymentIntent ID: %s", paymentIntent.ID)
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			log.Printf("Error parsing payment method: %v", err)
			http.Error(w, "Error parsing payment method", http.StatusInternalServerError)
			return
		}
		log.Printf("Payment method was attached to a customer! PaymentMethod ID: %s", paymentMethod.ID)
	// ... handle other event types
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

