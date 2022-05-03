DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS wallet_assets;
CREATE TABLE wallets (
     id INT AUTO_INCREMENT primary key NOT NULL,
	user_id int NOT NULL
);

CREATE TABLE wallet_assets (
   id INT AUTO_INCREMENT primary key NOT NULL,
	wallet_id int NOT NULL,
	currency varchar(255) NOT NULL DEFAULT 'USD',
	amount DECIMAL(19, 4),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE cards (
   id INT AUTO_INCREMENT primary key NOT NULL,
   wallet_id int NOT NULL,
   number varchar(255) NOT NULL,
   bank_id int NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO wallets(user_id) VALUES(1);
INSERT INTO wallets(user_id) VALUES(2);

INSERT INTO wallet_assets(wallet_id, amount) VALUES(1, '100');
INSERT INTO wallet_assets(wallet_id, amount) VALUES(2, '200');

INSERT INTO cards(wallet_id, number, bank_id) VALUES(1, '1234', 1);
INSERT INTO cards(wallet_id, number, bank_id) VALUES(2, '4321', 1);