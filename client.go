package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
)

const myip = "192.168.56.1"

func main() {
  conn, err := net.Dial("tcp", myip + ":6666")
  if err != nil {
    fmt.Printf("Error dialing\n")
    return
  }
  fmt.Printf("Enter username: ")
  for {
    reader := bufio.NewReader(os.Stdin)
    message, err := reader.ReadString('\n')
    if err != nil {
      fmt.Printf("Error in input\n")
      return
    }
    if len(message) > 0 {
      fmt.Fprintf(conn, message)
    }
  }
}