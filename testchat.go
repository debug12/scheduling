package main

import(
  "net"
  "fmt"
  "strings"
)

func main() {
  listener, err := net.Listen("tcp", ":6666")
  if err != nil {
    fmt.Println("Error accepting connection")
  }
  
  
  outgoingMessages := make(chan string)
  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Println("Error accepting connection")
    }
    fmt.Println("Accepted a connection")
    userBuffer := make([]byte, 10)
    for i := 0; i < 10; i++ {
      userBuffer[i] = ' '
    }
    conn.Read(userBuffer)
    user := strings.Fields(string(userBuffer))
    username := user[0]
    buffer := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
      buffer[i] = ' '
    }
    // Receive and print messages
    go func() {
      _, err := conn.Read(buffer)
      for err == nil {
        message := username + ": " + strings.TrimSpace(string(buffer))
        fmt.Println(message)
        outgoingMessages <- message
        for i := 0; i < 1024; i++ {
          buffer[i] = ' '
        }
        _, err := conn.Read(buffer)
        if err != nil {
          fmt.Println("Ending connection with", username)
          return
        }
      }
      fmt.Println("Exiting")
    }()
    // Relay messages to clients
    go func() {
      for {
        conn.Write([]byte(<- outgoingMessages))
      }
    }()
  }
}