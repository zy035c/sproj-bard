"""
This is the main file for the lamport clock project

example: python main.py -p 9090 9091 9092
"""

# set current path as the root path
import os
import sys
sys.path.append(os.path.dirname(os.path.realpath(__file__)))

import argparse
import subprocess
import platform
from myhttp import send_insert, send_get, concur_send

port_list = []
pid_list = []

def terminal_interact():
    """
    # this python program will remain running
    # on each stop, the user will be instructed to type numbers to execute specific commands
    # the user can see a list of ports and their pid
    # 1. type the port number to send a message to
    # 2. type the port number to receive a message from
    # 3. concurrently send a message to multiple ports

    # the user can type "exit" to stop the program
    """
    while True:
        print("-----------------------------------")
        print("port_list: ", port_list)
        print("pid_list: ", pid_list)
        print("[1] choose a port number to send a message to")
        print("[2] choose a port number to read most recent message from")
        print("[3] concurrently send a message to multiple ports")
        print("[4] get messages from all ports")
        print("[exit] stop the program")
        print("-----------------------------------")

        command = input("Enter command: ")
        if command == "exit":
            for pid in pid_list:
                subprocess.call(["kill", str(pid)])
            break
        elif command == "1":
            port = input("Enter port number: ")
            if port in port_list:
                message = input("Enter message: ")
                pid = pid_list[port_list.index(port)]
                send_insert(port, message, pid)
            else:
                print("port not found")
        elif command == "2":
            port = input("Enter port number: ")
            if port in port_list:
                pid = pid_list[port_list.index(port)]
                send_get(port, pid)
            else:
                print("port not found")
        elif command == "3":
            # loop: type port number and message and store in a dict
            # until "done" is typed
            # then call concur_send
            msgs = {}
            while True:
                port = input("([done] to finish input) Enter port number: ")
                if port == "done":
                    break
                if port not in port_list:
                    print("port not found")
                    continue
                message = input("Enter message: ")
                pid = pid_list[port_list.index(port)]
                msgs[port] = (message, pid)

            if len(msgs) > 0:
                concur_send(
                    list(msgs.keys()),
                    [msg[0] for msg in msgs.values()],
                    [msg[1] for msg in msgs.values()]
                )

            print(f"Successfully sent {len(msgs)} messages")
        elif command == "4":
            # get from all ports
            for port in port_list:
                pid = pid_list[port_list.index(port)]
                send_get(port, pid)
        else:
            print("Invalid command")

        print("\n")
    pass

if __name__ == "__main__":
    # parse input options: -p <a list of port>
    # example: python main.py -p 9090 9091 9092
    parser = argparse.ArgumentParser(description='Lamport Clock')
    parser.add_argument('-p', '--port', nargs='+', help='port number', required=True)

    args = parser.parse_args()
    port_list = args.port

    os_type = platform.system()

    print("-> Operating System:", os_type)

    # will start n golang server processes with the given port numbers
    # the command to start a server is go run server.go -port=9090 -ps=8000,9000,10000
    # where 9090 is the port number of the server, and 8000,9000,10000 are the port numbers of the other servers
    # code start below

    working_dir = os.path.dirname(os.path.abspath(__file__))
    for port in port_list:
        cmd = ["go", "run", "server.go", "-port="+port, "-ps="+",".join(filter(lambda x: x != port, port_list))]
        print(" ".join(cmd))

        # will open a new terminal for each server
        if os_type == "Linux":
            popen = subprocess.Popen(["gnome-terminal", "--", "bash", "-c", " ".join(cmd)])
        elif os_type == "Darwin":
            apple_script = f'tell application "Terminal" to do script "cd {working_dir};{" ".join(cmd)}"'
            # osascript AppleScript
            popen = subprocess.Popen(['osascript', '-e', apple_script])

            # popen = subprocess.Popen(["open", "-a", "Terminal", "bash", "-c", " ".join(cmd)])
        else:
            print("Unsupported operating system:", os_type)
            exit()

        pid_list.append(popen.pid)

    terminal_interact()
