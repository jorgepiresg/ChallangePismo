CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id uuid DEFAULT uuid_generate_v4 (),
    account_id  VARCHAR NOT NULL,
    operation_type_id INT NOT NULL,
    amount FLOAT NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (transaction_id)
);


