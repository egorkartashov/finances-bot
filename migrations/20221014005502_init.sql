-- +goose Up
-- +goose StatementBegin


-- CURRENCIES
CREATE TABLE currencies
(
    iso text primary key
);

INSERT INTO currencies (iso)
VALUES ('RUB'),
       ('USD'),
       ('EUR'),
       ('CNY');


-- USERS
CREATE TABLE users
(
    id bigint PRIMARY KEY,
    currency text REFERENCES currencies (iso)
);

CREATE TABLE month_limits
(
    user_id bigint REFERENCES users(id),
    category text not null,
    sum decimal(12, 6) not null check (sum > 0),
    PRIMARY KEY (user_id, category)
);

-- EXCHANGE RATES
CREATE TABLE exchange_rates_to_rub
(
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    rate decimal(12, 6) NOT NULL,
    from_currency text REFERENCES currencies (iso) NOT NULL,
    date date NOT NULL,
    UNIQUE (from_currency, date)
);


-- EXPENSES
CREATE TABLE expenses
(
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id int REFERENCES users (id) NOT NULL,
    category text NOT NULL,
    sum_rub decimal(18, 6),
    date date NOT NULL
);

-- Используем btree индекс по user_id и date,
-- потому что при построении отчета нам нужно взять траты для конкретного юзера с датой >= заданной и просуммировать.
-- Такой индекс позволит ускорить извлечение этих трат, т.к. ускоряет поиск через сравнение (=, >=)
CREATE INDEX ON expenses USING btree (user_id, date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE expenses;
DROP TABLE exchange_rates_to_rub;
DROP TABLE users;
DROP TABLE month_limits;
DROP TABLE currencies;

-- +goose StatementEnd
