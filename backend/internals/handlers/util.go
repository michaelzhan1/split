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

// resolve balances
func calculate(members []database.Member) []IOU {
	pos := []Member{}
	neg := []Member{}

	for _, member := range members {
		if member.Balance > 0 {
			pos = append(pos, Member{
				ID:      member.ID,
				Name:    member.Name,
				Balance: member.Balance,
			})
		} else if member.Balance < 0 {
			neg = append(neg, Member{
				ID:      member.ID,
				Name:    member.Name,
				Balance: -member.Balance,
			})
		}
	}


	ious := []IOU{}
	i, j := 0, 0
	for i < len(pos) && j < len(neg) {
		minAmount := pos[i].Balance
		if neg[j].Balance < minAmount {
			minAmount = neg[j].Balance
		}

		ious = append(ious, IOU{
			FromID: neg[j].ID,
			ToID:   pos[i].ID,
			Amount: minAmount,
		})

		pos[i].Balance -= minAmount
		neg[j].Balance -= minAmount

		if pos[i].Balance == 0 {
			i++
		}
		if neg[j].Balance == 0 {
			j++
		}
	}

	return ious
}
