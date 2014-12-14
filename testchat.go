package main

import(
  "net"
  "fmt"
  "strings"
)

func main() {
  listener, err := net.Listen("tcp", ":6666")
  if err != nil {
    fmt.Println("Error accepting connection\n")
  }
  
  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println("Error accepting connection\n")
    }
    fmt.Println("Accepted a connection\n")
    userbuffer := make([]byte, 10)
    conn.Read(userbuffer)
    user := strings.Fields(string(userbuffer))
    //fmt.Println(user)
    go func() {
      for {
        buffer := make([]byte, 80)
        conn.Read(buffer)
        fmt.Println(user[0] + ":", strings.TrimSpace(string(buffer)))
      }
    }()
  }
}