import os
from connected_graph import generate_graph
import docker


current_dir = os.getcwd()
tag = "golang-image"

def create_container(image_name, addr, id, env_vars={}, client=docker.from_env()):

    container_name = f"container-{id + 1}".lower()
    print(f"Creating container: {container_name}")

    print("Will set local addr to", addr)
    print("Will set env_vars to", env_vars)

    container = client.containers.run(
        tag, detach=True, name=container_name, environment=env_vars
    )

    container.exec_run(f"ifconfig eth0 {addr}")

    print("------------------success------------------")

def create_cluster(image_name, numNode=20, maxNeighbor=5):
    graph, plot_thread = generate_graph(numNode, maxNeighbor)

    plot_thread.start()
    plot_thread.join()

    baseAddr = "172.18.0.x:9000"

    print("Graph:")
    for i, elem in enumerate(graph):
        print(f'[{i}] {elem}')

    client = docker.from_env()
    # network = client.networks.create("go_net")

    dockerfile_path = os.path.join(current_dir, image_name)
    image, build_logs = client.images.build(
        path=current_dir,
        dockerfile=dockerfile_path,
        tag=tag,
    )

    for node in range(numNode):
        addr_list = [getLocalAddr(baseAddr, adj) for adj in graph[node]]
        addr_local = getLocalAddr(baseAddr, node)

        create_container(
            image_name=image_name,
            addr=addr_local[:-5],
            id=node,
            env_vars={
                'G_ADDR': addr_local,
                'G_REMOTE_ADDR_LIST': addr_list,
            },
            client=client,
        )


def getLocalAddr(base_addr, node) -> str:
    replaced = base_addr.replace("x", str(node))
    return replaced

if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description='Create multiple containers from the same Docker image.')
    parser.add_argument('image_name', type=str, help='Name of the Docker image')
    parser.add_argument('num_containers', type=int, help='Number of containers to create')
    parser.add_argument(
        'num_max_neighbors',
        type=int,
        help='Maximum number of neighbors a node can have',
    )
    args = parser.parse_args()

    create_cluster(image_name=args.image_name, numNode=args.num_containers, maxNeighbor=args.num_max_neighbors)
