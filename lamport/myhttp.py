import requests


def send_insert(port, text, proc_id):
    """
        send POST request to http://localhost:port/insert-data
        with following json:
        {
            "req-type": "insert",
            "proc-id": 1286,
            "order": {
                "product_name": "Hello World 2",
                "timestamp": 1
            }
        }

        timestamp be arbitrary
    """
    url = "http://localhost:" + str(port) + "/insert-data"
    data = {
        "req-type": "insert",
        "proc-id": proc_id,
        "order": {
            "product_name": text,
            "timestamp": 1
        }
    }
    resp = requests.post(url, json=data)
    pass


def send_get(port, proc_id):
    """
        send GET request to http://localhost:port/get-data
        with no json
    """
    url = "http://localhost:" + str(port) + "/get-data"
    data = {
        "req-type": "get",
        "proc-id": proc_id
    }
    resp = requests.get(url, json=data)
    # parse and print the response
    print(resp.text)
    pass