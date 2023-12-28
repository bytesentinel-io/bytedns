package resolver

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

type DnsForwarding struct {
	Enabled bool   `json:"enabled"`
	Server  string `json:"server"`
}

type DnsRegistry struct {
	Domains     []DnsDomain     `json:"domains"`
	Forwarding  DnsForwarding   `json:"forwarding"`
	RootServers []DnsRootServer `json:"rootServers"`
}

type DnsDomain struct {
	Name  string      `json:"name"`
	A     []DnsRecord `json:"a"`
	AAAA  []DnsRecord `json:"aaaa"`
	CNAME []DnsRecord `json:"cname"`
	MX    []DnsRecord `json:"mx"`
	NS    []DnsRecord `json:"ns"`
	PTR   []DnsRecord `json:"ptr"`
	SOA   []DnsRecord `json:"soa"`
	SRV   []DnsRecord `json:"srv"`
	TXT   []DnsRecord `json:"txt"`
}

type DnsRecord struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
	Pref  uint16 `json:"pref"`
}

type DnsRootServer struct {
	Host string `json:"host"`
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

var RootServers = []DnsRootServer{
	{
		Host: "a.root-servers.net",
		IPv4: "198.41.0.4",
		IPv6: "2001:503:ba3e::2:30",
	},
	{
		Host: "b.root-servers.net",
		IPv4: "170.247.170.2",
		IPv6: "2001:500:84::b",
	},
	{
		Host: "c.root-servers.net",
		IPv4: "192.33.4.12",
		IPv6: "2001:500:2::c",
	},
	{
		Host: "d.root-servers.net",
		IPv4: "199.7.91.13",
		IPv6: "2001:500:2d::d",
	},
	{
		Host: "e.root-servers.net",
		IPv4: "192.203.230.10",
		IPv6: "2001:500:a8::e",
	},
	{
		Host: "f.root-servers.net",
		IPv4: "192.5.5.241",
		IPv6: "2001:500:2f::f",
	},
	{
		Host: "g.root-servers.net",
		IPv4: "192.112.36.4",
		IPv6: "2001:500:12::d0d",
	},
	{
		Host: "h.root-servers.net",
		IPv4: "198.97.190.53",
		IPv6: "2001:500:1::53",
	},
	{
		Host: "i.root-servers.net",
		IPv4: "192.36.148.17",
		IPv6: "2001:7fe::53",
	},
	{
		Host: "j.root-servers.net",
		IPv4: "192.58.128.30",
		IPv6: "2001:503:c27::2:30",
	},
	{
		Host: "k.root-servers.net",
		IPv4: "193.0.14.129",
		IPv6: "2001:7fd::1",
	},
	{
		Host: "l.root-servers.net",
		IPv4: "199.7.83.42",
		IPv6: "2001:500:9f::42",
	},
	{
		Host: "m.root-servers.net",
		IPv4: "202.12.27.33",
		IPv6: "2001:dc3::35",
	},
}

var Registry DnsRegistry

func Listen(host string, port int, r DnsRegistry) error {
	Registry = r
	fmt.Printf("Starting DNS server on %s:%d\n", host, port)
	cnx, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(host),
		Port: port,
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer cnx.Close()

	fmt.Printf("Listening on %s:%d\n", host, port)

	for {
		handleConnection(cnx)
	}
}

func handleConnection(cnx *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, addr, err := cnx.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	fmt.Printf("[→] Received %d bytes from %s\n", n, addr.String())

	handleQuery(cnx, addr, buffer[:n])
}

func handleQuery(cnx *net.UDPConn, addr *net.UDPAddr, buffer []byte) {
	var msg dnsmessage.Message
	err := msg.Unpack(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	record, domain := ExtractQuery(msg.Questions[0].Name.String())
	recordType := msg.Questions[0].Type
	fmt.Printf("[*] Query (%s) for %s in %s\n", recordType, record, domain)

	// Check if domain is in registry
	dnsDomain := SearchDomain(domain)

	if dnsDomain.Name == "" {
		fmt.Printf("[!] Domain %s not found\n", domain)
		answer := createNXDOMAINMessage(msg.Header.ID, msg.Questions[0], domain)
		SendResponse(cnx, addr, answer, false)
		return
	}

	// Check if record is in domain
	_, answers := checkRecords(dnsDomain, record, msg.Questions[0], recordType)

	if len(answers) == 0 {
		fmt.Printf("[!] Record %s not found in domain %s\n", record, domain)
		answer := createNXDOMAINMessage(msg.Header.ID, msg.Questions[0], domain)
		SendResponse(cnx, addr, answer, false)
		return
	}

	response := createMessage(msg.Header.ID, msg.Questions[0], answers, domain, true, true)

	SendResponse(cnx, addr, response, true)
}

func SendResponse(cnx *net.UDPConn, addr *net.UDPAddr, msg dnsmessage.Message, success bool) {
	buffer, err := msg.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}
	if success {
		fmt.Printf("[←] Sending response for %s\n", msg.Questions[0].Name.String())
	} else {
		fmt.Println("[←] Sending NXDOMAIN response")
	}
	_, err = cnx.WriteToUDP(buffer, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createMessage(id uint16, q dnsmessage.Question, answers []dnsmessage.Resource, domain string, response bool, authoritative bool) dnsmessage.Message {
	authority := []dnsmessage.Resource{
		{
			Header: dnsmessage.ResourceHeader{
				Name:  q.Name,
				Type:  dnsmessage.TypeNS,
				Class: dnsmessage.ClassINET,
				TTL:   300,
			},
			Body: &dnsmessage.NSResource{
				NS: dnsmessage.MustNewName("ns1." + domain),
			},
		},
	}

	return dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:                 id,
			Response:           response,
			Authoritative:      authoritative,
			Truncated:          false,
			RecursionDesired:   true,
			RecursionAvailable: true,
			RCode:              dnsmessage.RCodeSuccess,
		},
		Questions: []dnsmessage.Question{
			q,
		},
		Answers:     answers,
		Authorities: authority,
	}
}

func checkARecords(d DnsDomain, record string) DnsRecord {
	for _, r := range d.A {
		if r.Name == record {
			return r
		}
	}
	return DnsRecord{}
}

func checkAAAARecords(d DnsDomain, record string) DnsRecord {
	for _, r := range d.AAAA {
		if r.Name == record {
			return r
		}
	}
	return DnsRecord{}
}

func checkMXRecords(d DnsDomain, record string) DnsRecord {
	for _, r := range d.MX {
		if r.Name == record {
			return r
		}
	}
	return DnsRecord{}
}

func checkRecords(d DnsDomain, record string, q dnsmessage.Question, recordType ...dnsmessage.Type) ([]DnsRecord, []dnsmessage.Resource) {
	aRecords := []DnsRecord{}
	aaaaRecords := []DnsRecord{}
	mxRecords := []DnsRecord{}
	answers := []dnsmessage.Resource{}
	for _, t := range recordType {
		if t == dnsmessage.TypeA || t == dnsmessage.TypeALL {
			aRecords = append(aRecords, checkARecords(d, record))
			for _, r := range aRecords {
				if r.Name != "" {
					answers = append(answers, createARecord(r.Name, r.Value, r.TTL, q))
				}
			}
		}
		if t == dnsmessage.TypeAAAA || t == dnsmessage.TypeALL {
			aaaaRecords = append(aaaaRecords, checkAAAARecords(d, record))
			for _, r := range aaaaRecords {
				if r.Name != "" {
					answers = append(answers, createAAAARecord(r.Name, r.Value, r.TTL, q))
				}
			}
		}
		if t == dnsmessage.TypeMX || t == dnsmessage.TypeALL {
			mxRecords = append(mxRecords, checkMXRecords(d, record))
			for _, r := range mxRecords {
				if r.Name != "" {
					answers = append(answers, createMXRecord(r.Name, r.Value, r.TTL, r.Pref, q))
				}
			}
		}
	}
	records := append(aRecords, mxRecords...)
	return records, answers
}

func createARecord(name string, value string, ttl int, q dnsmessage.Question) dnsmessage.Resource {
	ip := net.ParseIP(value)
	ip = ip.To4()
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  q.Name,
			Type:  dnsmessage.TypeA,
			Class: dnsmessage.ClassINET,
			TTL:   uint32(ttl),
		},
		Body: &dnsmessage.AResource{
			A: [4]byte{ip[0], ip[1], ip[2], ip[3]},
		},
	}
}

