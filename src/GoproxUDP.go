package main

import (
	"fmt"
	"net"
	"runtime"
)

var (
	AddrFromClient = "localhost:80"
	TARGSRV        = "192.168.1.1:80"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fromClient, err := net.ResolveUDPAddr("udp", AddrFromClient)
	if err != nil {
		fmt.Println("Error > ", err)
		return
	}
	conClientReq, err := net.ListenUDP("udp", fromClient)
	if err != nil {
		fmt.Println("Error > ", err)
	}
	proxyToSRV, err := net.ResolveUDPAddr("udp", TARGSRV)
	if err != nil {
		fmt.Println("Error > ", err)
	}
	conProxyToSRV, err := net.DialUDP("udp", nil, proxyToSRV)
	if err != nil {
		fmt.Println("Error > ", err)
	}
	defer conClientReq.Close()
	defer conProxyToSRV.Close()
	for {
		receivedRequest := make([]byte, 1024)
		length, fromClient, err := conClientReq.ReadFromUDP(receivedRequest)
		if err != nil {
			fmt.Println("Error > ", err)
		}
		go sendToSRV(receivedRequest[:length], fromClient, conClientReq, conProxyToSRV)
	}
}

func sendToSRV(buf []byte, addrOfClient *net.UDPAddr, conClientReq *net.UDPConn, conProxyToSRV *net.UDPConn) {
	_, err := conProxyToSRV.Write([]byte(buf))
	if err != nil {
		fmt.Println("Error > ", err)
	}
	responseToClient(addrOfClient, conProxyToSRV, conClientReq)
}

func responseToClient(addrOfClient *net.UDPAddr, conProxyToSRV *net.UDPConn, conClientReq *net.UDPConn) {
	response := make([]byte, 1024)
	length, _, err := conProxyToSRV.ReadFromUDP(response)
	if err != nil {
		fmt.Println("Error > ", err)
		return
	}
	_, err = conClientReq.WriteToUDP(response[:length], addrOfClient)
	if err != nil {
		fmt.Println("Error > ", err)
		return
	}
}
