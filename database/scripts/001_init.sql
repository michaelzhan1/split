DROP TABLE IF EXISTS party;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS member_payment;

CREATE TABLE party (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE member (
    id SERIAL PRIMARY KEY,
    party_id INTEGER REFERENCES party (id)
        ON DELETE CASCADE,
    name TEXT NOT NULL
);

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    party_id INTEGER REFERENCES party (id)
        ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    payer_id INTEGER REFERENCES member (id)
        ON DELETE RESTRICT
);

CREATE TABLE member_payment (
    member_id INTEGER REFERENCES member (id) ON DELETE RESTRICT,
    payment_id INTEGER REFERENCES payment (id) ON DELETE CASCADE,
    PRIMARY KEY (member_id, payment_id)
);

-- payment has to have at least 1 member associated
CREATE OR REPLACE FUNCTION check_payment_has_members()
RETURNS TRIGGER AS
$$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM member_payment WHERE payment_id = NEW.id
    ) THEN
        RAISE EXCEPTION 'Payment % must have at least one associated member', NEW.id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE CONSTRAINT TRIGGER ensure_payment_has_members
AFTER INSERT OR UPDATE ON payment
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION check_payment_has_members();
