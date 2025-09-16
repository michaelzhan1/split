# Split

This is a clone of splitwise, which lets members in a group sort out how much they owe each other based on a set of shared purchases.

## Repo structure
This repo is built on a React frontend in `frontend/` and a go server in `backend/`.

## Sample Calls
```bash
# Parties
curl -s -X POST localhost:3000/parties -H "Content-Type: application/json" -d '{"name": "Trip to Vegas"}' | jq
curl -s localhost:3000/parties/2 | jq
curl -s -X PATCH localhost:3000/parties/2 -H "Content-Type: application/json" -d '{"name": "Trip to New York"}' | jq
curl -s -X DELETE localhost:3000/parties/2

# Members
curl -s localhost:3000/parties/2/members | jq
curl -s -X POST localhost:3000/parties/2/members -H "Content-Type: application/json" -d '{"name": "Alice"}' | jq
curl -s -X PATCH localhost:3000/parties/2/members/2 -H "Content-Type: application/json" -d '{"name": "Bob"}' | jq
curl -s -X DELETE localhost:3000/parties/2/members/2

# Payments
curl -s localhost:3000/parties/2/payments | jq
curl -s -X POST localhost:3000/parties/2/payments -H "Content-Type: application/json" -d '{"amount": 100, "description": "Hotel", "payer_id": 2, "payee_ids": [2,3,4]}' | jq
curl -s -X PATCH localhost:3000/parties/2/payments/2 -H "Content-Type: application/json" -d '{"amount": 150, "description": "Dinner"}' | jq
curl -s -X DELETE localhost:3000/parties/2/payments/2
curl -s -X DELETE localhost:3000/parties/2/payments

# Calculate
curl -s -X POST localhost:3000/parties/2/calculate | jq
```