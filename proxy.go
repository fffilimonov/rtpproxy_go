package main

import (
    "fmt"
    "net"
    "strings"
    "sync"
)

func proxyProc(ID string, proxy_set map[string]ProxyControl, External_ip string, Internal_ip string) {
    Sin := make(chan *net.UDPAddr,10)
    Sout := make(chan *net.UDPAddr,10)
    Sinq := make(chan bool,10)
    Soutq := make(chan bool,10)
    CNin := make(chan string,10)
    CNout := make(chan string,10)
    Datain := make(chan []byte,1024)
    Dataout := make(chan []byte,1024)
    var wg sync.WaitGroup
    ctmp := proxy_set[ID]
    var SIDE bool //true if INVITE (UEIc) from local, false if INVITE (UIEc) from external
loop:
    for {
        select {
            case in := <- ctmp.Incoming:
fmt.Println("Got ", in)
                tmp := proxy_set[ID]
                if (in=="Uc") {
                    if(tmp.Inport=="") {
                        ServerAddrin,_ := net.ResolveUDPAddr("udp",ConcatStr (":",External_ip,"0"))
                        ServerConnin, err := net.ListenUDP("udp", ServerAddrin)
                        if err != nil {
                            fmt.Println("Error: ",err)
                        }
                        tmpinport:=ServerConnin.LocalAddr()
                        inport:=strings.Split(tmpinport.String(),":")
                        tmp.Inport=inport[1]
                        go proxyStreamReader(ServerConnin,Sin,Datain,wg)
                        go proxyStreamSender(ServerConnin,Sin,Sinq,Dataout,wg,CNin)
                        remote_addr,_:=net.ResolveUDPAddr("udp",ConcatStr (":",tmp.InviteAddr,tmp.InvitePort))
                        Sout<-remote_addr
                    }
                    tmp.Outgoing<-tmp.Inport
                }
                if (in=="UEIc") {
                    if(tmp.Inport=="") {
                        ServerAddrin,_ := net.ResolveUDPAddr("udp",ConcatStr (":",External_ip,"0"))
                        ServerConnin, err := net.ListenUDP("udp", ServerAddrin)
                        if err != nil {
                            fmt.Println("Error: ",err)
                        }
                        tmpinport:=ServerConnin.LocalAddr()
                        inport:=strings.Split(tmpinport.String(),":")
                        tmp.Inport=inport[1]
                        go proxyStreamReader(ServerConnin,Sin,Datain,wg)
                        go proxyStreamSender(ServerConnin,Sin,Sinq,Dataout,wg,CNin)
                        remote_addr,_:=net.ResolveUDPAddr("udp",ConcatStr (":",tmp.InviteAddr,tmp.InvitePort))
                        Sout<-remote_addr
                        CNin<-"true"
                    }
                    tmp.Outgoing<-ConcatStr (" ",tmp.Inport,External_ip)
                    SIDE=true
                }
                if (in=="UIEc") {
                    if(tmp.Inport=="") {
                        ServerAddrin,_ := net.ResolveUDPAddr("udp",ConcatStr (":",Internal_ip,"0"))
                        ServerConnin, err := net.ListenUDP("udp", ServerAddrin)
                        if err != nil {
                            fmt.Println("Error: ",err)
                        }
                        tmpinport:=ServerConnin.LocalAddr()
                        inport:=strings.Split(tmpinport.String(),":")
                        tmp.Inport=inport[1]
                        go proxyStreamReader(ServerConnin,Sin,Datain,wg)
                        go proxyStreamSender(ServerConnin,Sin,Sinq,Dataout,wg,CNin)
                        remote_addr,_:=net.ResolveUDPAddr("udp",ConcatStr (":",tmp.InviteAddr,tmp.InvitePort))
                        Sout<-remote_addr
                        CNin<-"false"
                    }
                    tmp.Outgoing<-ConcatStr (" ",tmp.Inport,Internal_ip)
                    SIDE=false
                }
                if (in=="Lc") {
                    if(tmp.Outport=="") {
                        ServerAddrout,_ := net.ResolveUDPAddr("udp",ConcatStr (":",External_ip,"0"))
                        ServerConnout, err := net.ListenUDP("udp", ServerAddrout)
                        if err != nil {
                            fmt.Println("Error: ",err)
                        }
                        tmpoutport:=ServerConnout.LocalAddr()
                        outport:=strings.Split(tmpoutport.String(),":")
                        tmp.Outport=outport[1]
                        go proxyStreamReader(ServerConnout,Sout,Dataout,wg)
                        go proxyStreamSender(ServerConnout,Sout,Soutq,Datain,wg,CNout)
                        remote_addr,_:=net.ResolveUDPAddr("udp",ConcatStr (":",tmp.OkAddr,tmp.OkPort))
                        Sin<-remote_addr
                    }
                    tmp.Outgoing<-tmp.Outport
                }
                if (in=="LEIc") {
                    if(tmp.Outport=="") {
                        ServerAddrout,_ := net.ResolveUDPAddr("udp",ConcatStr (":",Internal_ip,"0"))
                        ServerConnout, err := net.ListenUDP("udp", ServerAddrout)
                        if err != nil {
                            fmt.Println("Error: ",err)
                        }
                        tmpoutport:=ServerConnout.LocalAddr()
                        outport:=strings.Split(tmpoutport.String(),":")
                        tmp.Outport=outport[1]
                        go proxyStreamReader(ServerConnout,Sout,Dataout,wg)
                        go proxyStreamSender(ServerConnout,Sout,Soutq,Datain,wg,CNout)
                        remote_addr,_:=net.ResolveUDPAddr("udp",ConcatStr (":",tmp.OkAddr,tmp.OkPort))
                        Sin<-remote_addr
                        CNout<-"false"
                    }
                    tmp.Outgoing<-ConcatStr (" ",tmp.Outport,Internal_ip)
                }
                if (in=="LIEc") {
                    if(tmp.Outport=="") {
                        ServerAddrout,_ := net.ResolveUDPAddr("udp",ConcatStr (":",External_ip,"0"))
                        ServerConnout, err := net.ListenUDP("udp", ServerAddrout)
                        if err != nil {
                            fmt.Println("Error: ",err)
                        }
                        tmpoutport:=ServerConnout.LocalAddr()
                        outport:=strings.Split(tmpoutport.String(),":")
                        tmp.Outport=outport[1]
                        go proxyStreamReader(ServerConnout,Sout,Dataout,wg)
                        go proxyStreamSender(ServerConnout,Sout,Soutq,Datain,wg,CNout)
                        remote_addr,_:=net.ResolveUDPAddr("udp",ConcatStr (":",tmp.OkAddr,tmp.OkPort))
                        Sin<-remote_addr
                        CNout<-"true"
                    }
                    tmp.Outgoing<-ConcatStr (" ",tmp.Outport,External_ip)
                }
                if (in=="upload") {
                    if (SIDE) {
                        CNin<-"true"
                        CNout<-"false"
                        fmt.Println("got upload int")
                    } else {
                        CNin<-"false"
                        CNout<-"true"
                        fmt.Println("got upload ext")
                    }
                }
                if (in=="download") {
                    if (SIDE) {
                        CNin<-"false"
                        CNout<-"no"
                        fmt.Println("got download int")
                    } else {
                        CNin<-"no"
                        CNout<-"false"
                        fmt.Println("got download ext")
                    }
                }
                if (in=="quit") {
                    Sinq<-true
                    Soutq<-true
                    break loop
                }
                proxy_set[ID] = tmp
        }
    }
    wg.Wait()
    fmt.Println("Stop proxyProc")
    ctmp.Outgoing<-"1"
}

