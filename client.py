from flask import Flask, jsonify, render_template
from flask_cors import CORS
import requests

app = Flask(__name__)
CORS(app)

@app.route('/')
def index():
    return render_template('websockets.html')

@app.route('/get_history')
def get_history():
    url = "http://localhost:8080/history"
    response = requests.get(url)
    if response.status_code == 200:
        return jsonify(response.json())
    else:
        return jsonify({"error": "Failed to fetch chat history"}), 500

if __name__ == '__main__':
    app.run(port=8787,debug=True)