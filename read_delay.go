package main

//import "net" 
import "log"
import "fmt" 
import "time"
import "flag"
import "runtime"
import "github.com/go-amp/amp"

var NUM_REQUESTS *int

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

func server() {
    amp.ListenTCP(":8000")
    KeepAlive()
}

func client() {
    c, err := amp.ConnectTCP(*isClientHost)
    if err != nil { return }    
    go send_requests(c)
    KeepAlive()    
}

func send_requests(c *amp.Client) {
    test_start := time.Now()       
    for i := 1; i <= *NUM_REQUESTS; i++ {
        send := []byte{0,1,97,0,6,54,54,50,55,49,54,0,1,98,0,1,48,0,4,95,97,115,107,0,5,97,49,99,98,99,0,8,95,99,111,109,109,97,110,100,0,3,83,117,109,0,0}
        //log.Println("writing",send)
        _, err := c.Conn.Write(send)
        if err != nil { log.Println("err",err) }
        sent_count++
        runtime.Gosched()
    }
    now := time.Now()
    fmt.Printf("time taken -- %f\n", float32(now.Sub(test_start))/1000000000.0)
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

