DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS users_payment;

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups (id)
        ON DELETE CASCADE,
    name TEXT NOT NULL,
    balance NUMERIC NOT NULL DEFAULT 0
);

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    group_id INTEGER REFERENCES groups (id)
        ON DELETE CASCADE,
    description TEXT NOT NULL,
    amount NUMERIC NOT NULL CHECK (amount > 0),
    payer_id INTEGER REFERENCES users (id)
        ON DELETE RESTRICT
);

CREATE TABLE users_payment (
    user_id INTEGER REFERENCES users (id) ON DELETE RESTRICT,
    payment_id INTEGER REFERENCES payment (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, payment_id)
);

-- payment has to have at least 1 user associated
CREATE OR REPLACE FUNCTION check_payment_has_users()
RETURNS TRIGGER AS
$$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM users_payment WHERE payment_id = NEW.id
    ) THEN
        RAISE EXCEPTION 'Payment % must have at least one associated user', NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER ensure_payment_has_users
AFTER INSERT OR UPDATE ON payment
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION check_payment_has_users();

-- payees and payer have to be in the same group as the payment
CREATE OR REPLACE FUNCTION check_payment_users_in_same_group()
RETURNS TRIGGER AS 
$$
DECLARE
    payer_group_id INTEGER;
    payee_group_id INTEGER;
BEGIN
    -- check payer is in the same group as the payment
    SELECT group_id INTO payer_group_id FROM users WHERE id = NEW.payer_id;

    IF payer_group_id IS NULL THEN
        RAISE EXCEPTION 'Payer % does not exist', NEW.payer_id;
    END IF;

    IF payer_group_id != NEW.group_id THEN
        RAISE EXCEPTION 'Payer % is not in the same group as the payment (group %)', NEW.payer_id, NEW.group_id;
    END IF;

    -- check all associated payees are in the same group
    FOR payee_group_id IN
        SELECT u.group_id
        FROM users_payment up
        JOIN users u ON u.id = up.user_id
        WHERE up.payment_id = NEW.id
    LOOP
        IF payee_group_id != NEW.group_id THEN
            RAISE EXCEPTION 'One or more users in users_payment are not in the same group as the payment (group %)', NEW.group_id;
        END IF;
    END LOOP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER ensure_payment_users_in_same_group
AFTER INSERT OR UPDATE ON payment
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION check_payment_users_in_same_group()