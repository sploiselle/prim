package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// HeadEdge contains the HeadVertex and Distance from the Vertex.ID to which they belong
type Edge struct {
	HeadVertexID int
	HeadVertex   *Vertex
	TailVertexID int
	TailVertex   *Vertex
	Cost         float64
}

func (e Edge) String() string {
	return fmt.Sprintf("\nHeadVertex:\t%d\nTailVertex:\t%d\nDistance:\t%f\n", e.HeadVertexID, e.TailVertexID, e.Cost)
}

//A Vertex is a point on a graph, which contains []Edges to other Vertexes.
type Vertex struct {
	ID           int
	EdgesFrom    []*Edge
	EdgesTo      []*Edge
	CheapestEdge float64 // The cheapest edge that crosses the boundary
	Index        int     // Index of item in heap
}

func (v Vertex) String() string {
	return fmt.Sprintf("\nid:\t%d\nedges to: %v\nedges from: %v\nCheapestEdge:\t%g\nIndex:\t%d\n\n\n\n", v.ID, v.EdgesFrom, v.EdgesTo, v.CheapestEdge, v.Index)
}

// AddFromEdge appends an EdgesFromVertex to a Vertex's []Edges
func (v *Vertex) AddFromEdge(e *Edge) {
	v.EdgesFrom = append(v.EdgesFrom, e)
}

// AddFromEdge appends an EdgesFromVertex to a Vertex's []Edges
func (v *Vertex) AddToEdge(e *Edge) {
	v.EdgesTo = append(v.EdgesTo, e)
}

// VertexMap tracks Vertexes using [Vertex.ID]*Vertex
var VertexMap = make(map[int]*Vertex)

// A VertexHeap returns the Vertex from (V - X) with
// the lowest value of CheapestEdge
type VertexHeap []*Vertex

func (vh VertexHeap) Len() int { return len(vh) }

func (vh VertexHeap) Less(i, j int) bool {
	return vh[i].CheapestEdge < vh[j].CheapestEdge
}

func (vh VertexHeap) Swap(i, j int) {
	vh[i], vh[j] = vh[j], vh[i]
	vh[i].Index = i
	vh[j].Index = j
}

// Push adds Vertexes to VertexHeaps
func (vh *VertexHeap) Push(x interface{}) {
	n := len(*vh)
	v := x.(*Vertex)
	v.Index = n
	*vh = append(*vh, v)
}

// Pop returns the Vertex with the lowest value of CheapestEdge
func (vh *VertexHeap) Pop() interface{} {
	old := *vh
	n := len(old)
	v := old[n-1]
	v.Index = -1 // for safety, identify it's no longer in heap
	*vh = old[0 : n-1]
	return v
}

func (vh VertexHeap) Peek() {
	n := len(vh)

	if n == 0 {
		return
	}

	// fmt.Printf("\n\nNext item is\n%v with CheapestEdge %v", vh[0].ID, vh[0].CheapestEdge)
}

// Update modifies the DGS of a Vertex in the heap.
func (vh *VertexHeap) UpdateHeap(v *Vertex) {

	// fmt.Printf("\n\nUpdating vertices that point to %v\n", v.ID)

	for _, e := range v.EdgesTo {

		// fmt.Printf("\nBEFORE UPDATE:\ne.TailVertex.ID:\t\t%v\ne.Cost:\t\t\t\t%v\ne.TailVertex.CheapestEdge:\t%v", e.TailVertex.ID, e.Cost, e.TailVertex.CheapestEdge)

		if e.TailVertex.Index > -1 && (e.Cost < e.TailVertex.CheapestEdge) {
			e.TailVertex.CheapestEdge = e.Cost
			heap.Fix(vh, e.TailVertex.Index)

			// fmt.Printf("\nUPDATED:\ne.TailVertex.ID:\t\t%v\ne.Cost:\t\t\t\t%v\ne.TailVertex.CheapestEdge:\t%v", e.TailVertex.ID, e.Cost, e.TailVertex.CheapestEdge)
		} else {
			// fmt.Printf("\nNOT UPDATED:\ne.TailVertex.ID:\t%v\ne.TailVertex.Index:\t%v\ne.TailVertex.CheapestEdge:\t%v", e.TailVertex.ID, e.TailVertex.Index, e.TailVertex.CheapestEdge)
		}
	}

	for _, e := range v.EdgesFrom {

		// fmt.Printf("\nBEFORE UPDATE:\ne.TailVertex.ID:\t\t%v\ne.Cost:\t\t\t\t%v\ne.TailVertex.CheapestEdge:\t%v", e.TailVertex.ID, e.Cost, e.TailVertex.CheapestEdge)

		if e.HeadVertex.Index > -1 && (e.Cost < e.HeadVertex.CheapestEdge) {
			e.HeadVertex.CheapestEdge = e.Cost
			heap.Fix(vh, e.HeadVertex.Index)

			// fmt.Printf("\nUPDATED:\ne.HeadVertex.ID:\t\t%v\ne.Cost:\t\t\t\t%v\ne.HeadVertex.CheapestEdge:\t%v", e.HeadVertex.ID, e.Cost, e.HeadVertex.CheapestEdge)
		} else {
			// fmt.Printf("\nNOT UPDATED:\ne.HeadVertex.ID:\t%v\ne.HeadVertex.Index:\t%v\ne.HeadVertex.CheapestEdge:\t%v", e.HeadVertex.ID, e.HeadVertex.Index, e.HeadVertex.CheapestEdge)
		}
	}

	// vh.Peek()
}

