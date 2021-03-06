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
  c.Quit <- true
  Log("> Exiting")
}

func main() {
  listener, err := net.Listen("tcp", ":48322")
  if err != nil {
    Log("> Error listening for connection")
    return
  }
  listOfClients := list.New()
  
  for {
    conn, err := listener.Accept()
    if err != nil {
      Log("> Error accepting connection")
      return
    }
    Log("> Accepted a connection")
    userBuffer := make([]byte, 10)
    for i := 0; i < 10; i++ {
      userBuffer[i] = ' '
    }
    _, errUsername := conn.Read(userBuffer)
    if errUsername != nil {
      continue
    }
    user := strings.Fields(string(userBuffer))
    username := user[0]
    skipOneIteration := false
    for it := listOfClients.Front(); it != nil; it = it.Next(){
      client := it.Value.(Client)
      if bytes.Equal([]byte(client.Username), []byte(username)){
        conn.Close()
        skipOneIteration = true
        break
      }
    }
    if skipOneIteration {
      continue
    }
    newClient := Client{conn, username, make(chan string), make(chan bool)}
    listOfClients.PushBack(newClient)
    
    var csvUsernamesBuffer bytes.Buffer
    csvUsernamesBuffer.WriteString("Currently connected: ")
    for it := listOfClients.Front(); it != nil; it = it.Next() {
      if it != listOfClients.Front() {
        csvUsernamesBuffer.WriteString(", ")
      }
      client := it.Value.(Client)
      csvUsernamesBuffer.WriteString(client.Username)
    }
    newClient.Conn.Write([]byte(csvUsernamesBuffer.String()))
    
    // Receive messages from newClient
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