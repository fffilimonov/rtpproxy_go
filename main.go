package main

import (
    "flag"
    "fmt"
    "net"
    "strings"
)

var Debug bool

//just wrapper
func ConcatStr (sep string, args ... string) string {
    return strings.Join(args, sep)
}

func main() {

//arguments
    Control_socket := flag.String("control", "127.0.0.1:7722", "control UDP socket")
    External_ip := flag.String("external", "", "external address")
    Internal_ip := flag.String("internal", "", "internal address")
    Debug_arg := flag.Bool("debug", false, "debug output")
    flag.Parse()

//not working without external IP
    if ( *External_ip == "" ) {
        flag.PrintDefaults()
        return
    }

//without Internal IP use External
    if ( *Internal_ip == "" ) {
        *Internal_ip = *External_ip
    }
    Debug = *Debug_arg

//start listening control socket
    ServerAddr, err := net.ResolveUDPAddr("udp", *Control_socket)
    if err != nil {
        fmt.Println("Control socket: ", err)
        return
    }
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    if err != nil {
        fmt.Println("Control socket: ", err)
        return
    }
    defer ServerConn.Close()

//chan for commands from kamailio with buffer 100 commands
    Command_chan := make(chan ProxyCommand, 100)

//chan for response to kamailio with buffer 100 commands
    Response_chan := make(chan ProxyCommand, 100)

//start comand handler
    go ReqHandler (Command_chan, Response_chan, *External_ip, *Internal_ip)
    go RespSender (Response_chan, ServerConn)
//read commands from kamailio with buffer 1024 bytes
    Buffer := make([]byte, 1024)
    for {
        Nbytes, Raddr, err := ServerConn.ReadFromUDP(Buffer)
        if err != nil {
            fmt.Println("Read from control socket: ", err)
            break
        }
        var Recv ProxyCommand
        Recv.command = string(Buffer[0:Nbytes])
        Recv.raddr = Raddr
        Command_chan <- Recv
        if Debug {
            fmt.Println("Received: ",Recv.command, " from ",Recv.raddr)
        }
    }
}
