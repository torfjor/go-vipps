package ecom

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

// HandleConsentRemoval returns a convenience http.HandlerFunc for receiving
// requests for user consent removals from Vipps. `cb` is called with the uid
// of the user to that wishes to have its consents and data removed.
func HandleConsentRemoval(cb func(uid string)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.Header().Set("Allow", "DELETE")
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
		uid := path.Base(r.URL.Path)
		cb(uid)
	}

	return fn
}

// HandleShippingDetails returns a convenience http.HandlerFunc for responding
// to requests from Vipps to calculate shipping costs for an order.
//
// The provided authToken, if not empty, will be matched with the
// `Authorization` header of the incoming requests. If they don't match, the
// request will fail.
func HandleShippingDetails(authToken string, cb func(orderId string, req ShippingCostRequest) (ShippingCostResponse, error)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("ALLOW", "POST")
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
		if authToken != "" && r.Header.Get("Authorization") != authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		pathBySegments := strings.Split(r.URL.Path, "/")
		orderId := pathBySegments[len(pathBySegments)-2]

		bodyDec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		req := ShippingCostRequest{}
		err := bodyDec.Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sh, err := cb(orderId, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(sh)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(j)
	}
	return fn
}

// HandleTransactionUpdate returns a convenience http.HandlerFunc for getting
// notifications from Vipps about transaction updates.
//
// The provided authToken, if not empty, will be matched with the
// `Authorization` header of the incoming requests. If they don't match, the
// request will fail.
//
// `cb` will be called with the TransactionUpdate. Please be aware that for
// payments of PaymentTypeRegular, the fields `ShippingDetails` and
// `UserDetails` will be nil.
func HandleTransactionUpdate(authToken string, cb func(t TransactionUpdate)) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("ALLOW", "POST")
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}
		if authToken != "" && r.Header.Get("Authorization") != authToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var t TransactionUpdate
		bodyDec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err := bodyDec.Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cb(t)
	}
	return fn
}
