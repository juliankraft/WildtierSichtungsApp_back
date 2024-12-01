import mysql.connector
import csv
import json

# Load the database connection details from the JSON file
with open('./db_setup/db_admin.json', 'r') as config_file:
    config = json.load(config_file)

# Connect to the MariaDB database
connection = mysql.connector.connect(
    host=config['host'],
    port=config['port'],
    user=config['user'],
    password=config['password'],
    database=config['database']
)

# Purge existing tables
def purge_tables():
    drop_sichtungen_table = "DROP TABLE IF EXISTS sichtungen;"
    drop_tierarten_table = "DROP TABLE IF EXISTS tierarten;"
    drop_users_table = "DROP TABLE IF EXISTS users;"

    cursor = connection.cursor()
    try:
        cursor.execute(drop_sichtungen_table)
        cursor.execute(drop_tierarten_table)
        cursor.execute(drop_users_table)
        connection.commit()
        print("Existing tables purged.")
    except mysql.connector.Error as e:
        print(f"Error purging tables: {e}")
        connection.rollback()
    finally:
        cursor.close()

# Create tables
def create_tables():
    create_users_table = """
    CREATE TABLE users (
        user_id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL UNIQUE,
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    """
    create_tierarten_table = """
    CREATE TABLE tierarten (
        tierart_id INT AUTO_INCREMENT PRIMARY KEY,
        tierart VARCHAR(255) NOT NULL UNIQUE,
        wissenschaftlicher_name VARCHAR(255) NOT NULL,
        familie VARCHAR(255) NOT NULL,
        gattung VARCHAR(255) NOT NULL,
        art VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    """
    create_sichtungen_table = """
    CREATE TABLE sichtungen (
    sichtungen_id INT AUTO_INCREMENT PRIMARY KEY,
    -- user_id INT NOT NULL,
    tierart_id INT NOT NULL,
    anzahl_maennlich INT,
    anzahl_weiblich INT,
    anzahl_unbekannt INT,
    sichtung_date DATETIME NOT NULL,
    sichtung_location GEOMETRY NOT NULL,
    sichtung_bemerkung TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (tierart_id) REFERENCES tierarten(tierart_id)
    );
    """
    cursor = connection.cursor()
    try:
        cursor.execute(create_users_table)
        cursor.execute(create_tierarten_table)
        cursor.execute(create_sichtungen_table)
        connection.commit()
        print("Tables created successfully.")
    except mysql.connector.Error as e:
        print(f"Error creating tables: {e}")
        connection.rollback()
    finally:
        cursor.close()

# Insert data into tierarten from CSV
def insert_tierarten():
    insert_tierarten_query = """
    INSERT INTO tierarten (tierart, wissenschaftlicher_name, familie, gattung, art)
    VALUES (%s, %s, %s, %s, %s)
    ON DUPLICATE KEY UPDATE tierart=VALUES(tierart);
    """
    cursor = connection.cursor()
    try:
        with open('./db_setup/tierarten.csv', 'r', encoding='utf-8') as csvfile:
            csvreader = csv.reader(csvfile)
            next(csvreader)  # Skip the header row
            data = [tuple(row) for row in csvreader]

        cursor.executemany(insert_tierarten_query, data)
        connection.commit()
        print(f"Inserted or updated {len(data)} rows in the tierarten table.")
    except mysql.connector.Error as e:
        print(f"Error inserting data into tierarten: {e}")
        connection.rollback()
    finally:
        cursor.close()

if __name__ == "__main__":
    try:
        purge_tables()      # Drop existing tables
        create_tables()     # Create new tables
        insert_tierarten()  # Insert tierarten data
    except any as e:
        print(f"An error occurred: {e}")
    finally:
        print("DB setup completed")
        connection.close()  # Ensure connection is always closed
