CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    pwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (user_name, first_name, last_name, email, pwd)
VALUES ('kraftjul', 'Julian', 'Kraft', 'kraftjul@students.zhaw.ch', 'passwort');

CREATE TABLE tierarten (
    tierart_id INT AUTO_INCREMENT PRIMARY KEY,
    tierart VARCHAR(255) NOT NULL UNIQUE,
    wissenschaftlicher_name VARCHAR(255) NOT NULL,
    familie VARCHAR(255) NOT NULL,
    gattung VARCHAR(255) NOT NULL,
    art VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sichtungen (
    sichtungen_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    tierart_id INT NOT NULL,
    anzahl_maennlich INT,
    anzahl_weiblich INT,
    anzahl_unbekannt INT,
    sichtung_date DATE NOT NULL,
    sichtung_location GEOMETRY NOT NULL,
    sichtung_bemerkung TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (tierart_id) REFERENCES tierarten(tierart_id)
);