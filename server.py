from flask import Flask, request

app = Flask(__name__)

@app.route('/receive', methods=['POST'])
def receive_zip():
    file = request.files['file']
    file.save('received.zip')
    return 'Successfully received zip file.'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8000)
