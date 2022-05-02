DROP TABLE  IF EXISTS transaction;
CREATE TABLE transactions (
    id INT AUTO_INCREMENT primary key NOT NULL,
	type varchar(255) NOT NULL,
	from_wallet_id int,
	to_wallet_id int,
	ref_id varchar(255),
	amount DECIMAL(19, 4),
	currency varchar(255) NOT NULL DEFAULT 'USD',
	state varchar(255) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO transactions(type, ref_id, to_wallet_id, amount, state) VALUES ('DEPOSIT', 'b1', 1, '100', 'COMPLETED');
INSERT INTO transactions(type, ref_id, to_wallet_id, amount, state) VALUES ('DEPOSIT', 'b2', 2, '200', 'COMPLETED');