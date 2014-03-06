package main

import "net"
import "log"
import "runtime"
import "bufio"
import "time"
import "encoding/binary"

const PREFIXLENGTH = 2

func init() {
    procs := runtime.NumCPU()
    log.Println("setting GOMAXPROCS to",procs)
    runtime.GOMAXPROCS(procs)
}

func pack(m map[string][]byte) *[]byte {
    l := 0
    for k, v := range m {         
        l += len(k)
        l += PREFIXLENGTH
        l += len(v)
        l += PREFIXLENGTH        
    }
    
    var r = make([]byte, l + PREFIXLENGTH)
    i := 0
    for k, v := range m {
        //k        
        l = len(k)
        binary.BigEndian.PutUint16(r[i:i+PREFIXLENGTH], uint16(l))
        i += PREFIXLENGTH
        copy(r[i:i+l],k)
        i += l
        //v        
        l = len(v)
        binary.BigEndian.PutUint16(r[i:i+PREFIXLENGTH], uint16(l))
        i += PREFIXLENGTH
        copy(r[i:i+l],v)
        i += l
    }
    return &r
}

func main() {
    service := "10.38.207.141:8000"
    serverAddr, err := net.ResolveTCPAddr("tcp4", service)
    if err != nil { log.Fatal(err) }
    conn, err := net.DialTCP("tcp4", nil, serverAddr)    
    if err != nil { log.Fatal(err) }
    
    //conn.Write("hello")
    writer := bufio.NewWriter(conn)
    m := make(map[string][]byte)    
      
    //buf := make([]byte, 500)
    startTime := time.Now()
    for i := 0; i < 10000; i++ {
        m["i"] = []byte("i need a longer packet omg so what is umm.... this over here how long does it take")    
        buf := *pack(m)          
        writer.Write(buf)
    }
    writer.Flush()
    endTime := time.Now()
    log.Println("ElapsedTime:", endTime.Sub(startTime))
    log.Println("done..")
}
