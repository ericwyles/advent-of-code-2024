package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type GateType int

const (
	INPUT GateType = iota
	AND
	OR
	XOR
)

func (g GateType) String() string {
	switch g {
	case INPUT:
		return "INPUT"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	default:
		return "UNKNOWN"
	}
}

func (n *LogicGateNode) OutputValString() string {
	if n.OutputVal {
		return "1"
	}
	return "0"
}

type LogicGateNode struct {
	IDVal     int64
	Name      string
	GateKind  GateType
	OutputVal bool
}

// ID satisfies Gonum's graph.Node interface.
func (n LogicGateNode) ID() int64 {
	return n.IDVal
}

func (n *LogicGateNode) Evaluate(in1, in2 bool) bool {
	switch n.GateKind {
	case AND:
		return in1 && in2
	case OR:
		return in1 || in2
	case XOR:
		return in1 != in2
	}
	return false
}

var swaps map[string]string

func main() {
	swaps = make(map[string]string)
	args := os.Args

	populateSwaps(swaps, args[1:])
	swaps["khg"] = "tvb"
	swaps["tvb"] = "khg"
	swaps["z12"] = "vdc"
	swaps["vdc"] = "z12"
	swaps["nhn"] = "z21"
	swaps["z21"] = "nhn"
	swaps["gst"] = "z33"
	swaps["z33"] = "gst"

	outputNames := getOutputNames()
	fmt.Printf("All outputs: %v\n", outputNames)

	fmt.Printf("Swaps to make: %v\n", swaps)
	g, nodeMap := readInput()

	start1 := time.Now()
	fmt.Printf("Graph has %d nodes and %d edges.\n", g.Nodes().Len(), g.Edges().Len())
	elapsed1 := time.Since(start1)

	processCircuit(g)

	xBinary, xDecimal := getBinaryAndDecimalValues("x", nodeMap)
	yBinary, yDecimal := getBinaryAndDecimalValues("y", nodeMap)
	zBinary, zDecimal := getBinaryAndDecimalValues("z", nodeMap)
	fmt.Printf("x binary=[%s] x decimal=[%d]\n", xBinary, xDecimal)
	fmt.Printf("y binary=[%s] y decimal=[%d]\n", yBinary, yDecimal)
	fmt.Printf("z binary=[%s] z decimal=[%d]\n", zBinary, zDecimal)
	fmt.Println("Results:")
	target := xDecimal + yDecimal
	actual := zDecimal
	fmt.Printf("Target: %d\n", target)
	fmt.Printf("Actual: %d\n", actual)
	fmt.Printf(" Delta: %d\n", zDecimal-(xDecimal+yDecimal))

	// if target == actual {
	// 	panic("found it")
	// }

	start2 := time.Now()
	ExportGraphToStyledGraphviz(g, nodeMap)
	elapsed2 := time.Since(start2)

	fmt.Printf("Calculation time: Part 1 [%d]ms, Part 2 [%d] ms\n", elapsed1.Milliseconds(), elapsed2.Milliseconds())
}

func populateSwaps(swaps map[string]string, s []string) {
	for _, arg := range s {
		swap := strings.Split(arg, ",")
		swaps[swap[0]] = swap[1]
		swaps[swap[1]] = swap[0]
	}
}

