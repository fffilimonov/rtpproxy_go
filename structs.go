package main

import (
    "net"
)

//struct for proxy instance
type ProxyControl struct {
    Incoming chan string //chan for recieving commands
    Outgoing chan string //chan for sending answers
    Inport string //local port for calling party
    Outport string //local port for called party
    InviteAddr string //address from INVITE
    InvitePort string //port from INVITE
    OkAddr string //address from 200OK
    OkPort string //port from 2000OK
}

//struct for recieved command and source
type ProxyCommand struct {
    command string //command
    raddr *net.UDPAddr //source of command
}
