package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const host = "us.pool.ntp.org:123"
const ntpEpochOffset = 2208988800

type packet struct {
	Settings       uint8  // leap yr indicator, ver number, and mode
	Stratum        uint8  // stratum of local clock
	Poll           int8   // poll exponent
	Precision      int8   // precision exponent
	RootDelay      uint32 // root delay
	RootDispersion uint32 // root dispersion
	ReferenceID    uint32 // reference id
	RefTimeSec     uint32 // reference timestamp sec
	RefTimeFrac    uint32 // reference timestamp fractional
	OrigTimeSec    uint32 // origin time secs
	OrigTimeFrac   uint32 // origin time fractional
	RxTimeSec      uint32 // receive time secs
	RxTimeFrac     uint32 // receive time frac
	TxTimeSec      uint32 // transmit time secs
	TxTimeFrac     uint32 // transmit time frac
}


func main() {
	fmt.Println("Hello, NTP!")

	conn, err := net.Dial("udp", host)
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(
		time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal("failed to set deadline: ", err)
	}

	req := &packet{Settings: 0x1B}

	if err := binary.Write(conn, binary.BigEndian, req); err != nil {
		log.Fatalf("failed to send request: %v", err)
	}

	rsp := &packet{}
	if err := binary.Read(conn, binary.BigEndian, rsp); err != nil {
		log.Fatalf("failed to read server response: %v", err)
	}

	secs := float64(rsp.TxTimeSec) - ntpEpochOffset

	fmt.Printf("Current time: %v\n", time.Unix(int64(secs), 0))
}
