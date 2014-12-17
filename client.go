package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "strings"
)

const DNS = "ec2-54-149-89-198.us-west-2.compute.amazonaws.com"


func Log(v ...interface{}) {
  fmt.Println(v...)
}

func main() {
  conn, err := net.Dial("tcp", DNS + ":48322")
  if err != nil {
    Log("Error dialing")
    return
  }
  
  quit := make(chan bool)
  // Send messages
  go func() {
    fmt.Printf("Enter username: ")
    reader := bufio.NewReader(os.Stdin)
    username, err := reader.ReadString('\n')
    if err != nil {
      Log("Invalid username")
      return
    }
    conn.Write([]byte(username))
    for {
      reader := bufio.NewReader(os.Stdin)
      message, err := reader.ReadString('\n')
      if err != nil {
        Log("\nExiting")
        quit <- true
      }
      if !strings.EqualFold(message, "\n") {
        conn.Write([]byte(message))
      }
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
      Log("\r" + receivedMessage)
      for i := 0; i < 1024; i++ {
        receivedMessageBuffer[i] = ' '
      }
      _, err = conn.Read([]byte(receivedMessageBuffer))
    }
    if err != nil {
      Log("\nExiting")
      quit <- true
    }
  }()
  if ( <- quit) {
    return
  }
}