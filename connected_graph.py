import multiprocessing
import random
import threading

import networkx as nx
import matplotlib.pyplot as plt

def generate_graph(n, m):

    graph = [[] for _ in range(n)]

    father = [i for i in range(n)]
    edges = []

    def find(a):
        if father[a] != a:
            father[a] = find(father[a])
        return father[a]

    def union(a, b):
        a = find(a)
        b = find(b)

        father[b] = a

    for i in range(n):
        num_neighbors = random.randint(1, m)
        neighbors = random.sample(
            [x for x in range(n) if x != i],
            num_neighbors
        )

        for neighbor in neighbors:
            if len(graph[i]) == m:
                break
            if i not in graph[neighbor] and len(graph[neighbor]) < m:
                graph[i].append(neighbor)
                graph[neighbor].append(i)
                union(i, neighbor)
                edges.append((i, neighbor))

    components = set()

    for i in range(n):
        components.add(find(i))

    if len(components) > 1:
        print(f"The graph is not connected, retry. comp={len(components)}")
        return generate_graph(n, m)

    print("graph successfully generated.")

    def target():
        plot_graph_structure(edges=edges, n=n)
    graph_thread = multiprocessing.Process(target=target)

    return graph, graph_thread

def plot_graph_structure(edges, n):
    G = nx.Graph()
    
    G.add_nodes_from(range(n))
    G.add_edges_from(edges)

    pos = nx.spring_layout(G)
    nx.draw(
        G, pos, with_labels=True,
        node_color='skyblue',
        node_size=500,
        edge_color='black',
        linewidths=1,
        font_size=10
    )
    plt.title("Undirected Graph")
    plt.ion()
    plt.show()
    plt.pause(1)