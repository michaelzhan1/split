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
    description TEXT,
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
