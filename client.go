package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "strings"
)

const serverIP = "192.168.56.1"

func Logln(v ...interface{}) {
  fmt.Println(v...)
}

func main() {
  conn, err := net.Dial("tcp", serverIP + ":6666")
  if err != nil {
    Logln("Error dialing")
    return
  }
  
  quit := make(chan bool)
  // Send messages
  go func() {
    fmt.Printf("Enter username: ")
    reader := bufio.NewReader(os.Stdin)
    username, err := reader.ReadString('\n')
    if err != nil {
      Logln("Invalid username")
    }
    conn.Write([]byte(username))
    for {
      //fmt.Printf("You: ")
      reader := bufio.NewReader(os.Stdin)
      message, err := reader.ReadString('\n')
      if err != nil {
        Logln("\nExiting")
        quit <- true
      }
      conn.Write([]byte(message))
    }
  }()
  // Receive messages
  go func() {
    receivedMessageBuffer := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
      receivedMessageBuffer[i] = ' '
    }
    _, err = conn.Read([]byte(receivedMessageBuffer))
    for err == nil {
      receivedMessage := strings.TrimSpace(string(receivedMessageBuffer))
      Logln(receivedMessage)
      for i := 0; i < 1024; i++ {
        receivedMessageBuffer[i] = ' '
      }
      _, err = conn.Read([]byte(receivedMessageBuffer))
    }
    if err != nil {
      Logln("\nExiting")
      quit <- true
    }
  }()
  if ( <- quit) {
    return
  }
}