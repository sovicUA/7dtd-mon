package sq7dtd

import (
	"bytes"
	"log"
	"fmt"
	"encoding/binary"
	"net"
	"strconv"
	"errors"
	"time"
	"encoding/json"
)


type A2S_INFO struct {
	Key			string
	Host		string
	Port		int
	//Header	int32
	ReturnCode	byte
	Protocol	byte
	ServerName	string
	World		string
	DescShort	string
	DescLong	string
	Players     byte
	PlayersMAX  byte
	Bots		byte
	ServerType  byte
	Environment byte
	Visibility  byte
	Version     string
}

var server A2S_INFO

func ServerType(t byte) string {
	switch _t := string(t); _t {
	case "d":
		return "Dedicated server"
	case "l":
		return "Non-dedicated server"
	default:
		return "Unknown server type"
	}
}

func Environment(e byte) string {
	switch _e := string(e); _e {
	case "l":
		return "Linux"
	case "w":
		return "Windows"
	case "m":
		return "Mac"
	default:
		return "Unknown operation system"
	}
}

func Protected(p byte) string {
	if p == 1 {
		return "Private"
	} else {
		return "Public"
	}
}

func B2S(b []byte, i int) (string, int) {
	c := bytes.IndexByte(b[i:], 0) + i
	return string(b[i:c]), c + 1
}

func Parse(b []byte) {
	i := 6 // Index of data server name
	// Header
	// binary.Read(bytes.NewReader(b[:4]), binary.BigEndian, &server.Header)
	// ReturnCode | Always equal to 'I' (0x49) 
	server.ReturnCode = b[4]
	// Protocol
	server.Protocol = b[5]
	// Server Name
	server.ServerName, i = B2S(b, i)
	// Map
	server.World, i = B2S(b, i)
	// Game Short Description
	server.DescShort, i = B2S(b, i)
	// Game Long Description
	server.DescLong, i = B2S(b, i)
	i++; i++ // +2 byte for next data
	// Online Players
	server.Players = b[i]; i++
	// MAX Players
	server.PlayersMAX = b[i]; i++
	// Bots
	server.Bots = b[i]; i++
	// Server Type
	server.ServerType = b[i]; i++
	// Server Environment
	server.Environment = b[i]; i++
	// Protected
	server.Visibility = b[i]; i++
	i++
	// Version game
	server.Version, i = B2S(b, i)
}

func Query(host string, port int) error {
	// Server Query "A2S_INFO" message
	message := []byte("\xFF\xFF\xFF\xFFTSource Engine Query\x00")
	
	server.Host = host; server.Port = port	
	service := host + ":" + strconv.Itoa(port)
	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Timeout: 3 sec	
	deadline := time.Now().Add(3 * time.Second)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		log.Fatal(err)
		return errors.New("Failed to set timeout")
	}

	defer conn.Close()
	/***
	log.Printf("Established connection to %s \n", service)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())
	log.Printf("Send A2S_INFO message... \n") 
	***/
	_, err = conn.Write(message)
	if err != nil {
		log.Fatal(err)
		return errors.New("Failed to send message")
	}

	buffer := make([]byte, 1400)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return errors.New("Failed read data")
	}

	var Header int32
	binary.Read(bytes.NewReader(buffer[:4]), binary.BigEndian, &Header)
	if Header == -1 && buffer[4] == 73 {
		Parse(buffer)
		return nil
	}
	
	log.Println("UDP Server : ", addr)
	log.Println("Received from UDP server : ", string(buffer[:n]))

	for i := 0; i < len(buffer); i++ {
		fmt.Printf("%x ", buffer[i])
	}
	
	return errors.New("Failed to parse data")
}

func Json() string {
	serverJson, err :=  json.Marshal(server)
	if err != nil {
        log.Println(err)
		return "{\"ReturnCode\":-1}"
    }
    return string(serverJson)
}

func String() string {
	return `
	Key: ` + server.Key + `
	Host: ` + server.Host + `
	Port: ` + strconv.Itoa(server.Port) + `
	Return Code: ` + strconv.Itoa(int(server.ReturnCode)) + `
	Protocol: ` + strconv.Itoa(int(server.Protocol)) + `	
	Server Name: ` + server.ServerName + `	
	World: ` + server.World + `		
	Game Short Description: ` + server.DescShort + `	
	Game Long Description: ` + server.DescLong + `	
	Players Online: ` + strconv.Itoa(int(server.Players)) + `     
	Limit Players: ` + strconv.Itoa(int(server.PlayersMAX)) + ` 
	Bots: ` + strconv.Itoa(int(server.Bots)) + `		
	Server Type: ` + ServerType(server.ServerType) + `  
	Environment: ` + Environment(server.Environment) + `
	Visibility: ` + Protected(server.Visibility) + ` 
	Version: ` + server.Version + `
	`
}