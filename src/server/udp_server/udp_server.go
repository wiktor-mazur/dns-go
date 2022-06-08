package udp_server

import (
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"github.com/wiktor-mazur/dns-go/src/protocol"
	"github.com/wiktor-mazur/dns-go/src/resolver"
	"log"
	"net"
	"strconv"
	"sync"
)

type Config struct {
	ListenIP   net.IP
	ListenPort int
}

type UDPServer struct {
	cfg      Config
	resolver *resolver.Resolver
	conn     *net.UDPConn
}

func New(cfg Config, resolver *resolver.Resolver) *UDPServer {
	return &UDPServer{cfg: cfg, resolver: resolver}
}

func (v *UDPServer) Start(wg *sync.WaitGroup) error {
	defer wg.Done()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: v.cfg.ListenPort,
		IP:   v.cfg.ListenIP,
	})
	if err != nil {
		return err
	}

	v.conn = conn

	log.Printf("UDP server listening at %s\n", conn.LocalAddr().String())

	for {
		// @TODO: Implement timeout to reading from UDP socket
		rawQueryPacket := make([]byte, buffer.DNSBufferSize+1)
		packetLength, clientAddr, err := conn.ReadFromUDP(rawQueryPacket)

		if packetLength > buffer.DNSBufferSize {
			// @TODO: Respond to client with FORMERR when that happens
			log.Printf("Packet received from [%s] is too large (%d/%d bytes).", clientAddr, packetLength, buffer.DNSBufferSize)
			continue
		}

		if err != nil {
			// @TODO: Respond to client with FORMERR when that happens
			log.Printf("Could not read packet received from [%s]: %s", clientAddr, err.Error())
			continue
		}

		go v.handleRequest(clientAddr, rawQueryPacket[:packetLength])
	}
}

func (v *UDPServer) handleRequest(clientAddr *net.UDPAddr, rawPacket []byte) {
	queryPacket, err := protocol.DnsPacketFromRawBuffer(rawPacket)
	if err != nil {
		logFormat := "[%s] Received query from [%s] but could not parse the packet: %s"

		if queryPacket != nil {
			log.Printf(logFormat, strconv.Itoa(int(queryPacket.Header.ID)), clientAddr, err.Error())

			response, _ := v.resolver.QueryToErrResponse(queryPacket, common.FORMERR).ToRawBuffer()

			go v.sendResponse(clientAddr, queryPacket.Header.ID, response)
		} else {
			log.Printf(logFormat, "?", clientAddr, err.Error())
		}

		return
	}

	log.Printf("[%d] Received query from [%s]: %s", queryPacket.Header.ID, clientAddr, queryPacket.Questions[0].CompactString())

	responsePacket, err := v.resolver.ResolveQuery(queryPacket)
	if err != nil {
		log.Printf("[%d] Could not perform the lookup: %s", queryPacket.Header.ID, err.Error())

		response := v.resolver.QueryToErrResponse(queryPacket, common.SERVFAIL)
		respBuf, _ := response.ToRawBuffer()

		go v.sendResponse(clientAddr, queryPacket.Header.ID, respBuf)

		return
	}

	responsePacketBuf, err := responsePacket.ToRawBuffer()
	if err != nil {
		log.Printf("[%d] Could not serialize response packet: %s", queryPacket.Header.ID, err.Error())

		resp := v.resolver.QueryToErrResponse(queryPacket, common.SERVFAIL)
		respBuf, _ := resp.ToRawBuffer()

		go v.sendResponse(clientAddr, queryPacket.Header.ID, respBuf)

		return
	}

	go v.sendResponse(clientAddr, responsePacket.Header.ID, responsePacketBuf)
}

func (v *UDPServer) sendResponse(addr *net.UDPAddr, ID uint16, data []byte) {
	_, err := v.conn.WriteToUDP(data, addr)
	if err != nil {
		log.Printf("[%d] Couldn't send response %v", ID, err)
	} else {
		log.Printf("[%d] Response sent to client", ID)
	}
}