func createAAAARecord(name string, value string, ttl int, q dnsmessage.Question) dnsmessage.Resource {
	ip := net.ParseIP(value)
	ip = ip.To16()
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  q.Name,
			Type:  dnsmessage.TypeAAAA,
			Class: dnsmessage.ClassINET,
			TTL:   uint32(ttl),
		},
		Body: &dnsmessage.AAAAResource{
			AAAA: [16]byte{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15]},
		},
	}
}

func createMXRecord(name string, value string, ttl int, pref uint16, q dnsmessage.Question) dnsmessage.Resource {
	if !strings.HasSuffix(name, ".") {
		name += "."
	}
	if !strings.HasSuffix(value, ".") {
		value += "."
	}
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  q.Name,
			Type:  dnsmessage.TypeMX,
			Class: dnsmessage.ClassINET,
			TTL:   uint32(ttl),
		},
		Body: &dnsmessage.MXResource{
			Pref: pref,
			MX:   dnsmessage.MustNewName(value),
		},
	}
}

func createNXDOMAINMessage(id uint16, q dnsmessage.Question, domain string) dnsmessage.Message {
	return dnsmessage.Message{
		Header: dnsmessage.Header{
			ID:                 id,
			Response:           true,
			Authoritative:      true,
			Truncated:          false,
			RecursionDesired:   true,
			RecursionAvailable: true,
			RCode:              dnsmessage.RCodeNameError, // NXDOMAIN
		},
		Questions: []dnsmessage.Question{
			q,
		},
		Answers:     nil,
		Authorities: nil,
	}
}
