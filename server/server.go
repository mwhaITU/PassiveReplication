package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	// this has to be the same as the go.mod module,
	// followed by the path to the folder the proto file is in.
	gRPC "github.com/mwhaITU/PassiveReplication/proto"

	"google.golang.org/grpc"
)

type Server struct {
	gRPC.UnimplementedTemplateServer        // You need this line if you have a server
	name                             string // Not required but useful if you want to name your server
	port                             string // Not required but useful if your server needs to know what port it's listening to

	isPrimary      bool
	conns          []*grpc.ClientConn
	incrementValue int64 // value that clients can increment.
	timeSinceBeat  int32
	mutex          sync.Mutex // used to lock the server to avoid race conditions.
}

// Connection to a secondary server
type SecondaryServer struct {
	Address string
	Conn    *grpc.ClientConn
}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var serverName = flag.String("name", "default", "Senders name") // set with "-name <name>" in terminal
var port = flag.String("port", "5400", "Server port")           // set with "-port <port>" in terminal

func main() {

	// setLog() //uncomment this line to log to a log.txt file instead of the console

	// This parses the flags and sets the correct/given corresponding values.
	flag.Parse()
	fmt.Println(".:server is starting:.")

	// starts a goroutine executing the launchServer method.
	go launchServer()

	// This makes sure that the main method is "kept alive"/keeps running
	for {
		time.Sleep(time.Second * 5)
	}
}

func launchServer() {
	log.Printf("Server %s: Attempts to create listener on port %s\n", *serverName, *port)

	// Create listener tcp on given port or default port 5400
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err) //If it fails to listen on the port, run launchServer method again with the next value/port in ports array
		return
	}

	// makes gRPC server using the options
	// you can add options here if you want or remove the options part entirely
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// makes a new server instance using the name and port from the flags.
	server := &Server{
		name:           *serverName,
		port:           *port,
		incrementValue: 0, // gives default value, but not sure if it is necessary
		isPrimary:      *port == "5400",
	}

	if server.isPrimary {
		log.Printf("I AM PRIMARY!!!")
		go server.SendHeartbeats()
	} else {
		log.Printf("I am secondary...")
		go server.CheckForHeartbeat()
	}

	gRPC.RegisterTemplateServer(grpcServer, server) //Registers the server to the gRPC server.

	if server.isPrimary {
		// List of secondary servers
		secondaryServers := []SecondaryServer{
			{Address: "localhost:5401"},
			{Address: "localhost:5402"},
			{Address: "localhost:5403"},
		}
		// Dial the secondary servers
		server.conns, err = DialSecondaryServers(secondaryServers)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			for _, conn := range server.conns {
				conn.Close()
			}
		}()
	}

	log.Printf("Server %s: Listening on port %s\n", *serverName, *port)

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
	// code here is unreachable because grpcServer.Serve occupies the current thread.
}

// Dial a secondary server and return the connection
func DialSecondaryServer(address string) (*grpc.ClientConn, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("could not connect to secondary server: %v", err)
	}

	return conn, nil
}

// Dial all secondary servers and return the connections
func DialSecondaryServers(secondaryServers []SecondaryServer) ([]*grpc.ClientConn, error) {
	conns := make([]*grpc.ClientConn, len(secondaryServers))

	for i, server := range secondaryServers {
		// Dial the secondary server
		conn, err := DialSecondaryServer(server.Address)
		if err != nil {
			return nil, err
		}

		// Save the connection
		conns[i] = conn
	}

	return conns, nil
}

func (s *Server) SendHeartbeat(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	for i, conn := range s.conns {
		currentServer := gRPC.NewTemplateClient(conn)
		ack, err := currentServer.ReceiveHeartbeat(context.Background(), Amount)
		if err != nil {
			log.Printf("Client %s: no response from the server, attempting to reconnect", Amount.GetClientName())
			log.Println(err)
			continue
		}
		// check if the server has handled the request correctly
		if ack.NewValue >= Amount.GetValue() {
			fmt.Printf("Success, the new value is now %d\n on secondary server %v", ack.NewValue, i+1)
			fmt.Println()
		} else {
			// something could be added here to handle the error
			// but hopefully this will never be reached
			fmt.Printf("Oh no something went wrong on server %v :(", i+1)
			fmt.Println()
		}
	}
	return &gRPC.Ack{NewValue: s.incrementValue}, nil
}

func (s *Server) SendHeartbeats() {
	for {
		Amount := gRPC.Amount{
			ClientName: *serverName,
			Value:      s.incrementValue,
		}
		s.SendHeartbeat(context.Background(), &Amount)
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) ReceiveHeartbeat(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	s.timeSinceBeat = 0
	if Amount.GetValue() != s.incrementValue {
		s.incrementValue = Amount.GetValue()
		log.Printf("New value is %v", s.incrementValue)
	}
	return &gRPC.Ack{NewValue: s.incrementValue}, nil
}

func (s *Server) CheckForHeartbeat() {
	for {
		if s.timeSinceBeat >= 10 {
			StartElection()
		}
		s.timeSinceBeat += 1
		time.Sleep(time.Second * 1)
	}
}

func StartElection() {

}

// ChangeServerPort changes the port that a server is listening on
func ChangeServerPort(server *grpc.Server, newPort int) error {
	// Stop the server
	server.Stop()

	// Create a new listener on the new port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", newPort))
	if err != nil {
		return err
	}

	// Start the server on the new listener
	go server.Serve(listener)

	return nil
}

// The method format can be found in the pb.go file. If the format is wrong, the server type will give an error.
func (s *Server) Increment(ctx context.Context, Amount *gRPC.Amount) (*gRPC.Ack, error) {
	// locks the server ensuring no one else can increment the value at the same time.
	// and unlocks the server when the method is done.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// increments the value by the amount given in the request,
	// and returns the new value.
	s.incrementValue += int64(Amount.GetValue())
	log.Printf("New value is %v", s.incrementValue)
	return &gRPC.Ack{NewValue: s.incrementValue}, nil
}
