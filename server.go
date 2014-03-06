package main

import "net"
import "log"
import "runtime"
import "bufio"
import "time"
import "encoding/binary"

const PREFIXLENGTH = 2

func scan(reader *bufio.Reader, v []byte) error {
    i := 0
    for {
        n, err := reader.Read(v[i:])
        if err != nil { return err }
        i += n
        if i == len(v) { return nil }
    }
    return nil
}

func get(reader *bufio.Reader, m map[string][]byte) error {
    prefix := make([]byte, 2)
    l := 0
    var err error
    for {                
        //k        
        err = scan(reader, prefix[:])        
        if err != nil { return err }
        l = int(binary.BigEndian.Uint16(prefix)) 
        // indicates end of message 
        if l == 0 { return nil }              
        
        k := make([]byte, l)
        err = scan(reader, k[:])                
        if err != nil { return err }
        
        //v        
        err = scan(reader, prefix[:])                
        if err != nil { return err }
        l = int(binary.BigEndian.Uint16(prefix))        
                
        v := make([]byte, l)
        err = scan(reader, v[:])                
        if err != nil { return err }
        
        // assign
        m[string(k)] = v        
    }
    return nil
}

func read(conn net.Conn) {
    reader := bufio.NewReader(conn)
    
    startTime := time.Now()    
    for {    
        m := make(map[string][]byte)
        err := get(reader, m)
        if err != nil { log.Println(err); break }
        log.Println(m)        
    }        
    endTime := time.Now()    
    log.Println("ElapsedTime:", endTime.Sub(startTime))    
}

func init() {
    procs := runtime.NumCPU()
    log.Println("setting GOMAXPROCS to",procs)
    runtime.GOMAXPROCS(procs)
}


func main() {
    service := ":8000"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service) 
    if err != nil { log.Fatal(err) }
    l, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
    if err != nil { log.Fatal(err) }
    defer l.Close()
    
    for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		
		go read(conn)
	}
}
