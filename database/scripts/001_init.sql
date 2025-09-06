DROP TABLE IF EXISTS party;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS member_payment;

CREATE TABLE party (
    id SERIAL PRIMARY KEY
);

CREATE TABLE member (
    id SERIAL PRIMARY KEY,
    party_id INTEGER REFERENCES party (id),
    name TEXT NOT NULL
);

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    party_id INTEGER REFERENCES party (id),
    amount INTEGER NOT NULL CHECK (amount > 0),
    payer_id INTEGER REFERENCES member (id)
);

CREATE TABLE member_payment (
    member_id INTEGER REFERENCES member (id),
    payment_id INTEGER REFERENCES payment (id),
    PRIMARY KEY (member_id, payment_id)
);