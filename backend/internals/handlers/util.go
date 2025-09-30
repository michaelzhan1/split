package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/michaelzhan1/split/internals/database"
)

func withGroupID(r *http.Request) (int, *HttpError) {
	groupIDStr := chi.URLParam(r, "group_id")
	if groupIDStr == "" {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Empty or missing group ID",
		}
	}
	groupIDInt, err := strconv.Atoi(groupIDStr)
	if err != nil || groupIDInt <= 0 {
		return 0, &HttpError{
			Code:    http.StatusBadRequest,
			Message: "Bad group ID",
		}
	}
	return groupIDInt, nil
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

func toGroupView(group database.Group) Group {
	return Group{
		ID:   group.ID,
		Name: group.Name,
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

// resolve balances
func calculate(users []database.User) []IOU {
	pos := []User{}
	neg := []User{}

	for _, user := range users {
		if user.Balance > 0 {
			pos = append(pos, User{
				ID:      user.ID,
				Name:    user.Name,
				Balance: user.Balance,
			})
		} else if user.Balance < 0 {
			neg = append(neg, User{
				ID:      user.ID,
				Name:    user.Name,
				Balance: -user.Balance,
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
