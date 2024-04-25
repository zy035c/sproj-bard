from connected_graph import generate_graph
import docker


def create_container(image_name, id, env_vars={}, client=docker.from_env()):
    container_name = f"container-{id + 1}".lower()
    print(f"Creating container: {container_name}")
    client.containers.run(image_name, detach=True, name=container_name, environment=env_vars)


def create_cluster(image_name, numNode=20, maxNeighbor=5):
    graph = generate_graph(numNode, maxNeighbor)

    base_port = 9000

    client = docker.from_env()

    for node in range(numNode):
        port_list = [str(base_port + adj) for adj in graph[node]]
        port_my = str(base_port + node)

        create_container(
            image_name=image_name,
            id=node,
            env_vars={
                'PORT': port_my,
                'GPLIST': port_list
            },
            client=client
        )


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description='Create multiple containers from the same Docker image.')
    parser.add_argument('image_name', type=str, help='Name of the Docker image')
    parser.add_argument('num_containers', type=int, help='Number of containers to create')
    args = parser.parse_args()

    create_cluster(image_name=args.image_name, numNode=args.num_containers)