var vh VertexHeap

// This example creates a VertexHeap with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {

	readFile(os.Args[1])

	makeVertexHeap()

	minspancost := prim()

	fmt.Printf("\n\n%v", int(minspancost))
}

func readFile(filename string) {

	file, err := os.Open(filename) //should read in file named in CLI

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Scan first line
	if scanner.Scan() {

		discard := scanner.Text()

		fmt.Println(discard)

		// _, err = strconv.Atoi(firstLine)

		// if err != nil {
		// 	fmt.Errorf("Couldn't convert %v", firstLine)
		// }
	}

	for scanner.Scan() {

		thisLine := strings.Fields(scanner.Text())

		HeadVertexID, err := strconv.Atoi(thisLine[0])
		TailVertexID, err := strconv.Atoi(thisLine[1])
		EdgeCost, err := strconv.ParseFloat(thisLine[2], 64)

		if err != nil {
			fmt.Printf("couldn't convert number: %v\n", err)
			return
		}

		h, ok := VertexMap[HeadVertexID]

		if !ok {

			h = &Vertex{
				HeadVertexID,  // Vertex ID		ID           int
				*new([]*Edge), // Edges From		EdgesFrom    []*Edge
				*new([]*Edge), // Edges To		EdgesTo      []*Edge
				math.Inf(1),   // Cheapest edge	CheapestEdge float64 // The cheapest edge that crosses the boundary
				-1,            // Index			Index        int     // Index of item in heap
			}

			VertexMap[HeadVertexID] = h
		}

		t, ok := VertexMap[TailVertexID]

		if !ok {

			t = &Vertex{
				TailVertexID,  // Vertex ID		ID           int
				*new([]*Edge), // Edges From		EdgesFrom    []*Edge
				*new([]*Edge), // Edges To		EdgesTo      []*Edge
				math.Inf(1),   // Cheapest edge	CheapestEdge float64
				-1,            // Index			Index        int
			}

			VertexMap[TailVertexID] = t
		}

		EdgeOneDir := &Edge{HeadVertexID, h, TailVertexID, t, EdgeCost}
		EdgeOtherDir := &Edge{TailVertexID, t, HeadVertexID, h, EdgeCost}

		h.AddToEdge(EdgeOneDir)
		h.AddFromEdge(EdgeOtherDir)
		t.AddFromEdge(EdgeOneDir)
		t.AddToEdge(EdgeOtherDir)

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func makeVertexHeap() {

	vh = make(VertexHeap, len(VertexMap))

	i := 0

	for _, v := range VertexMap {
		v.Index = i
		vh[i] = v
		i++
	}

	heap.Init(&vh)
}

func prim() float64 {

	var MinSpanCost float64

	MinSpanCost = 0.0

	workingVertex := heap.Pop(&vh).(*Vertex)

	vh.UpdateHeap(workingVertex)

	for vh.Len() > 0 {

		workingVertex = heap.Pop(&vh).(*Vertex)

		MinSpanCost += workingVertex.CheapestEdge

		vh.UpdateHeap(workingVertex)

	}

	return MinSpanCost

}
