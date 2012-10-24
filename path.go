// Cobbled together from https://github.com/sperre/astar.

package rog

import (
	"container/heap"
)

type Walkable interface {
	GetRoughness() int
}

type WalkableMap interface {
	Width() int
	Height() int
	Roughness(x, y int) int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Node

/* sort.Interface */
func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use smaller than here.
	return pq[i].f < pq[j].f
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].heap_index = i
	pq[j].heap_index = j
}

/* heap.interface */
func (pq *PriorityQueue) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	// To simplify indexing expressions in these methods, we save a copy of the
	// slice object. We could instead write (*pq)[i].
	a := *pq
	n := len(a)
	a = a[0 : n+1]
	item := x.(*Node)
	item.heap_index = n
	a[n] = item
	*pq = a
}

func (pq *PriorityQueue) Pop() interface{} {
	a := *pq
	n := len(a)
	item := a[n-1]
	item.heap_index = -1 // for safety
	*pq = a[0 : n-1]
	return item
}

/* Node */

func (pq *PriorityQueue) PushNode(n *Node) {
	heap.Push(pq, n)
}

func (pq *PriorityQueue) PopNode() *Node {
	return heap.Pop(pq).(*Node)
}

func (pq *PriorityQueue) RemoveNode(n *Node) {
	heap.Remove(pq, n.heap_index)
}

/*** Helper functions ***/

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

/*** MapData type and related consts ***/

// Tile information
const (
	PATH_MIN = 1 << iota
	PATH_MAX
)

// Tile movement costs
const (
	COST_STRAIGHT = 1000
	COST_DIAGONAL = 1414
)

// TODO: move to serialization module

// func str_map(data MapData, nodes []*Node) string {
// 	var result string
// 	for i, row := range data {
// 		for j, cell := range row {
// 			added := false
// 			for _, node := range nodes {
// 				if node.X == i && node.Y == j {
// 					result += "o"
// 					added = true
// 					break
// 				}
// 			}
// 			if !added {
// 				switch cell {
// 				case PATH_MIN:
// 					result += "."
// 				case PATH_MAX:
// 					result += "#"
// 				default: //Unknown
// 					result += "?"
// 				}
// 			}
// 		}
// 		result += "\n"
// 	}
// 	return result
// }

/*** Node type ***/

// X and Y are coordinates, parent is a link to where we came from. cost are the
// estimated cost from start along the best known path. h is the heuristic value
// (air line distance to goal).
type Node struct {
	X, Y       int
	parent     *Node
	f, g, h    int
	heap_index int // only used and maintained by pqueue
}

// Create a new Node
func NewNode(x, y int) *Node {
	node := &Node{
		X:      x,
		Y:      y,
		parent: nil,
		f:      0, // f = g + h
		g:      0, // Cost from node to start
		h:      0, // Estimated cost from node to finish
	}
	return node
}

// Return string representation of the node
func (n *Node) String() string {
	return ""
	//	return fmt.Sprintf("<Node x:%d y:%d addr:%d>", n.X, n.Y, &n)
}

/*** nodeList type ***/

type nodeList struct {
	nodes      map[int]*Node
	rows, cols int
}

func newNodeList(rows, cols int) *nodeList {
	return &nodeList{
		nodes: make(map[int]*Node, rows*cols),
		rows:  rows,
		cols:  cols,
	}
}

func (n *nodeList) addNode(node *Node) {
	n.nodes[node.X+node.Y*n.rows] = node
}

func (n *nodeList) getNode(x, y int) *Node {
	return n.nodes[x+y*n.rows]
}

func (n *nodeList) removeNode(node *Node) {
	delete(n.nodes, node.X+node.Y*n.rows)
}

func (n *nodeList) hasNode(node *Node) bool {
	if n.getNode(node.X, node.Y) != nil {
		return true
	}
	return false
}

/*** Graph type ***/

// Start, stop nodes and a slice of nodes
type Graph struct {
	nodes *nodeList // Used to avoid duplicated nodes!
	wmap WalkableMap
}

// Return a Graph from a map of coordinates (those that are passible)
func NewGraph(wmap WalkableMap) *Graph {
	//var start, stop *Node
	return &Graph{
		nodes: newNodeList(wmap.Width(), wmap.Height()),
		wmap:  wmap,
	}
}

