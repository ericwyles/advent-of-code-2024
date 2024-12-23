package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type NamedNode struct {
	IDVal int64  // The numeric ID used by Gonum
	Name  string // The human-readable name, e.g. "Node1", "Server42", etc.
}

// ID satisfies Gonum's graph.Node interface.
func (n NamedNode) ID() int64 {
	return n.IDVal
}

func main() {
	g := readInput()

	start1 := time.Now()

	t := 0
	triangles := findTriangles(g)
	for _, tri := range triangles {
		a, b, c := tri[0].(*NamedNode).Name, tri[1].(*NamedNode).Name, tri[2].(*NamedNode).Name
		if a[0] == 't' || b[0] == 't' || c[0] == 't' {
			t++
		}
	}

	fmt.Printf("Graph has %d nodes and %d edges.\n", g.Nodes().Len(), g.Edges().Len())
	fmt.Printf("Found %d triangles.\n", len(triangles))
	fmt.Printf("Part 1 - Groups of 3 with T: %d\n", t)
	elapsed1 := time.Since(start1)

	start2 := time.Now()
	clique := findMaximumClique(g)
	names := make([]string, len(clique))
	for i, node := range clique {
		names[i] = node.(*NamedNode).Name
	}
	slices.Sort(names)
	password := names[0]
	for i := 1; i < len(names); i++ {
		password += fmt.Sprintf(",%s", names[i])
	}
	fmt.Printf("Part 2 - Password is: %s\n", password)
	elapsed2 := time.Since(start2)

	fmt.Printf("Calculation time: Part 1 [%d]ms, Part 2 [%d] ms\n", elapsed1.Milliseconds(), elapsed2.Milliseconds())
}

func readInput() *simple.UndirectedGraph {
	scanner := bufio.NewScanner(os.Stdin)

	// Create an unweighted undirected graph
	g := simple.NewUndirectedGraph()

	// We’ll store each string node in a map so we can handle
	// gonum's integer-based Node IDs.
	nodeMap := make(map[string]*NamedNode)
	var nextID int64

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			computers := strings.Split(line, "-")
			a, b := computers[0], computers[1]

			// Ensure we have a node in the graph for 'a'
			if _, exists := nodeMap[a]; !exists {
				nodeMap[a] = &NamedNode{
					IDVal: nextID,
					Name:  a,
				}
				g.AddNode(nodeMap[a])
				nextID++
			}

			// Ensure we have a node in the graph for 'b'
			if _, exists := nodeMap[b]; !exists {
				nodeMap[b] = &NamedNode{
					IDVal: nextID,
					Name:  b,
				}
				g.AddNode(nodeMap[b])
				nextID++
			}

			// Now create an edge between these two nodes
			aNode := nodeMap[a]
			bNode := nodeMap[b]
			// Since it’s undirected, just set the edge once
			g.SetEdge(g.NewEdge(aNode, bNode))
		}
	}

	return g
}

// findTriangles returns all cliques of size 3 (triangles) in g.
func findTriangles(g *simple.UndirectedGraph) [][]graph.Node {
	found := make(map[[3]int64]bool)

	var triangles [][]graph.Node

	// Convert the graph's node iterator to a slice
	allNodes := graph.NodesOf(g.Nodes())

	for _, n := range allNodes {
		// Get neighbors of n
		neighbors := graph.NodesOf(g.From(n.ID()))

		// Check each pair of neighbors
		for i := 0; i < len(neighbors); i++ {
			for j := i + 1; j < len(neighbors); j++ {
				u := neighbors[i]
				v := neighbors[j]

				// If there's an edge between these two neighbors, we have a triangle
				if g.HasEdgeBetween(u.ID(), v.ID()) {
					// Create a sorted ID slice so we identify the same triangle consistently
					trio := []int64{n.ID(), u.ID(), v.ID()}
					sort.Slice(trio, func(a, b int) bool { return trio[a] < trio[b] })
					var triKey [3]int64
					copy(triKey[:], trio)

					// Only add if we haven’t already seen this exact triangle
					if !found[triKey] {
						found[triKey] = true

						// Convert ID back to graph.Node using g.Node()
						triNodes := []graph.Node{
							g.Node(trio[0]),
							g.Node(trio[1]),
							g.Node(trio[2]),
						}
						triangles = append(triangles, triNodes)
					}
				}
			}
		}
	}
	return triangles
}

