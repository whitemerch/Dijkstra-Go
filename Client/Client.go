//go run Client.go <graph.txt> <port>

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//This Client program will send the graph file to the server, receive the solutions and write in a .txt file(gonna make a python program for the gui of the graph)
//The function to retrive the arguments given and to see if everything is alright with them
func args() (int, string) {
	if len(os.Args) != 3 {
		log.Fatal("Your input is invalid. You either put too much or less arguments than needed. It should be in this form: go run Client.go <graph.txt> <port>")
		//Fatal is equivalent to os.Exit()
		//It exits the program if you give more or less arguments than needed
	} else {
		fmt.Printf("#DEBUG ARGS Port Number : %s\n", os.Args[2])
		portNumber, err := strconv.Atoi(os.Args[2])
		//It converts our port string into an int
		if err != nil {
			log.Fatal("The port number you gave is invalid. It should be a number")
			//If the port given isn't convertible(E.g '/' not convertible into an int)
		} else {
			filename := os.Args[1]
			_, err := os.Stat(filename)
			//os.Stat is to see if the file exists. If it does, err=nil
			if err != nil {
				log.Fatalf("File %v not found ", filename)
				//File not found
			} else {
				return portNumber, filename
			}
		}
	}
	return -1, ""
	//We are obligated to give a return there
}

func main() {
	start := time.Now()
	port, filename := args()
	fmt.Printf("#DEBUG DIALING TCP Server on port %d\n", port)
	portString := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	//Sprintf is a function that takes two strings and fuse them into one string as we want here
	//strconv.Itoa convert the port which is an int to a string
	fmt.Printf("#DEBUG MAIN PORT STRING |%s|\n", portString)
	connection, err := net.Dial("tcp", portString)
	//The Dial function connects to a server
	if err != nil {
		log.Fatal("#DEBUG MAIN could not connect")
	} else {
		defer connection.Close()
		//The connection closes at the end once every instruction is done
		server := bufio.NewReader(connection)
		//We create a reader there to listen to what the server has to say
		//It's important to put the reader there so afterwards we can get the answer of the server for our problem
		fmt.Printf("#DEBUG MAIN connected\n")
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(("Failed opening file"))
			//Fatalf is equivalent to os.Exit()
		} else {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			//We create a new scanner to scan our file(we give it a buffer)
			for scanner.Scan() { //The loop works till there is no more lines to scan
				txt := scanner.Text()
				//Text reads each line
				io.WriteString(connection, txt+"\n") //We send to the server each line of the text and we go to the next line
			}
			err := scanner.Err()
			//To see if an error occurs during the scanning of the file
			if err != nil {
				log.Fatal("An error occured during the reading of the file")
			} else {
				solutioname := fmt.Sprintf("solution__%s", filename)
				//We did this to get different .txt files depending on the graph given and not always do it on the same .txt file
				solution, err := os.OpenFile(solutioname, os.O_CREATE|os.O_RDWR, 0755)
				if err != nil {
					log.Fatalf("Failed creating file: %s", err)
				} else {
					defer solution.Close()
					results, err := server.ReadString('#') //We take what is being sent by the server with the reader we created before
					if err != nil {
						log.Fatalf("End of the program \n")
					}
					results = strings.TrimSuffix(results, "#")
					//fmt.Println(results)
					_, err = solution.WriteString(results)
					fmt.Println(time.Since(start))
					if err != nil {
						log.Fatal("Failed writing the solutions file")
					}
				}
			}
		}
	}
}