func proxyStreamReader(conn *net.UDPConn,rch chan *net.UDPAddr, datach chan []byte, wg sync.WaitGroup) {
    wg.Add(1)
    defer wg.Done()
    buf := make([]byte, 1024)
    var raddr *net.UDPAddr
    fmt.Println("Start reader", conn.LocalAddr())
    for {
        n,addr,err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error: ",err)
            break
        }
        if (raddr==nil) {
            raddr=addr
            fmt.Println("Start recieving reader ",conn.LocalAddr(),raddr)
            rch<-raddr
        }
        datach<-buf[0:n]
    }
    fmt.Println("Reader stop", raddr)
}

func proxyStreamSender(conn *net.UDPConn, rch chan *net.UDPAddr, qch chan bool, datach chan []byte, wg sync.WaitGroup, controlCN chan string) {
    wg.Add(1)
    defer wg.Done()
    var raddr *net.UDPAddr
    var i int = 0
    var cCN string
loop:
    for {
        select {
            case cCN = <-controlCN:
                fmt.Println("CN sender ",conn.LocalAddr(),raddr,cCN)
            case raddr = <-rch:
                fmt.Println("Start sender ",conn.LocalAddr(),raddr)
            case data := <-datach:
                if (cCN=="true") {
                    CN := make([]byte, 17)
                    CN[0] = 0x80
                    CN[1] = 0x0d
                    CN[3] = byte(i)
                    CN[12] = 0x72
                    CN[13] = 0x6c
                    CN[14] = 0x72
                    CN[15] = 0x77
                    CN[16] = 0x68
                    _,err := conn.WriteToUDP(CN, raddr)
                    i++
                    if err != nil {
                        fmt.Print("Couldn't send response ", err)
                    }
                }
                if (cCN=="false") {
                    _,err := conn.WriteToUDP(data, raddr)
                    if err != nil {
                        fmt.Print("Couldn't send response ", err)
                    }
                }
            case <-qch:
                conn.Close()
                break loop
        }
    }
    fmt.Println("Sender stop ",raddr)
}
