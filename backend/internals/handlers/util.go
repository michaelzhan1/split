package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/michaelzhan1/split/internals/database"
)

func withPartyId(r *http.Request) (int, *HttpError) {
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

func withMemberId(r *http.Request) (int, *HttpError) {
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
		res = append(res, Payment{
			ID:          payment.ID,
			Description: payment.Description,
			Amount:      payment.Amount,
			Payer:       payment.Payer,
			Payees:      payment.Payees,
		})
	}
	return res
}
