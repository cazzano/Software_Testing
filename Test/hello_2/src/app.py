from flask import Flask, jsonify, request
from flask_cors import CORS
import sqlite3
import os

app = Flask(__name__)
CORS(app)  # Enable CORS for all routes

# Database setup
DB_PATH = 'languages.db'

def init_db():
    # Create the database if it doesn't exist
    if not os.path.exists(DB_PATH):
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        
        # Create table
        cursor.execute('''
        CREATE TABLE languages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            color TEXT NOT NULL
        )
        ''')
        
        # Insert initial data
        initial_data = [
            ('Python', '4F75A1'),  # Blue
            ('JavaScript', 'F1DC5D'),  # Yellow
            ('Swift', 'EB735F'),  # Red/Orange
            ('Rust', '8E6CE1'),  # Purple
            ('Kotlin', '8E6CE1'),  # Purple
            ('Dart', '74B8DF'),  # Light Blue
            ('Go', '6BAFC6')  # Teal
        ]
        
        cursor.executemany('INSERT INTO languages (name, color) VALUES (?, ?)', initial_data)
        conn.commit()
        conn.close()

# Initialize the database
init_db()

# Routes
@app.route('/api/languages', methods=['GET'])
def get_languages():
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row  # This enables name-based access to columns
    cursor = conn.cursor()
    
    cursor.execute('SELECT * FROM languages')
    languages = [dict(row) for row in cursor.fetchall()]
    
    conn.close()
    return jsonify(languages)

@app.route('/api/languages/<int:language_id>', methods=['GET'])
def get_language(language_id):
    conn = sqlite3.connect(DB_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    cursor.execute('SELECT * FROM languages WHERE id = ?', (language_id,))
    language = cursor.fetchone()
    
    conn.close()
    
    if language:
        return jsonify(dict(language))
    return jsonify({"error": "Language not found"}), 404

@app.route('/api/languages', methods=['POST'])
def add_language():
    data = request.json
    
    if not data or 'name' not in data or 'color' not in data:
        return jsonify({"error": "Name and color are required"}), 400
    
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    cursor.execute(
        'INSERT INTO languages (name, color) VALUES (?, ?)',
        (data['name'], data['color'])
    )
    
    language_id = cursor.lastrowid
    conn.commit()
    conn.close()
    
    return jsonify({"id": language_id, "name": data['name'], "color": data['color']}), 201

@app.route('/api/languages/<int:language_id>', methods=['PUT'])
def update_language(language_id):
    data = request.json
    
    if not data:
        return jsonify({"error": "No data provided"}), 400
    
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    # Check if language exists
    cursor.execute('SELECT * FROM languages WHERE id = ?', (language_id,))
    if not cursor.fetchone():
        conn.close()
        return jsonify({"error": "Language not found"}), 404
    
    # Update fields that are provided
    updates = []
    values = []
    
    if 'name' in data:
        updates.append('name = ?')
        values.append(data['name'])
    
    if 'color' in data:
        updates.append('color = ?')
        values.append(data['color'])
    
    if updates:
        values.append(language_id)
        cursor.execute(
            f'UPDATE languages SET {", ".join(updates)} WHERE id = ?',
            values
        )
        conn.commit()
    
    # Get updated language
    cursor.execute('SELECT * FROM languages WHERE id = ?', (language_id,))
    conn.row_factory = sqlite3.Row
    language = dict(cursor.fetchone())
    
    conn.close()
    return jsonify(language)

@app.route('/api/languages/<int:language_id>', methods=['DELETE'])
def delete_language(language_id):
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    # Check if language exists
    cursor.execute('SELECT * FROM languages WHERE id = ?', (language_id,))
    if not cursor.fetchone():
        conn.close()
        return jsonify({"error": "Language not found"}), 404
    
    cursor.execute('DELETE FROM languages WHERE id = ?', (language_id,))
    conn.commit()
    conn.close()
    
    return jsonify({"message": "Language deleted successfully"})

if __name__ == '__main__':
    app.run(debug=True, port=5000)
