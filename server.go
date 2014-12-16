package main

import(
  "net"
  "fmt"
  "strings"
  "bytes"
  "container/list"
)

func Log(v ...interface{}) {
  fmt.Println(v...)
}

type Client struct {
  Conn net.Conn
  Username string
  OutgoingMessages chan string
  Quit chan bool
}

func (c Client) receiveMessages() {
  buffer := make([]byte, 1024)
  for i := 0; i < 1024; i++ {
    buffer[i] = ' '
  }
  _, err := c.Conn.Read(buffer)
  for err == nil {
    message := c.Username + ": " + strings.TrimSpace(string(buffer))
    //Log(message)
    c.OutgoingMessages <- message
    for i := 0; i < 1024; i++ {
      buffer[i] = ' '
    }
    _, err := c.Conn.Read(buffer)
    if err != nil {
      Log("> Ending connection with", c.Username)
      c.Quit <- true
      c.OutgoingMessages <- ("> " + c.Username + " has left\n")
      return
    }
  }
  Log("> Exiting")
}

func main() {
  listener, err := net.Listen("tcp", ":6666")
  if err != nil {
    Log("> Error accepting connection")
  }
  listOfClients := list.New()
  
  for {
    conn, err := listener.Accept()
    if err != nil {
      Log("> Error accepting connection")
    }
    Log("> Accepted a connection")
    userBuffer := make([]byte, 10)
    for i := 0; i < 10; i++ {
      userBuffer[i] = ' '
    }
    conn.Read(userBuffer)
    user := strings.Fields(string(userBuffer))
    username := user[0]
    for it := listOfClients.Front(); it != nil; it = it.Next(){
      client := it.Value.(Client)
      if bytes.Equal([]byte(client.Username), []byte(username)){
        conn.Close()
        break
      }
    }
    newClient := Client{conn, username, make(chan string), make(chan bool)}
    listOfClients.PushBack(newClient)
    
    // Receive and print messages from newClient
    go newClient.receiveMessages()
    
    // Relay from newClient to others
    go func() {
      for {
        select {
          case outgoingMessage := <- newClient.OutgoingMessages:
            for it := listOfClients.Front(); it != nil; it = it.Next(){
              client := it.Value.(Client)
              if !bytes.Equal([]byte(client.Username), []byte(newClient.Username)) {
                client.Conn.Write([]byte(outgoingMessage))
              }
            }
          case <- newClient.Quit:
            for it := listOfClients.Front(); it != nil; it = it.Next() {
              client := it.Value.(Client)
              if bytes.Equal([]byte(client.Username), []byte(newClient.Username)) {
                listOfClients.Remove(it)
              }
            }
            outgoingMessage := <- newClient.OutgoingMessages
            for it := listOfClients.Front(); it != nil; it = it.Next() {
              client := it.Value.(Client) 
              client.Conn.Write([]byte(outgoingMessage))
            }
            return
        }
      }
    }()
    
    newClient.OutgoingMessages <- ("> " + newClient.Username + " has entered\n")
  }
}