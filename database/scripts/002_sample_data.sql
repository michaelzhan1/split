-- insert group and user
WITH new_group AS (
    INSERT INTO groups (name)
    VALUES ('Test group name')
    RETURNING id AS group_id
), new_user AS (
    INSERT INTO users (group_id, name)
    SELECT group_id, 'Test user name'
    FROM new_group
    RETURNING id AS user_id, group_id
), new_payment AS (
    INSERT INTO payment (group_id, description, amount, payer_id)
    SELECT group_id, 'test description', 100, user_id
    FROM new_user
    RETURNING id as payment_id
)
INSERT INTO users_payment (user_id, payment_id)
SELECT nu.user_id, np.payment_id
FROM new_user AS nu
CROSS JOIN new_payment AS np;

INSERT INTO groups (name)
VALUES ('Another test group name');
