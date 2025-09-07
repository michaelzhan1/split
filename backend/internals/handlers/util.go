package handlers

import "github.com/michaelzhan1/split/internals/database"

func toPartyView(party database.Party) Party {
	return Party{Name: party.Name}
}

func toMemberList(members []database.Member) []Member {
	res := make([]Member, 0, len(members))
	for _, member := range members {
		res = append(res, Member{Name: member.Name})
	}
	return res
}