func ExportGraphToStyledGraphviz(g *simple.DirectedGraph, nodeMap map[string]*LogicGateNode) error {
	var builder strings.Builder

	// Start the DOT graph definition
	builder.WriteString("digraph G {\n")

	// Helper to write a subgraph
	writeSubgraph := func(subgraphName, color string, nodes []string) {
		builder.WriteString(fmt.Sprintf("  subgraph %s {\n", subgraphName))
		builder.WriteString(fmt.Sprintf("    node [style=filled,color=%s];\n", color))
		for _, node := range nodes {
			builder.WriteString(fmt.Sprintf("    %s;\n", node))
		}
		builder.WriteString("  }\n")
	}

	// Categorize nodes
	var inputX, inputY, gatesAnd, gatesOr, gatesXor, outputs []string

	for _, node := range nodeMap {
		switch {
		case strings.HasPrefix(node.Name, "x"):
			inputX = append(inputX, node.Name)
		case strings.HasPrefix(node.Name, "y"):
			inputY = append(inputY, node.Name)
		case strings.HasPrefix(node.Name, "z"):
			outputs = append(outputs, node.Name)
		case node.GateKind == AND:
			gatesAnd = append(gatesAnd, node.Name)
		case node.GateKind == OR:
			gatesOr = append(gatesOr, node.Name)
		case node.GateKind == XOR:
			gatesXor = append(gatesXor, node.Name)
		}
	}

	// Write subgraphs for each category
	writeSubgraph("input_x", "lightgrey", inputX)
	writeSubgraph("input_y", "lightgrey", inputY)
	writeSubgraph("gates_and", "lightgreen", gatesAnd)
	writeSubgraph("gates_or", "yellow", gatesOr)
	writeSubgraph("gates_xor", "lightskyblue", gatesXor)
	writeSubgraph("output_z", "lightgrey", outputs)

	// Write edges
	edges := g.Edges()
	for edges.Next() {
		edge := edges.Edge()
		from := edge.From().(*LogicGateNode)
		to := edge.To().(*LogicGateNode)
		builder.WriteString(fmt.Sprintf("  %s -> %s;\n", from.Name, to.Name))
	}

	// Close the DOT graph definition
	builder.WriteString("}\n")

	// Write the graph to a file
	return writeGraphToFile("graph.dot", builder.String())
}

// Helper function to write content to a file
func writeGraphToFile(filename, content string) error {
	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	return err
}

func getBinaryAndDecimalValues(prefix string, nodeMap map[string]*LogicGateNode) (string, int) {
	var zNames []string
	for name := range nodeMap {
		if strings.HasPrefix(name, prefix) {
			zNames = append(zNames, name)
		}
	}

	sort.Slice(zNames, func(i, j int) bool {
		return zNames[i] > zNames[j]
	})

	zString := ""
	for _, name := range zNames {
		node := nodeMap[name]
		zString += node.OutputValString()
	}

	zInt, _ := strconv.ParseInt(zString, 2, 64)

	return zString, int(zInt)

}

func processCircuit(g *simple.DirectedGraph) {
	// Get topologically sorted nodes
	sorted, _ := topo.Sort(g)

	for _, n := range sorted {
		// Safely assert the type
		gate, ok := n.(*LogicGateNode)
		if !ok {
			continue
		}

		// Skip input gates
		if gate.GateKind == INPUT {
			continue
		}

		// Get the input nodes
		inputs := graph.NodesOf(g.To(gate.ID()))
		if len(inputs) != 2 {
			panic(fmt.Sprintf("wrong number of inputs for gate %s: expected 2, got %d", gate.Name, len(inputs)))
		}

		// Check for nil inputs
		if inputs[0] == nil || inputs[1] == nil {
			panic(fmt.Sprintf("nil input detected for gate %s", gate.Name))
		}

		// Safely assert the type of input nodes
		input1, ok1 := inputs[0].(*LogicGateNode)
		input2, ok2 := inputs[1].(*LogicGateNode)
		if !ok1 || !ok2 {
			panic(fmt.Sprintf("input nodes for gate %s are not of type *LogicGateNode", gate.Name))
		}

		// Execute the logic operation
		gate.OutputVal = executeBooleanLogic(gate.GateKind, input1, input2)
	}
}

func executeBooleanLogic(gateType GateType, node1, node2 *LogicGateNode) bool {
	if gateType == AND {
		return node1.OutputVal && node2.OutputVal
	} else if gateType == OR {
		return node1.OutputVal || node2.OutputVal
	} else if gateType == XOR {
		return node1.OutputVal != node2.OutputVal
	} else {
		panic(fmt.Sprintf("Unknown gate type %v\n", gateType))
	}
}

func getOutputNames() []string {
	var outputNames []string

	file, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read each line from the file
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "->") {
			arrowParts := strings.Split(line, "->")

			outputNames = append(outputNames, strings.TrimSpace(arrowParts[1]))
		}
	}

	return outputNames
}

