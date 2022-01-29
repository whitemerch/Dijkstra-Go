package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Edge struct {
	from string
	to   string
	cost int
}

type way struct {
	from string
	cost int
}

type result struct {
	from string
	to   string
	list []string
	cost string
}

//go run Server_string_go.go <port>
func getArgs() int {
	if len(os.Args) != 2 {
		log.Fatal("Your input is invalid. You either put too much or less arguments than needed. It should be in this form: go run Server.go <port>")
	} else {
		//We convert the port given which is a string into an int
		port, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal("The port given is invalid")
		} else {
			return port
		}
	}
	return -1
	//Should never be returned
}

func reverse(slice []string) []string { //There are no implemented library to reverse slices
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func getNeighbors(node string, list []Edge) []Edge { //[{A B 1} {A C 2} ... ]
	//We build a slice of the neighbours of the node we give in parenthesis
	var neighbors []Edge
	for _, elt := range list { //Since the slice graph contains all the edges we have ([{A B 8] {B C 9}...]for e.g), we run throug our slice, elt here being a complete edge({A B 9})
		if elt.from == node { //If elt.from means the first node of the edge type, if it corresponds to our point given then we add it to our neighbour slice
			neighbors = append(neighbors, elt)
		}
	}
	return neighbors
}
func getAllNeighbors(list []Edge, nodes []string) map[string][]Edge {
	allNeighbors := make(map[string][]Edge)
	for _, node := range nodes {
		allNeighbors[node] = getNeighbors(node, list)
	} //For each node in our nodes list, we will use our previous function to get every edge of it so we will get a map like that for e.g
	//{A:[{A B 1} {A C 2}... ],B:[...]}
	//It will help us by sorting everything being given by the .txt file
	return allNeighbors
}

//This function is needed when we have to build our path back by finding between the paths which one was the shortest so that we choose it
func getMin(graphPart []way) way {

	minKey := 0
	minValue := graphPart[0].cost
	for k, elt := range graphPart {
		if elt.cost < minValue {
			minValue = elt.cost
			minKey = k
		}
	}
	return graphPart[minKey]
}

// This function gets us the nearest node from a certain node and the cost of the actual state
func getMinDijk(table map[string][]way, done map[string]int) (string, int) {
	min := -1
	minPoint := ""
	minKey := 0
	for point, i := range table { //point the key, i the value(slice of ways)
		if _, ok := done[point]; !ok {
			for k, chm := range i {
				if min != -1 && chm.cost < min {
					min = chm.cost
					minPoint = point
					minKey = k
				}
				if min == -1 {
					min = chm.cost
					minPoint = point
					minKey = k
				}
			}
		}
	}
	return minPoint, minKey
}

//Our principle function which gets every shortest path from a node given to all the others in the graph
func oneDijkstra(from string, wg *sync.WaitGroup, list []Edge, nodes []string, neighbors map[string][]Edge, dataCh chan<- []result) {
	defer wg.Done()                 //We close our waitgroup once everything is done
	table := make(map[string][]way) //Our table which contains our work(similar to table from the pdf)
	done := make(map[string]int)    //Once we passed through a certain node, we add it to this table(it's like the scratches under the node we chose in the pdf)

	table[from] = append(table[from], way{from, 0})
	//{1:[{1 0}]}

	for i := 0; i < len(nodes); i++ { //Number of steps according to the number of nodes(like our explanation in the pdf)

		pt, k := getMinDijk(table, done) //It will give us the nearest node of our actual node with its cost
		smallestWay := table[pt][k]

		done[pt] = i
		for _, direction := range neighbors[pt] {
			if _, ok := done[direction.to]; !ok { //If the direction we can go to is already in done, then ok will give true, !ok false and will not enter the condition
				table[direction.to] = append(table[direction.to], way{direction.from, direction.cost + smallestWay.cost})
			}
		}
	}
	var results []result
	for _, node := range nodes {
		theway := []string{}
		theway = append(theway, node)
		n := node
		if len(table[n]) > 0 {
			for getMin(table[n]).from != from {
				theway = append(theway, getMin(table[n]).from)
				n = getMin(table[n]).from
			}
			theway = append(theway, from)
			theway = reverse(theway)
			//I reverse my slice to get my path in the good order
			results = append(results, result{from, node, theway, strconv.Itoa(getMin(table[node]).cost)})
		}
	}
	dataCh <- results
}

//This function will get us all the shortest paths from all the nodes to all the other nodes
//list []Edge here is our list containing all the characteristics given by the .txt file
//nodes []string gives us all the nodes in our graph
func Dijkstra(list []Edge, nodes []string) []result {
	var wg sync.WaitGroup // Waitgroup so that we won't get some things done before all the goroutines are done
	dataCh := make(chan []result)
	done := make(chan bool)
	var results []result
	start := time.Now()
	neighbors := getAllNeighbors(list, nodes)
	//We get all the neighbors for every node we have in our graph
	//{1:[{1 2 1},{1 3 2}],2:...}
	//start a goroutine that will read from the data channel
	go func() {
		for d := range dataCh {
			results = append(results, d...)
		}
		done <- true //this is used to wait until all data has been read from the channel
		//true once everything has been received
	}()

	wg.Add(len(nodes))
	for _, node := range nodes {
		go oneDijkstra(node, &wg, list, nodes, neighbors, dataCh)
	}

	wg.Wait()     //We wait until the waitgroup is empty
	close(dataCh) //this closes the dataCh channel, which will make the for-range loop exit once all the data has been read
	<-done        //Once done is true close the dataCh
	//We listen to the channel until it tells us true so we can continue
	fmt.Println(time.Since(start))
	return results
}

func handleConnection(connection net.Conn, connum int) {
	defer connection.Close()
	fmt.Println("Connection accepted")

	connReader := bufio.NewReader(connection)
	//We put a reader on our connection
	var list []Edge
	var nodes []string
	for {
		inputLine, err := connReader.ReadString('\n')
		//While listening to our connection, we receive the line being sent by the client
		//We read until there is a \n
		if err != nil {
			fmt.Printf("#DEBUG %d RCV ERROR no panic, just a client\n", connum)
			fmt.Printf("Error :|%s|\n", err.Error())
			break
		}
		inputLine = strings.TrimSuffix(inputLine, "\n")
		//We take away the /n of the line we receive
		//fmt.Printf("#DEBUG %d RCV |%s|\n", connum, inputLine)
		splitLine := strings.Split(inputLine, " ")
		//We split our string(First Node Last Node Cost) into a slice
		if splitLine[0] != "." { //The user should always mark the end of the .txt file by a point
			nodes = append(nodes, splitLine[0], splitLine[1])
			cost, _ := strconv.Atoi(splitLine[2])
			list = append(list, Edge{splitLine[0], splitLine[1], cost})
		} else {
			break //To break our infinite loop when there is no more lines
		}
	}
	nodes = unique(nodes)
	//Unique there sort the slice so we can eliminate the duplicate edges
	sort.Strings(nodes)
	//Sorting our slice in an ascending order
	routes := Dijkstra(list, nodes)
	io.WriteString(connection, fmt.Sprintf("%s # ", routes))
}

func unique(a []string) []string {

	check := make(map[string]int)
	res := make([]string, 0) //A new slice
	for _, val := range a {
		check[val] = 1 //Each letter of the original slice is being added in our map with a value of 1
		//When we build our map there is only one letter of each that's how we get our slice without duplicates
	}

	for letter := range check {
		res = append(res, letter)
	}
	//We get our slices without duplicates
	return res
}

func main() {
	port := getArgs()
	fmt.Printf("#DEBUG MAIN Creating TCP Server on port %d\n", port)
	portString := fmt.Sprintf(":%s", strconv.Itoa(port))
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)
	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Printf("#DEBUG MAIN Could not create listener\n")
		panic(err)
		//Panic is like log.fatal or os.Exit(1)
	}
	connum := 1
	for {
		fmt.Printf("#DEBUG MAIN Accepting next connection\n")
		conn, errconn := ln.Accept()
		//Once a connection is being made from our client file we accept it
		if errconn != nil {
			fmt.Printf("DEBUG MAIN Error when accepting next connection\n")
			panic(errconn)
		}
		//If we're here, we did not panic and conn is a valid handler to the new connection
		go handleConnection(conn, connum)
		connum += 1
	}
}
