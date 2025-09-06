-- insert party and member
WITH new_party AS (
    INSERT INTO party DEFAULT VALUES
    RETURNING id AS party_id
), new_member AS (
    INSERT INTO member (party_id, name)
    SELECT party_id, 'Test name'
    FROM new_party
    RETURNING id AS member_id, party_id
), new_payment AS (
    INSERT INTO payment (party_id, amount, payer_id)
    SELECT party_id, 100, member_id
    FROM new_member
    RETURNING id as payment_id
)
INSERT INTO member_payment (member_id, payment_id)
SELECT nm.member_id, np.payment_id
FROM new_member AS nm
CROSS JOIN new_payment AS np;
