package main

import (
	"github.com/gorilla/mux"
	"github.com/torfjor/go-vipps/recurring"
	"github.com/unrolled/render"
	"net/http"
)

type handler struct {
	r     *render.Render
	vipps recurring.Client
}

func (h *handler) listAgreements(w http.ResponseWriter, r *http.Request) {
	var agreementStatus recurring.AgreementStatus
	ctx := r.Context()
	status := r.URL.Query().Get("status")
	switch status {
	case "active":
		agreementStatus = recurring.AgreementStatusActive
	case "stopped":
		agreementStatus = recurring.AgreementStatusStopped
	case "pending":
		agreementStatus = recurring.AgreementStatusPending
	case "expired":
		agreementStatus = recurring.AgreementStatusExpired
	default:
		agreementStatus = recurring.AgreementStatusActive
	}

	agreements, err := h.vipps.ListAgreements(ctx, agreementStatus)
	if err != nil {
		h.r.HTML(w, http.StatusInternalServerError, "error", err)
		return
	}
	h.r.HTML(w, http.StatusOK, "agreements", struct {
		Status     string
		Agreements []*recurring.Agreement
	}{status, agreements})
}

func (h *handler) getAgreement(w http.ResponseWriter, r *http.Request) {
	var statuses []recurring.ChargeStatus
	var chargeStatus recurring.ChargeStatus
	id := mux.Vars(r)["id"]
	ctx := r.Context()

	status := r.URL.Query().Get("status")
	switch status {
	case "pending":
		chargeStatus = recurring.ChargeStatusPending
	case "due":
		chargeStatus = recurring.ChargeStatusDue
	case "reserved":
		chargeStatus = recurring.ChargeStatusReserved
	case "charged":
		chargeStatus = recurring.ChargeStatusCharged
	case "cancelled":
		chargeStatus = recurring.ChargeStatusCancelled
	case "failed":
		chargeStatus = recurring.ChargeStatusFailed
	case "partial":
		chargeStatus = recurring.ChargeStatusPartiallyRefunded
	case "refunded":
		chargeStatus = recurring.ChargeStatusRefunded
	case "processing":
		chargeStatus = recurring.ChargeStatusProcessing
	default:
		chargeStatus = ""
	}

	if len(chargeStatus) > 0 {
		statuses = append(statuses, chargeStatus)
	}

	agreement := &recurring.Agreement{}
	charges := []*recurring.Charge{}

	a := make(chan *recurring.Agreement)
	c := make(chan []*recurring.Charge)
	errors := make(chan error)

	go func() {
		agreement, err := h.vipps.GetAgreement(ctx, id)
		if err != nil {
			errors <- err
		}
		a <- agreement
	}()

	go func() {
		charge, err := h.vipps.ListCharges(ctx, id, statuses...)
		if err != nil {
			errors <- err
		}
		c <- charge
	}()

	for i := 0; i < 2; i++ {
		select {
		case agreement = <-a:
		case charges = <-c:
		case e := <-errors:
			h.r.HTML(w, http.StatusInternalServerError, "error", e)
			return
		}
	}

	h.r.HTML(w, http.StatusOK, "agreement", struct {
		Agreement *recurring.Agreement
		Charges   []*recurring.Charge
	}{agreement, charges})
}
