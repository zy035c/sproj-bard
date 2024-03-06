"""
This is the main file for the lamport clock

example: python main.py -p 9090 9091 9092
"""

# set current path as the root path
import os
import sys
sys.path.append(os.path.dirname(os.path.realpath(__file__)))

import argparse
import subprocess
from myhttp import send_insert, send_get

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
        print("port_list: ", port_list)
        print("pid_list: ", pid_list)
        print("1. type the port number to send a message to")
        print("2. type the port number to receive a message from")
        print("3. concurrently send a message to multiple ports")
        print("4. type 'exit' to stop the program")

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
        # elif command == "3":
        #     port = input("Enter port numbers (separated by space): ")
        #     port = port.split(" ")
        #     for p in port:
        #         if p not in port_list:
        #             print("port not found")
        #             break
        #     else:
        #         message = input("Enter message: ")
        #         for p in port:
        #             subprocess.call(["go", "run", "client.go", "-port="+p, "-m="+message])
        else:
            print("invalid command")
    pass

if __name__ == "__main__":
    # parse input options: -p <a list of port>
    # example: python main.py -p 9090 9091 9092
    parser = argparse.ArgumentParser(description='Lamport Clock')
    parser.add_argument('-p', '--port', nargs='+', help='port number', required=True)

    args = parser.parse_args()
    port_list = args.port

    print("port_list: ", port_list)
    
    # will start n golang server processes with the given port numbers
    # the command to start a server is go run server.go -port=9090 -ps=8000,9000,10000
    # where 9090 is the port number of the server, and 8000,9000,10000 are the port numbers of the other servers
    # code start below
    for port in port_list:
        cmd = ["go", "run", "server.go", "-port="+port, "-ps="+",".join(filter(lambda x: x != port, port_list))]
        print(" ".join(cmd))

        # will open a new terminal for each server
        popen = subprocess.Popen(["gnome-terminal", "--", "bash", "-c", " ".join(cmd)])
        pid_list.append(popen.pid)

    terminal_interact()
