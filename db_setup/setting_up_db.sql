CREATE TABLE users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

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
    anzahl_maennlich INT NOT NULL,
    anzahl_weiblich INT NOT NULL,
    anzahl_unbekannt INT NOT NULL,
    sichtung_date DATE NOT NULL,
    sichtung_location GEOMETRY NOT NULL,
    sichtung_bemerkung TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (tierart_id) REFERENCES tierarten(tierart_id)
);