// Get or create a *Node based on x, y coordinates. Avoids duplicated nodes!
func (g *Graph) Node(x, y int) *Node {
	//Check if node is already in the graph
	var node *Node
	node = g.nodes.getNode(x, y)

	if node == nil && (g.wmap.Roughness(x, y) < PATH_MAX) {
		//Create a new node and add it to the graph
		node = NewNode(x, y)
		g.nodes.addNode(node)
	}
	return node
}

/* Astar func */

func retracePath(current_node *Node) []*Node {
	var path []*Node
	path = append(path, current_node)
	for current_node.parent != nil {
		path = append(path, current_node.parent)
		current_node = current_node.parent
	}
	//Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// Diagonal/Chebyshev distance is used.
func Heuristic(tile, stop *Node) (h int) {
	h_diag := min(abs(tile.X-stop.X), abs(tile.Y-stop.Y))
	h_stra := abs(tile.X-stop.X) + abs(tile.Y-stop.Y)
	h = COST_DIAGONAL*h_diag + COST_STRAIGHT*(h_stra-2*h_diag)

	/* TODO: Breaking ties:
	dx1 := tile.X - stop.X
	dy1 := tile.Y - stop.Y
	dx2 := start.X - stop.X
	dy2 := start.Y - stop.Y
	cross := abs(dx1*dy2 - dx2*dy1)
	h += cross * COST_DIAGONAL/100
	*/
	return
}

// 8 directions adjecentDirs and costs
var adjecentDirs8 = [][3]int{
	{-1, -1, COST_DIAGONAL}, {-1, 0, COST_STRAIGHT}, {-1, 1, COST_DIAGONAL},
	{0, -1, COST_STRAIGHT}, {0, 1, COST_STRAIGHT},
	{1, -1, COST_DIAGONAL}, {1, 0, COST_STRAIGHT}, {1, 1, COST_DIAGONAL},
}

// 4 directions adjecentDirs and costs
var adjecentDirs4 = [][3]int{
	{-1, 0, COST_STRAIGHT},
	{0, -1, COST_STRAIGHT}, {0, 1, COST_STRAIGHT},
	{1, 0, COST_STRAIGHT},
}

// A* search algorithm. See http://en.wikipedia.org/wiki/A*_search_algorithm
func Astar(wmap WalkableMap, startx, starty, stopx, stopy int, dir8 bool) []*Node {
	graph := NewGraph(wmap)
	rows, cols := wmap.Width(), wmap.Height()

	// Create lists
	closedSet := newNodeList(rows, cols)
	openSet := newNodeList(rows, cols)
	pq := make(PriorityQueue, 0, rows*cols) // heap, used to find minF

	// Move in 8 or 4 directions?
	var adjecentDirs [][3]int
	if dir8 {
		adjecentDirs = adjecentDirs8
	} else {
		adjecentDirs = adjecentDirs4
	}

	// TODO: GUARD: startx... stopy inside array range?

	// Add start node to the task list
	start := NewNode(startx, starty)
	stop := NewNode(stopx, stopy)
	openSet.addNode(start)
	pq.PushNode(start)

	for len(openSet.nodes) != 0 {
		// Get the node with the min H
		//current := openSet.minF()
		current := pq.PopNode()
		openSet.removeNode(current)
		closedSet.addNode(current)

		if current.X == stop.X && current.Y == stop.Y {
			// Finished, return shortest path
			//fmt.Println(str_map(map_data,  retracePath(current)))
			return retracePath(current)
		}

		for _, adir := range adjecentDirs {
			x, y := (current.X + adir[0]), (current.Y + adir[1])

			// Check if x, y is inside the map:
			if (x < 0) || (x >= rows) || (y < 0) || (y >= cols) {
				continue
			}

			neighbor := graph.Node(x, y)
			if neighbor == nil || closedSet.hasNode(neighbor) {
				// Wall, or old node
				continue
			}

			g_score := current.g + adir[2]

			if !openSet.hasNode(neighbor) {
				// Add new interesting node
				neighbor.parent = current
				neighbor.g = g_score
				neighbor.f = neighbor.g + Heuristic(neighbor, stop)
				openSet.addNode(neighbor)
				pq.PushNode(neighbor)
			} else if g_score < neighbor.g {
				// Update, old node
				pq.RemoveNode(neighbor)
				neighbor.parent = current
				neighbor.g = g_score
				neighbor.f = neighbor.g + Heuristic(neighbor, stop)
				pq.PushNode(neighbor)
			}

		}
	}

	return nil
}

// ----
