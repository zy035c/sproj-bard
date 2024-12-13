package utils

type Graph struct {
	nodes map[int][]int // adjacency list
}

func NewGraph() *Graph {
	return &Graph{nodes: make(map[int][]int)}
}

func Append_if_not_exist(val int, slc []int) []int {
	for _, v := range slc {
		if v == val {
			return slc
		}
	}
	slc = append(slc, val)
	return slc
}

// Function to count independent cycles in a directed graph
func CountCycles(adj map[int][]int) int {
	visited := make(map[int]bool)
	stack := []int{}
	cycles := 0

	// Helper function for DFS
	var dfs func(node, start int)
	dfs = func(node, start int) {
		// Mark the current node as visited and add it to the stack
		visited[node] = true
		stack = append(stack, node)

		for _, neighbor := range adj[node] {
			if neighbor == start {
				// Found a cycle, count it and skip further exploration of this path
				cycles++
				continue
			} else if !visited[neighbor] {
				dfs(neighbor, start)
			}
		}

		// Backtrack: remove the node from the stack and mark it as unvisited
		stack = stack[:len(stack)-1]
		visited[node] = false
	}

	// Explore each node as a potential cycle start
	for node := range adj {
		for _, neighbor := range adj[node] {
			if neighbor > node {
				// Only start cycles from the smallest node to avoid duplicates
				dfs(node, node)
			}
		}
	}

	return cycles
}
