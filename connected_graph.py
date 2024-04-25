import random

def generate_graph(n, m):

    graph = [[] for _ in range(n)]

    for i in range(n):
        num_neighbors = min(random.randint(1, m), n - 1)

        neighbors = random.sample(range(n), num_neighbors)

        for neighbor in neighbors:
            if neighbor != i:
                graph[i].append(neighbor)
                graph[neighbor].append(i)

    return graph
