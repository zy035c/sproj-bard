import requests
import aiohttp
import asyncio


def send_insert(port, text, proc_id):
    """
        send POST request to http://localhost:port/insert-data
        with following json:
        {
            "req_type": "insert",
            "order": {
                "product_name": "Hello World 2",
                "timestamp": 1,
                "proc_id": 1286
            }
        }

        timestamp be arbitrary
    """
    url = "http://localhost:" + str(port) + "/insert-data"
    data = {
        "req_type": "insert",
        "order": {
            "product_name": text,
            "timestamp": 1,
            "proc_id": proc_id
        }
    }
    print("-> Insert Post Json: ", data)
    resp = requests.post(url, json=data)
    print("-> Insert Response: ", resp.json())
    pass


def send_get(port, proc_id):
    """
        send GET request to http://localhost:port/get-data
        with no json
    """
    url = "http://localhost:" + str(port) + "/get-data"
    data = {
        "req_type": "get",
    }
    resp = requests.get(url, json=data)
    # parse and print the response
    print("-> Get Response: ", resp.json())
    pass

def concur_send(ports, message_list, proc_ids):
    """
        use asyncio and aiohttp to concurrently send messages to multiple ports
    """

    tasks = []
    # create a list of tasks
    for msg, port, proc_id in zip(message_list, ports, proc_ids):

        async def post_task(proc_id_, port_, msg_):

            url = "http://localhost:" + str(port_) + "/insert-data"
            data = {
                "req_type": "insert",
                "order": {
                    "product_name": msg_,
                    "timestamp": 1,  # arbitrary
                    "proc_id": 1  # arbitrary
                }
            }
            async with aiohttp.ClientSession() as session:
                async with session.post(url, json=data) as response:
                    return await response.text()

        tasks.append(post_task(proc_id, port, msg))
        print(f"-> Gathered a task for port {port}")

    async def main():
        await asyncio.gather(*tasks)
    asyncio.run(main())
