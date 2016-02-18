package main
import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"net"
	"time"
	"strconv"
)

type Node struct {
	ip   string
	port string
	receive bool 
	send bool
	back bool
}
var sendBackFlag bool
var iniFlag bool
var sendIniMessFlag bool
var Iam Node
var Parent Node
var neighbors []Node
var leader string

func main() {
	filename := `configuration.conf`
	fmt.Println("Start..."  )
	// Find Iam,Initiator,Neighbors
	readFile(filename)
	fmt.Printf("I am %s:%s \niniFlag is: %t \nAll my neighbors are: %v \n" , Iam.ip , Iam.port , iniFlag , neighbors)
	go server(Iam)
	
	if checkNeighborServer(neighbors) {
		done := false
		for {
			if iniFlag && !sendIniMessFlag {
				fmt.Println("Start to send message from initiator: "  )
				iniMss := "&Iam="+Iam.ip+":"+Iam.port+"&leader="+leader+"&back=false"
				sendMssToAllNeighbors(iniMss)
				sendIniMessFlag = true
				fmt.Printf("All my neighbors are: %v \n" , neighbors)
			}else{
				time.Sleep(3000 * time.Millisecond)
				fmt.Print("." )
				if checkReceiveFromAll() {
					if iniFlag {
						fmt.Println("\nDone "  )
						done = true
					}else{
						m := "&Iam="+Iam.ip+":"+Iam.port+"&leader="+leader+"&back=true"
						sendMessage(m,Parent)
						fmt.Println("\nDone "  )
						done = true
					}
				}
			}
			if done {
				fmt.Printf( "Neighbors: %v \nMy Leader is: %s" , neighbors,leader )
				break
			}
		}
	}
	
}

func readFile(fileName string){
	f, _ := os.Open(fileName)
	defer f.Close()
	r := bufio.NewReaderSize(f, 2*1024)
	line, isPrefix, err := r.ReadLine()
	i := 1
	for err == nil && !isPrefix {
		s := string(line)
		if i == 1 {
				// Find Iam
				t :=strings.Split(s, ":")
				Iam = Node{t[0],t[1],false,false,false}
		}else{
			k :=strings.Split(s, ":")
			if k[0] == "initiator" {
					// Find if initiator
					iniFlag = true
					leader = k[2]
					fmt.Println("I've just set the Leader: " , k[2])
				}else{
					// Find neighbors
					neighbors = append(neighbors, Node{k[0],k[1],false,false,false})
				}		
		}
		i++
		line, isPrefix, err = r.ReadLine()		
	}
}

func analizMessage(message string) map[string]string{
	ms :=strings.Split(message, "&")
	GetMessage := make(map[string]string)
	
	msIam :=strings.Split(ms[1], "=")
	mx :=strings.Split(msIam[1], ":")
	GetMessage["ip"] =  mx[0]
	GetMessage["port"] =  mx[1]
	
	getLeader :=strings.Split(ms[2], "=")
	GetMessage["leader"] =  getLeader[1]
	
	getBack :=strings.Split(ms[3], "=")
	GetMessage["back"] =  getBack[1]
		
	fmt.Println("GetMessage-> ",GetMessage)
	return GetMessage
}

func server(s Node) {
	fmt.Printf("Launching server... %s:%s \n" , s.ip,s.port)
	ln, _ := net.Listen("tcp", s.ip+":"+s.port)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n') 
		if string(message) != "" {
			fmt.Println("->", string(message))
			doIt(analizMessage(message))
		}
	}	
	
}

func checkNeighborServer(n []Node) bool{
	for i:=0; i < len(n);i++{
		for {
			conn, err := net.Dial("tcp", n[i].ip+":"+n[i].port)
			fmt.Println("Looking for " + n[i].ip+":"+n[i].port)
			time.Sleep(3000 * time.Millisecond)
			if err == nil {
				conn.Close()
				break
			}
		}
	}
	
	return true
}

func sendMessage(s string, n Node){
	conn, _ := net.Dial("tcp", n.ip+":"+n.port)
	defer conn.Close()
	conn.Write([]byte(s))
	fmt.Printf("Message Sent to %s:%s \n" ,n.ip,n.port )	
}

func sendMssToAllNeighbors(ms string){

	for i:=0; i < len(neighbors);i++{
		if Parent.port != neighbors[i].port {
			sendMessage(ms,neighbors[i])
			neighbors[i].send = true	
		}
		
	}
}

func findNewLeadership( l string) bool{
	myReturn:=false
	if leader != "" {
		lastLeader,_ := strconv.Atoi(leader)
		newLeader,_ := strconv.Atoi(l)
		if  lastLeader < newLeader {
			fmt.Println("I'm going to change my leader to: ", l )
			leader = l
			myReturn = true
		}
	}else{
		// First leader
		leader = l
	}
	return myReturn
}
func deleteAllActivities(){
	fmt.Println("I'm going to delete all previous activities: " )
	for i:=0; i < len(neighbors);i++{
		fmt.Println("*** Delete all activities about " +neighbors[i].ip + ":" + neighbors[i].port )
		neighbors[i].receive = false
		neighbors[i].send = false
		neighbors[i].back = false
	}
	iniFlag = false
}
func doIt( ms map[string]string){
	//find sender
	Sender,id := findNodeBtwNeighbors(ms["ip"],ms["port"])
	if ms["back"] == "true" {
		neighbors[id].back = true
	}else{
		if findNewLeadership(ms["leader"]){
			//delete all activities
			deleteAllActivities()
			
			fmt.Println("I'm going to change my parent to: ",  ms["ip"],ms["port"])
			Parent =  Sender
			neighbors[id].receive = true
		}else{	
			if iniFlag{
				neighbors[id].receive = true
			}else{
				if Parent.ip == "" {
					Parent =  Sender
					sendMss := "&Iam="+Iam.ip+":"+Iam.port+"&leader="+leader+"&back=false"
					neighbors[id].receive = true
					sendMssToAllNeighbors(sendMss)
				}else{
					neighbors[id].receive = true
				}
			}
		}		
	}
}

func findNodeBtwNeighbors(ip string, port string) (Node , int){
	j := 0
	for i:=0; i < len(neighbors);i++{
		if neighbors[i].ip == ip && neighbors[i].port == port {
				j = i
		}
	}
	return neighbors[j],j
}

func checkReceiveFromAll() bool{
	Myreturn := true
	for _, n := range neighbors {
		if Parent.port != n.port   {
			if !n.back {
				Myreturn = false
			}
		}
	}
	return Myreturn
}
