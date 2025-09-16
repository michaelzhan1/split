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

func withMemberID(r *http.Request) (int, *HttpError) {
	memberIDStr := chi.URLParam(r, "member_id")
	if memberIDStr == "" {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Empty or missing member ID",
		}
	}
	memberIDInt, err := strconv.Atoi(memberIDStr)
	if err != nil || memberIDInt <= 0 {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Bad member ID",
		}
	}
	return memberIDInt, nil
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

func toMemberList(members []database.Member) []Member {
	res := make([]Member, 0, len(members))
	for _, member := range members {
		res = append(res, Member{
			ID:      member.ID,
			Name:    member.Name,
			Balance: member.Balance,
		})
	}
	return res
}

func toPaymentList(payments []database.Payment) []Payment {
	res := make([]Payment, 0, len(payments))
	for _, payment := range payments {
		payees := []Member{}
		for idx := range payment.PayeeIDs {
			payees = append(payees, Member{
				ID:      payment.PayeeIDs[idx],
				Name:    payment.PayeeNames[idx],
				Balance: payment.PayeeBalances[idx],
			})
		}

		res = append(res, Payment{
			ID:          payment.ID,
			Description: payment.Description,
			Amount:      payment.Amount,
			Payer: Member{
				ID:      payment.PayerID,
				Name:    payment.PayerName,
				Balance: payment.PayerBalance,
			},
			Payees: payees,
		})
	}
	return res
}
