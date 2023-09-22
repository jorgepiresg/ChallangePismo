CREATE TABLE IF NOT EXISTS operations_type (
    operation_type_id INT NOT NULL,
    description  VARCHAR NOT NULL,
    operation INT NOT NULL,  
    PRIMARY KEY (operation_type_id)
);

INSERT INTO operations_type
    (operation_type_id, description, operation)
VALUES
    (1, 'COMPRA A VISTA', -1),
    (2, 'COMPRA PARCELADA', -1),
    (3, 'SAQUE', -1),
    (4, 'PAGAMENTO', 1)
ON CONFLICT (operation_type_id) DO NOTHING;

