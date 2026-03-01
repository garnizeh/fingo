-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
user_id       UUID        NOT NULL,
name          TEXT        NOT NULL,
email         TEXT UNIQUE NOT NULL,
roles         TEXT[]      NOT NULL,
password_hash TEXT        NOT NULL,
    department    TEXT        NULL,
    enabled       BOOLEAN     NOT NULL,
date_created  TIMESTAMP   NOT NULL,
date_updated  TIMESTAMP   NOT NULL,

PRIMARY KEY (user_id)
);

-- Version: 1.02
-- Description: Create table products
CREATE TABLE products (
product_id   UUID           NOT NULL,
user_id      UUID           NOT NULL,
name         TEXT           NOT NULL,
cost         NUMERIC(10, 2) NOT NULL,
quantity     INT            NOT NULL,
date_created TIMESTAMP      NOT NULL,
date_updated TIMESTAMP      NOT NULL,

PRIMARY KEY (product_id),
FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.03
-- Description: Add products view.
CREATE OR REPLACE VIEW view_products AS
SELECT
    p.product_id,
    p.user_id,
p.name,
    p.cost,
p.quantity,
    p.date_created,
    p.date_updated,
    u.name AS user_name
FROM
    products AS p
JOIN
    users AS u ON u.user_id = p.user_id;

-- Version: 1.04
-- Description: Create table homes
CREATE TABLE homes (
    home_id       UUID       NOT NULL,
    type          TEXT       NOT NULL,
    user_id       UUID       NOT NULL,
    address_1     TEXT       NOT NULL,
    address_2     TEXT       NULL,
    zip_code      TEXT       NOT NULL,
    city          TEXT       NOT NULL,
    state         TEXT       NOT NULL,
    country       TEXT       NOT NULL,
    date_created  TIMESTAMP  NOT NULL,
    date_updated  TIMESTAMP  NOT NULL,

    PRIMARY KEY (home_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.05
-- Description: Create table audit
CREATE TABLE audit (
    id          UUID      NOT NULL,
    obj_id      UUID      NOT NULL,
    obj_domain  TEXT      NOT NULL,
    obj_name    TEXT      NOT NULL,
    actor_id    UUID      NOT NULL,
    action      TEXT      NOT NULL,
    data        JSONB     NULL,
    message     TEXT      NULL,
    timestamp   TIMESTAMP NOT NULL,

    PRIMARY KEY (id)
);

-- Version: 2.01
-- Description: Create table credit_cards
CREATE TABLE IF NOT EXISTS credit_cards (
    credit_card_id   UUID           NOT NULL,
    user_id          UUID           NOT NULL REFERENCES users(user_id),
    name             TEXT           NOT NULL,
    card_limit       NUMERIC(12, 2) NOT NULL,
    closing_day      INT            NOT NULL,
    due_day          INT            NOT NULL,
    last_four_digits CHAR(4)        NOT NULL,
    enabled          BOOLEAN        NOT NULL DEFAULT true,
    date_created     TIMESTAMPTZ    NOT NULL,
    date_updated     TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (credit_card_id),
    CONSTRAINT chk_closing_day CHECK (closing_day BETWEEN 1 AND 31),
    CONSTRAINT chk_due_day     CHECK (due_day BETWEEN 1 AND 31)
);

CREATE INDEX IF NOT EXISTS idx_credit_cards_user ON credit_cards (user_id);

-- Version: 2.02
-- Description: Create table invoices
CREATE TABLE IF NOT EXISTS invoices (
    invoice_id       UUID           NOT NULL,
    credit_card_id   UUID           NOT NULL REFERENCES credit_cards(credit_card_id),
    reference_month  DATE           NOT NULL,
    total_amount     NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status           TEXT           NOT NULL DEFAULT 'open',
    due_date         TIMESTAMPTZ    NOT NULL,
    date_created     TIMESTAMPTZ    NOT NULL,
    date_updated     TIMESTAMPTZ    NOT NULL,

    PRIMARY KEY (invoice_id),
    CONSTRAINT uq_invoice_card_month UNIQUE (credit_card_id, reference_month)
);