// findMaximumClique runs Bron–Kerbosch with pivoting to find all maximal cliques
// in g, then returns the single largest of those cliques.
//
// Reference: https://en.wikipedia.org/wiki/Bron–Kerbosch_algorithm
func findMaximumClique(g *simple.UndirectedGraph) []graph.Node {
	// Gather all nodes (P = the set of potential nodes)
	allNodes := graph.NodesOf(g.Nodes())
	// We’ll track them by ID in sets of int64 for speed.
	P := makeSet(allNodes)
	R := makeSet(nil) // empty set
	X := makeSet(nil) // empty set

	// We’ll store just the biggest clique we encounter.
	var maxClique []int64

	// The main recursive function with pivoting
	var bronKerboschPivot func(R, P, X map[int64]bool)
	bronKerboschPivot = func(R, P, X map[int64]bool) {
		// If P and X are both empty, R is a maximal clique
		if len(P) == 0 && len(X) == 0 {
			// If it's bigger than our current best, update
			if len(R) > len(maxClique) {
				maxClique = setToSlice(R)
			}
			return
		}

		// Choose a pivot (any node from P ∪ X)
		pivot := anyNode(union(P, X))

		// P \ N(pivot)
		// We'll only recurse on nodes in P that aren't neighbors of pivot
		pivotNeighbors := neighborsOf(g, pivot)
		toExplore := difference(P, pivotNeighbors)

		for n := range toExplore {
			// R ∪ {n}
			newR := union(R, singleton(n))
			// P ∩ N(n)
			nNeighbors := neighborsOf(g, n)
			newP := intersection(P, nNeighbors)
			// X ∩ N(n)
			newX := intersection(X, nNeighbors)

			bronKerboschPivot(newR, newP, newX)

			// Move n from P to X
			P = difference(P, singleton(n))
			X = union(X, singleton(n))

			if len(P) == 0 {
				return
			}
		}
	}

	// Run the algorithm
	bronKerboschPivot(R, P, X)

	// Convert IDs of our max clique back to []graph.Node
	// (We do this because we want the actual Node objects for printing, etc.)
	var largest []graph.Node
	for _, id := range maxClique {
		largest = append(largest, g.Node(id))
	}
	return largest
}

// neighborsOf returns a set of IDs for all nodes adjacent to the given node.
func neighborsOf(g *simple.UndirectedGraph, id int64) map[int64]bool {
	nbrs := make(map[int64]bool)
	for _, n := range graph.NodesOf(g.From(id)) {
		nbrs[n.ID()] = true
	}
	return nbrs
}

func makeSet(nodes []graph.Node) map[int64]bool {
	s := make(map[int64]bool)
	for _, n := range nodes {
		s[n.ID()] = true
	}
	return s
}

func singleton(id int64) map[int64]bool {
	return map[int64]bool{id: true}
}

func setToSlice(m map[int64]bool) []int64 {
	sl := make([]int64, 0, len(m))
	for id := range m {
		sl = append(sl, id)
	}
	// Sort for consistency
	sort.Slice(sl, func(i, j int) bool { return sl[i] < sl[j] })
	return sl
}

func union(a, b map[int64]bool) map[int64]bool {
	out := make(map[int64]bool, len(a)+len(b))
	for x := range a {
		out[x] = true
	}
	for x := range b {
		out[x] = true
	}
	return out
}

func intersection(a, b map[int64]bool) map[int64]bool {
	out := make(map[int64]bool)
	// iterate over smaller set for efficiency
	if len(a) < len(b) {
		for x := range a {
			if b[x] {
				out[x] = true
			}
		}
	} else {
		for x := range b {
			if a[x] {
				out[x] = true
			}
		}
	}
	return out
}

func difference(a, b map[int64]bool) map[int64]bool {
	out := make(map[int64]bool, len(a))
	for x := range a {
		if !b[x] {
			out[x] = true
		}
	}
	return out
}

func anyNode(m map[int64]bool) int64 {
	for x := range m {
		return x
	}
	return -1 // Should never happen if m isn't empty
}
