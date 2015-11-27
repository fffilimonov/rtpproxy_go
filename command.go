package main

import (
    "fmt"
    "net"
    "regexp"
    "strings"
)

func ReqHandler (Command_chan chan ProxyCommand, Response_chan chan ProxyCommand, External_ip string, Internal_ip string) {
    var V = regexp.MustCompile(`^V$`)
    var VF = regexp.MustCompile(`^VF$`)
    var D = regexp.MustCompile(`^D$`)
    var Uc = regexp.MustCompile(`^Uc.*$`)
    var UEIc = regexp.MustCompile(`^UEIc.*$`)
    var UIEc = regexp.MustCompile(`^UIEc.*$`)
    var Lc = regexp.MustCompile(`^Lc.*$`)
    var LEIc = regexp.MustCompile(`^LEIc.*$`)
    var LIEc = regexp.MustCompile(`^LIEc.*$`)
    var P1 = regexp.MustCompile(`^P-1$`)
    Proxy_set := make(map[string]ProxyControl)
    for {
        select {
            case Command := <-Command_chan:
                Splited := strings.Split(Command.command, " ")
                var Answer string
                switch {
                    case V.MatchString(Splited[1]):
                        Answer="20040107"
                    case VF.MatchString(Splited[1]):
                        Answer="1"
                    case D.MatchString(Splited[1]):
                        ID:=Splited[2]
                        tmp:=Proxy_set[ID]
                        if (tmp.Outgoing!=nil) {
                            tmp.Incoming <- "quit"
                            Answer = <- tmp.Outgoing
                            delete(Proxy_set,ID)
                            fmt.Println("In map: ",len(Proxy_set))
                        } else {
                            Answer="1"
                        }
                    case P1.MatchString(Splited[1]):
                        ID:=Splited[2]
                        tmp:=Proxy_set[ID]
                        if (tmp.Outgoing!=nil) {
                                tmp.Incoming <- Splited[3]
                        }
                        Answer="1"
                    case Uc.MatchString(Splited[1])||UEIc.MatchString(Splited[1])||UIEc.MatchString(Splited[1]):
                        ID:=Splited[2]
                        tmp:=Proxy_set[ID]
                        if (tmp.Incoming==nil) {
                            tmp = ProxyControl{make(chan string),make(chan string),"","","","","",""}
                            go proxyProc(ID, Proxy_set, External_ip, Internal_ip)
                        }
                        tmp.InviteAddr=Splited[3]
                        tmp.InvitePort=Splited[4]
                        Proxy_set[ID]=tmp
                        if(Uc.MatchString(Splited[1])) {
                            tmp.Incoming <- "Uc"
                        }
                        if(UEIc.MatchString(Splited[1])) {
                            tmp.Incoming <- "UEIc"
                        }
                        if(UIEc.MatchString(Splited[1])) {
                            tmp.Incoming <- "UIEc"
                        }
                        Answer = <- tmp.Outgoing
                    case Lc.MatchString(Splited[1])||LEIc.MatchString(Splited[1])||LIEc.MatchString(Splited[1]):
                        ID:=Splited[2]
                        tmp:=Proxy_set[ID]
                        if (tmp.Incoming!=nil) {
                            tmp.OkAddr=Splited[3]
                            tmp.OkPort=Splited[4]
                            Proxy_set[ID]=tmp
                            if(Lc.MatchString(Splited[1])) {
                                tmp.Incoming <- "Lc"
                            }
                            if(LEIc.MatchString(Splited[1])) {
                                tmp.Incoming <- "LEIc"
                            }
                            if(LIEc.MatchString(Splited[1])) {
                                tmp.Incoming <- "LIEc"
                            }
                            Answer = <- tmp.Outgoing
                        } else {
                            Answer="1"
                        }
                    default:
                        Answer="0"
                }
                Response_chan <- ProxyCommand{ConcatStr(" ",Splited[0],Answer),Command.raddr}
        }
    }
}

func RespSender (Response_chan chan ProxyCommand, ServerConn *net.UDPConn) {
    for {
        select {
            case Responce := <- Response_chan:
                _,err := ServerConn.WriteToUDP([]byte(Responce.command), Responce.raddr)
                if err != nil {
                    fmt.Print("Couldn't send response ", err)
                } else {
                    fmt.Println("Sent: ",Responce.command, "to ",Responce.raddr)
                }
        }
    }
}
