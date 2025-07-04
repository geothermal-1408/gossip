package main

import (
	"log"
	"net"
    "strings"
    "time"
)

const port = "6969"
const ratelimit = 1
const bannedLimit = 60*10.0

type MessageType int
const
(
    ClientConnected MessageType = iota+1
    ClientDisconnected
    NewMessage
)

type Message struct {
    Type MessageType
    Conn net.Conn
    Text string
}

type Client struct {
    Conn net.Conn
    LastMessage time.Time
    FlagCount int
}

func server(messages chan Message) {
    clients := map[string]*Client{}
    banned_user := map[string]time.Time{}
    
    for {
        msg := <- messages
        switch msg.Type {
        case ClientConnected:
            addr := msg.Conn.RemoteAddr().(*net.TCPAddr)
            bannedAt,banned := banned_user[addr.IP.String()]
            
            if banned {
                if time.Now().Sub(bannedAt).Seconds() >= bannedLimit {
                    delete(banned_user,addr.IP.String())
                    banned = false
                }
            }

            if !banned {
                clients[msg.Conn.RemoteAddr().String()]= &Client{
                    Conn: msg.Conn,
                    LastMessage: time.Now(),
                }
                log.Printf("\x1b[32mINFO\x1b[0m: connected to server: %s\n",msg.Conn.RemoteAddr().String())
            }else {
                msg.Conn.Close()
            }
            
        case ClientDisconnected:
            
            delete(clients,msg.Conn.RemoteAddr().String())
            log.Printf("\x1b[32mINFO\x1b[0m: disconnected from server: %s\n",msg.Conn.RemoteAddr().String())
            
        case NewMessage:
            authorAddr := msg.Conn.RemoteAddr().(*net.TCPAddr)
            now := time.Now()
            author := clients[authorAddr.String()]
            if time.Now().Sub(author.LastMessage).Seconds() >= ratelimit {
                author.LastMessage = now
                author.FlagCount = 0
                log.Printf("%s %s: \n",msg.Conn.RemoteAddr().String(),msg.Text)
                for _, client := range clients{
                    if client.Conn.RemoteAddr().String() != authorAddr.String() {
                        _,err := client.Conn.Write([]byte(msg.Text))
                        if err != nil {
                            log.Printf("\x1b[31mERROR\x1b[0m: can't send data from server to %s: %s\n",client.Conn.RemoteAddr(),err)
                        }
                    }
                }
            }else{
                author.FlagCount+=1
                if author.FlagCount >= 3 {
                    banned_user[authorAddr.IP.String()] = now
                    author.Conn.Close()
                }
            }
        }
    }
}

func client(conn net.Conn, messages chan Message){
	buffer := make([]byte,64)
    for {
        n,err := conn.Read(buffer);
        if err != nil {
            conn.Close()
            messages <- Message{
                Type: ClientDisconnected,
                Conn: conn,
            }
            return
        }
        text := string(buffer[0:n])
        //exit command 
        if strings.TrimSpace(text) == ":quit" {
            conn.Close()
            messages <- Message{
                Type: ClientDisconnected,
                Conn: conn,
            }
            return
        }
        messages <- Message{
            Type: NewMessage,
            Text: text,
            Conn: conn,           
        }
    }
}

func main() {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("\x1b[31mERROR\x1b[0m: could not listen to port %s: %s\n", port, err)
	}
	log.Printf("\x1b[32mINFO\x1b[0m: listening on %s\n", port)

    messages := make(chan Message)
    go server(messages)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("\x1b[31mERROR\x1b[0m: could not accept new connection: %s\n", err)
		}

        messages <- Message{
            Type:  ClientConnected,
            Conn: conn,
        }
		go client(conn,messages)
	}
}
