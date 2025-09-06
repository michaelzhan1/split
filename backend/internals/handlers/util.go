package handlers

import "github.com/michaelzhan1/split/internals/database"

func toPartyView(party database.Party) GetPartyResponse {
	return GetPartyResponse{Name: party.Name}
}