func readInput() (*simple.DirectedGraph, map[string]*LogicGateNode) {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(fmt.Sprintf("Failed to open file: %v", err))
	}
	defer file.Close()

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)

	g := simple.NewDirectedGraph()
	nodeMap := make(map[string]*LogicGateNode) // name -> node
	var nextID int64

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, ":") {
			err := parseInputNode(line, g, nodeMap, &nextID)
			if err != nil {
				panic(err)
			}
		}

		if strings.Contains(line, "->") {
			err := parseGateDefinitions(line, g, nodeMap, &nextID)
			if err != nil && err != io.EOF {
				panic(err)
			}
		}
	}

	return g, nodeMap
}

func parseInputNode(line string, g *simple.DirectedGraph, nodeMap map[string]*LogicGateNode, nextID *int64) error {
	// Each line looks like "x00: 1"
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid input node line: %s", line)
	}
	name := strings.TrimSpace(parts[0])
	valStr := strings.TrimSpace(parts[1])

	boolVal, err := parseBitToBool(valStr)
	if err != nil {
		return fmt.Errorf("error parsing bit for %s: %v", name, err)
	}

	// Create the node if not present
	node, exists := nodeMap[name]
	if !exists {
		node = &LogicGateNode{
			IDVal:     *nextID,
			Name:      name,
			GateKind:  INPUT,   // input node
			OutputVal: boolVal, // initial value
		}
		*nextID++
		nodeMap[name] = node
		g.AddNode(node)
	} else {
		// If the node already exists, override its GateKind & OutputVal
		node.GateKind = INPUT
		node.OutputVal = boolVal
	}
	return nil
}

func parseGateDefinitions(line string, g *simple.DirectedGraph, nodeMap map[string]*LogicGateNode, nextID *int64) error {
	// Example line: "x00 AND y00 -> z00"
	// We can split on " -> " first.
	arrowParts := strings.Split(line, "->")
	if len(arrowParts) != 2 {
		return fmt.Errorf("invalid gate definition line (missing '->'): %s", line)
	}

	lhs := strings.TrimSpace(arrowParts[0]) // "x00 AND y00"
	rhs := strings.TrimSpace(arrowParts[1]) // "z00"

	//fmt.Printf("rhs=%s\n", rhs)
	nameToSwap, exists := swaps[rhs]
	if exists {
		fmt.Printf("Replacing right hand side node %s with %s\n", rhs, nameToSwap)
		rhs = nameToSwap
	}

	tokens := strings.Fields(lhs) // ["x00", "AND", "y00"] or ["x02", "OR", "y02"]
	if len(tokens) != 3 {
		return fmt.Errorf("invalid LHS format: %s", lhs)
	}
	leftNodeName := tokens[0]
	opStr := tokens[1] // AND, OR, XOR
	rightNodeName := tokens[2]

	// Parse operator
	gateKind, err := parseGateType(opStr)
	if err != nil {
		return fmt.Errorf("unknown gate type %q: %v", opStr, err)
	}

	newGateName := rhs

	leftNode := ensureNodeExists(leftNodeName, nodeMap, g, nextID)
	rightNode := ensureNodeExists(rightNodeName, nodeMap, g, nextID)

	newGate, exists := nodeMap[newGateName]
	if !exists {
		newGate = &LogicGateNode{
			IDVal:     *nextID,
			Name:      newGateName,
			GateKind:  gateKind,
			OutputVal: false, // default, will be computed later
		}
		*nextID++
		nodeMap[newGateName] = newGate
		g.AddNode(newGate)
	} else {
		newGate.GateKind = gateKind
	}

	g.SetEdge(g.NewEdge(leftNode, newGate))
	g.SetEdge(g.NewEdge(rightNode, newGate))

	return nil
}

func ensureNodeExists(name string, nodeMap map[string]*LogicGateNode, g *simple.DirectedGraph, nextID *int64) *LogicGateNode {
	node, ok := nodeMap[name]
	if !ok {
		node = &LogicGateNode{
			IDVal:     *nextID,
			Name:      name,
			GateKind:  INPUT, // default to input if we don't have info yet
			OutputVal: false, // can be overwritten later
		}
		*nextID++
		nodeMap[name] = node
		g.AddNode(node)
	}
	return node
}

