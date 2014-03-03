package main

import "net" 
import "log"
import "fmt" 
import "time"
import "flag"
import "runtime"

var NUM_REQUESTS *int
const READ_BUFFER_SIZE int = 65535
var isClient *bool
var isServer *bool
var isClientHost *string
var sent_count int = 0
var requests_count int = 0

func KeepAlive() {
    for { 
        runtime.Gosched()
        time.Sleep(1 * time.Second) 
        if *isClient { 
            log.Println("sent",sent_count)             
        } else { log.Println("requests",requests_count) }
    }
}

type Client struct {
    Name *string
    Conn *net.TCPConn
}

func ClientCreator(name *string, conn *net.TCPConn) *Client {
    client := &Client{name, conn} 
    go client.Reader()
    return client
}

func connectionListener(netListen *net.TCPListener, service string) {    
    defer netListen.Close()
    log.Println("Waiting for clients") 
    for {
        conn, err := netListen.AcceptTCP()
        if err != nil {
            log.Println("Client error: ", err)
            break
        } else {            
            name := fmt.Sprintf("<%s<-%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())
            log.Println("AMP.connectionListener accepted",name)            
            
            ClientCreator(&name, conn)
            
        }
    }
}

func ListenTCP(service string) error {
    tcpAddr, err := net.ResolveTCPAddr("tcp", service) 
    if err != nil {
        log.Println("Error: Could not resolve address")
        return err
    } else {
        log.Println("ListenTCP",*tcpAddr)
        netListen, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
        if err != nil {
            log.Println("Error: could not listen")
            return err
        } else {
            go connectionListener(netListen, service)
       }
    }
    return nil
}

func (c *Client) Reader() {    
    buffer := make([]byte, READ_BUFFER_SIZE)
    for {
        log.Println("ready for new read..")
        readBytes, err := c.Conn.Read(buffer) 
        log.Println("received bytes",readBytes)
        if err != nil {
            log.Println("connection reader error!!",err)            
        }        
        time.Sleep(100 * time.Millisecond)                        
    }
}

func ConnectTCP(service string) (*Client, error) {    
    
    serverAddr, err := net.ResolveTCPAddr("tcp", service)
    if err != nil {
        log.Println("error!",err)
        return nil, err
    }
    conn, err := net.DialTCP("tcp", nil, serverAddr)
    if err != nil {
        log.Println("error!",err)
        return nil, err
    }
    name := fmt.Sprintf("<%s->%s>", conn.LocalAddr().String(), conn.RemoteAddr().String())    
    log.Println("AMP.ConnectTCP connected",name)
        
    newClient := ClientCreator(&name, conn)     
    return newClient, nil
}

func server() {
    ListenTCP(":8000")
    KeepAlive()
}

func client() {
    c, err := ConnectTCP(*isClientHost)
    log.Println("c",c,err)
}

func main() {
    isServer = flag.Bool("server", false, "use as a server")
    isClient = flag.Bool("client", false, "use as a client")
    isClientHost = flag.String("host","127.0.0.1:8000","host address")
    NUM_REQUESTS = flag.Int("num",100000,"number of requests to do")
    log.Println("hi")
    flag.Parse()    
    if *isServer {
        server()
    } else if *isClient {        
        client()
    } else { flag.Usage() }  
}
