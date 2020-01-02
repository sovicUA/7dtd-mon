package sq7dtd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type A2S_RULES struct {
	Key         string
	Host        string
	Port        int
	ReturnCode  byte
	RuleNumbers uint16
	RuleList    map[string]string
}

var rules A2S_RULES

func getCurrentTime() (int, int, int, error) {
	current, err := strconv.Atoi(rules.RuleList["CurrentServerTime"])
	if err != nil {
		log.Fatal(err)
		return 0, 0, 0, errors.New("Failed to convert current server time")
	}

	day := int((current / 24000) + 1)
	hour := int((current % 24000) / 1000)
	minutes := int(((current % 1000) * 60) / 1000)

	return day, hour, minutes, nil
}

func getRuleNumbers() int {
	return int(rules.RuleNumbers)
}

func ParseRules(b []byte) {
	var Name, Value string
	seek := 7 // Index of data server rules
	// Header
	// binary.Read(bytes.NewReader(b[:4]), binary.BigEndian, &server.Header)
	// ReturnCode | Always equal to 'E' (0x45)
	rules.ReturnCode = b[4]
	rules.RuleNumbers = binary.LittleEndian.Uint16(b[5:7])

	rules.RuleList = make(map[string]string)
	for i := 0; i < int(rules.RuleNumbers); i++ {
		Name, seek = B2S(b, seek)
		Value, seek = B2S(b, seek)
		rules.RuleList[Name] = Value
	}
}

func QueryRules(host string, port int) error {
	// Server Query "A2S_RULES" challenge number
	message := []byte("\xFF\xFF\xFF\xFFV\xFF\xFF\xFF\xFF")

	rules.Host = host
	rules.Port = port
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
	if Header == -1 && buffer[4] == 69 {
		ParseRules(buffer)
		return nil
	}

	log.Println("UDP Server : ", addr)
	log.Println("Received from UDP server : ", string(buffer[:n]))

	for i := 0; i < len(buffer); i++ {
		fmt.Printf("%x ", buffer[i])
	}

	return errors.New("Failed to parse data")
}