func parseBitToBool(bitStr string) (bool, error) {
	bitStr = strings.TrimSpace(bitStr)
	switch bitStr {
	case "1":
		return true, nil
	case "0":
		return false, nil
	}
	return false, fmt.Errorf("expected '0' or '1', got %q", bitStr)
}

func parseGateType(op string) (GateType, error) {
	switch strings.ToUpper(op) {
	case "AND":
		return AND, nil
	case "OR":
		return OR, nil
	case "XOR":
		return XOR, nil
	}
	return -1, fmt.Errorf("unrecognized operator %s", op)
}

// TraceGateDependencies identifies the direct dependencies of each gate
func TraceGateDependencies(g *simple.DirectedGraph, nodeMap map[string]*LogicGateNode) map[string][]string {
	dependencies := make(map[string][]string)

	for _, node := range nodeMap {
		inputs := graph.NodesOf(g.To(node.ID()))
		var inputNames []string
		for _, input := range inputs {
			inputNames = append(inputNames, input.(*LogicGateNode).Name)
		}
		dependencies[node.Name] = inputNames
	}

	return dependencies
}

// ComputeExpectedOutputs computes the expected output of all gates
func ComputeExpectedOutputs(g *simple.DirectedGraph, nodeMap map[string]*LogicGateNode, x, y string, dependencies map[string][]string) map[string]bool {
	expectedOutputs := make(map[string]bool)

	// Set expected outputs for input nodes
	for i, bit := range strings.Split(reverse(x), "") {
		nodeName := fmt.Sprintf("x%02d", i)
		expectedOutputs[nodeName] = bit == "1"
	}

	for i, bit := range strings.Split(reverse(y), "") {
		nodeName := fmt.Sprintf("y%02d", i)
		expectedOutputs[nodeName] = bit == "1"
	}

	// Traverse gates in topological order
	for _, node := range nodeMap {
		if node.GateKind == INPUT {
			continue // Skip input nodes, already handled
		}

		// Get inputs to this gate
		inputNames := dependencies[node.Name]
		if len(inputNames) < 2 {
			continue // Invalid gate
		}

		in1 := expectedOutputs[inputNames[0]]
		in2 := expectedOutputs[inputNames[1]]

		// Compute the expected output based on gate kind
		switch node.GateKind {
		case AND:
			expectedOutputs[node.Name] = in1 && in2
		case OR:
			expectedOutputs[node.Name] = in1 || in2
		case XOR:
			expectedOutputs[node.Name] = in1 != in2
		}
	}

	return expectedOutputs
}

// EvaluateCircuit traverses the graph and evaluates actual outputs
func EvaluateCircuit(g *simple.DirectedGraph, nodeMap map[string]*LogicGateNode) map[string]bool {
	outputs := make(map[string]bool)

	// Traverse gates in topological order (or manually)
	for _, node := range nodeMap {
		if node.GateKind == INPUT {
			outputs[node.Name] = node.OutputVal // Use preset value
		} else {
			// Get inputs
			inputs := graph.NodesOf(g.To(node.ID()))
			if len(inputs) < 2 {
				continue // Skip invalid gates
			}
			in1 := inputs[0].(*LogicGateNode).OutputVal
			in2 := inputs[1].(*LogicGateNode).OutputVal

			// Compute output
			switch node.GateKind {
			case AND:
				node.OutputVal = in1 && in2
			case OR:
				node.OutputVal = in1 || in2
			case XOR:
				node.OutputVal = in1 != in2
			}

			outputs[node.Name] = node.OutputVal
		}
	}

	return outputs
}

// CompareAllGates compares expected vs. actual outputs for all gates
func CompareAllGates(expected, actual map[string]bool, nodeMap map[string]*LogicGateNode) []string {
	mismatches := []string{}

	for name, expectedVal := range expected {
		if actual[name] != expectedVal {
			mismatches = append(mismatches, name)
		}
	}

	return mismatches
}

// reverse reverses a binary string (used for least significant bit first)
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
