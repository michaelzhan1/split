package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/michaelzhan1/split/internals/database"
)

func withPartyID(r *http.Request) (int, *HttpError) {
	partyIDStr := chi.URLParam(r, "party_id")
	if partyIDStr == "" {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Empty or missing party ID",
		}
	}
	partyIDInt, err := strconv.Atoi(partyIDStr)
	if err != nil || partyIDInt <= 0 {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Bad party ID",
		}
	}
	return partyIDInt, nil
}

func withUserID(r *http.Request) (int, *HttpError) {
	userIDStr := chi.URLParam(r, "user_id")
	if userIDStr == "" {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Empty or missing user ID",
		}
	}
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil || userIDInt <= 0 {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Bad user ID",
		}
	}
	return userIDInt, nil
}

func withPaymentID(r *http.Request) (int, *HttpError) {
	paymentIDStr := chi.URLParam(r, "payment_id")
	if paymentIDStr == "" {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Empty or missing payment ID",
		}
	}
	paymentIDInt, err := strconv.Atoi(paymentIDStr)
	if err != nil || paymentIDInt <= 0 {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Bad payment ID",
		}
	}
	return paymentIDInt, nil
}

func toPartyView(party database.Party) Party {
	return Party{
		ID:   party.ID,
		Name: party.Name,
	}
}

func toUserList(users []database.User) []User {
	res := make([]User, 0, len(users))
	for _, user := range users {
		res = append(res, User{
			ID:      user.ID,
			Name:    user.Name,
			Balance: user.Balance,
		})
	}
	return res
}

func toPaymentList(payments []database.Payment) []Payment {
	res := make([]Payment, 0, len(payments))
	for _, payment := range payments {
		payees := []User{}
		for idx := range payment.PayeeIDs {
			payees = append(payees, User{
				ID:      payment.PayeeIDs[idx],
				Name:    payment.PayeeNames[idx],
				Balance: payment.PayeeBalances[idx],
			})
		}

		res = append(res, Payment{
			ID:          payment.ID,
			Description: payment.Description,
			Amount:      payment.Amount,
			Payer: User{
				ID:      payment.PayerID,
				Name:    payment.PayerName,
				Balance: payment.PayerBalance,
			},
			Payees: payees,
		})
	}
	return res
}
