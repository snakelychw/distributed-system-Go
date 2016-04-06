package main

import(
  "net"
  "fmt"
  "bufio"
  "strings"
  "strconv"
  "io"
  "os"
)
var NewConn = make(chan net.Conn)
var GlobalIn = make(chan string)
var GlobalBuf = make(chan string)
//User class
type user struct{
  ID int
  connection net.Conn
  in chan string
  out chan string
  LocalReader *bufio.Reader
  LocalWriter *bufio.Writer
}

//ChatRoom class
type ChatRoom struct{
  count int
  UserList []*user
  in chan string
}
//User constructor
func NewUser(con net.Conn, count int) *user{
    temp := new(user)
    temp.ID = count
    fmt.Println("User",temp.ID,"joined Chitter.")
    temp.in = make(chan string)
    temp.out = make(chan string)
    temp.LocalReader = bufio.NewReader(con)
    temp.LocalWriter = bufio.NewWriter(con)
    go temp.ReadListen()
    go temp.WriteListen()
    return temp
}
//User input lisertener, connection lost is hanlded here
func (ur *user) ReadListen(){
  GlobalIn <- fmt.Sprintf("User %d joined. You can PM to him/her\n", ur.ID)
  for{
    msg, err := ur.LocalReader.ReadString('\n')
    if err == io.EOF{ // handle connection lost
      GlobalIn <- fmt.Sprintf("User %d left.\n", ur.ID)
      fmt.Println(fmt.Sprintf("User %d left Chitter.", ur.ID))
      ur.connection.Close()
      return
    }
    s := strings.SplitN(msg, ":",2)
    digi, ok := strconv.Atoi(strings.Trim(s[0], " "))
    if ok == nil{//handle pm
      msg = strconv.Itoa(digi)+":"+strconv.Itoa(ur.ID)+":"
      msg = msg+s[1]
    } else if strings.Trim(s[0], " ") == "all" { //handle all with "all:"
      msg = "all:"+strconv.Itoa(ur.ID)+":"
      msg = msg+s[1]
    } else if strings.Trim(s[0], " ") == "whoami" {//whoami
      msg = strconv.Itoa(ur.ID)+":chitter: "+ strconv.Itoa(ur.ID)+"\n"
    } else {//handle all without "all:"
      msg = "all:" +strconv.Itoa(ur.ID)+":" + msg
    }
    GlobalBuf <- msg
  }
}
//User output listener
func (ur *user) WriteListen(){
  for msg:= range ur.out{
    ur.LocalWriter.WriteString(msg)
    ur.LocalWriter.Flush()
  }
}

//ChatRoom constructor
func NewChatRoom() *ChatRoom{
  temp := new(ChatRoom)
  temp.UserList = make([]*user, 0)
  temp.count = 0
  return temp
}
//adding new user to  ChatRoom, listen to messeges
func (cr *ChatRoom) MainTask(){
  fmt.Println("Chitter started. Have Fun! ")
  for{
    select{
      case conn := <- NewConn:
        cr.UserList = append(cr.UserList, NewUser(conn, cr.count))
        cr.UserList[cr.count].connection = conn
        cr.count++
      case msg := <- GlobalBuf:
        s := strings.SplitN(msg, ":", 2)
        digi, ok := strconv.Atoi(strings.Trim(s[0], " "))
        if ok == nil && digi <= cr.count{
          cr.UserList[digi].out <- s[1]
          } else if s[0] == "all"{
            for _, u := range cr.UserList{
                u.out <- s[1]
            }
        }
      case content:=  <- GlobalIn:
        for _, u := range cr.UserList{
            u.out <- content
        }
      }
    }
}


func main() {
    chitter := NewChatRoom()
    port := os.Args[1]
    port = ":"+port
    listerner, err1 := net.Listen("tcp", port)
    if(err1 != nil){
      fmt.Println("Server Failed...")
    }
    go chitter.MainTask()

    for{
      conn, err2 := listerner.Accept()
      if err2 != nil{
        fmt.Println("Connection Failed...")
      }
      NewConn <- conn

    }
}
