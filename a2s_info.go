package 7dtd-mon

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
	"encoding/json"
)


type A2S_INFO struct {
	//Header		int32
	ReturnCode	byte
	Protocol	byte
	ServerName	string
	World		string
	DescShort	string
	DescLong	string
	Players     byte
	PlayersMAX  byte
	Bots		byte
	ServerType  string
	Environment string
	Visibility  string
	Version     string
}

var server A2S_INFO

func getServerType(t byte) string {
	switch _t := string(t); _t {
	case "d":
		return "Dedicated server"
	case "l":
		return "Non-dedicated server"
	default:
		return "Unknown server type"
	}
}

func getEnvironment(e byte) string {
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

func getProtected(p byte) string {
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

func parse(b []byte) {
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
	server.ServerType = getServerType(b[i]); i++
	// Server Environment
	server.Environment = getEnvironment(b[i]); i++
	// Protected
	server.Visibility = getProtected(b[i]); i++
	i++
	// Version game
	server.Version, i = B2S(b, i)
}

