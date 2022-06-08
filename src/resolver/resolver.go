package resolver

import (
	"fmt"
	"github.com/wiktor-mazur/dns-go/src/buffer"
	"github.com/wiktor-mazur/dns-go/src/common"
	"github.com/wiktor-mazur/dns-go/src/protocol"
	"log"
	"net"
)

type Config struct {
	InternetRootServer string
}

type Resolver struct {
	cfg Config
}

func New(cfg Config) *Resolver {
	return &Resolver{cfg: cfg}
}

func (v *Resolver) ResolveQuery(query *protocol.DnsPacket) (*protocol.DnsPacket, error) {
	responsePacket := protocol.NewDnsPacket()
	responsePacket.Header.ID = query.Header.ID
	responsePacket.Header.IsResponse = true
	responsePacket.Header.RecursionDesired = true
	responsePacket.Header.RecursionAvailable = true

	if len(query.Questions) == 0 || query.Header.QuestionsCount == 0 {
		responsePacket.Header.ResultCode = common.FORMERR
		return responsePacket, nil
	}

	if query.Header.IsResponse {
		responsePacket.Header.ResultCode = common.FORMERR
		return responsePacket, nil
	}

	question := query.Questions[0]

	lookup, err := v.LookupRecursive(query.Header.ID, question.Name, question.QueryType)
	if err != nil {
		return nil, err
	}

	responsePacket.Header.ResultCode = common.NOERROR
	responsePacket.AddQuestion(question)

	for _, v := range lookup.Answers {
		log.Printf("[%d] Answer %s", query.Header.ID, v.CompactString())
		responsePacket.AddAnswer(v)
	}

	for _, v := range lookup.Authorities {
		log.Printf("[%d] Authority %s", query.Header.ID, v.CompactString())
		responsePacket.AddAuthority(v)
	}

	for _, v := range lookup.Resources {
		log.Printf("[%d] Resource %s", query.Header.ID, v.CompactString())
		responsePacket.AddResource(v)
	}

	return responsePacket, nil
}

func (v *Resolver) QueryToErrResponse(query *protocol.DnsPacket, err common.ResultCode) *protocol.DnsPacket {
	response := *query

	response.Header.IsResponse = true
	response.Header.RecursionAvailable = true
	response.Header.ResultCode = err

	return &response
}

func (v *Resolver) Lookup(qName string, qType common.QueryType, serverAddr *net.UDPAddr) (*protocol.DnsPacket, error) {
	queryPacket := buildQueryPacket(qName, qType)

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	queryBuf, err := queryPacket.ToBuffer()
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(queryBuf.GetBytes())
	if err != nil {
		return nil, err
	}

	buf := make([]byte, buffer.DNSBufferSize)
	responseLength, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}

	responsePacket, err := protocol.DnsPacketFromRawBuffer(buf[:responseLength])
	if err != nil {
		return nil, err
	}

	return responsePacket, nil
}

func (v *Resolver) LookupRecursive(queryID uint16, qName string, qType common.QueryType) (*protocol.DnsPacket, error) {
	ns := v.cfg.InternetRootServer

	// @TODO: check max iterations and max recursive depth
	for {
		log.Printf("[%d] Attempting lookup of %s %s with ns %s", queryID, qType.String(), qName, ns)

		serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ns, 53))
		if err != nil {
			return nil, err
		}

		response, err := v.Lookup(qName, qType, serverAddr)
		if err != nil {
			return nil, err
		}

		isFinalResultFound := len(response.Answers) > 0 && response.Header.ResultCode == common.NOERROR
		if isFinalResultFound {
			return response, nil
		}

		nameNotExists := response.Header.ResultCode == common.NXDOMAIN
		if nameNotExists {
			return response, nil
		}

		// DNS server included name server's A record in the additional section, so we already have the IP
		resolvedNS := response.GetResolvedNS(qName)
		if resolvedNS != nil {
			newNsIP := resolvedNS.GetIP()
			ns = newNsIP.String()

			continue
		}

		// DNS server didn't include name server's IP, so we need to recursively find it
		unresolvedNS := response.GetUnresolvedNS(qName)
		if unresolvedNS == nil {
			// we didn't get any name servers, so we return what we have
			return response, nil
		}

		recursiveResponse, err := v.LookupRecursive(queryID, unresolvedNS.GetHost(), common.A)
		if err != nil {
			return nil, err
		}

		newNS := recursiveResponse.GetFirstARecord()

		if newNS != nil {
			// continue search with next NS
			newNsIP := newNS.GetIP()
			ns = newNsIP.String()
		} else {
			// there are no A records, we return what we have
			return response, nil
		}
	}
}

func buildQueryPacket(qName string, qType common.QueryType) *protocol.DnsPacket {
	result := protocol.NewDnsPacket()
	result.Header.RecursionDesired = true

	question := protocol.NewDnsQuestion()
	question.Name = qName
	question.QueryType = qType

	result.AddQuestion(*question)

	return result
}
