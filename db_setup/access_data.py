import mysql.connector
import json

# Load the database connection details from the JSON file
with open('./db_setup/db_readonly_config.json', 'r') as config_file:
    config = json.load(config_file)

# Connect to the MariaDB database
connection = mysql.connector.connect(
    host=config['host'],
    port=config['port'],
    user=config['user'],
    password=config['password'],
    database=config['database']
)

# Function to fetch data from the database
def get_data(statement):
    cursor = connection.cursor()
    try:
        cursor.execute(statement)
        data = cursor.fetchall()
        return data
    except mysql.connector.Error as e:
        print(f"Error fetching data: {e}")
    finally:
        cursor.close()

# Fetching data example

data = get_data("SELECT * FROM sichtungen;")

print(data)