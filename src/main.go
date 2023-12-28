package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bytesentinel-io/bytedns/resolver"
)

var registry = resolver.DnsRegistry{
	Domains: []resolver.DnsDomain{
		{
			Name: "bytie.lab.",
			A: []resolver.DnsRecord{
				{
					Name:  "@",
					Value: "10.100.0.1",
					TTL:   300,
				},
				{
					Name:  "pihole",
					Value: "10.100.0.2",
					TTL:   300,
				},
				{
					Name:  "guaca",
					Value: "10.100.0.5",
					TTL:   300,
				},
			},
			AAAA: []resolver.DnsRecord{
				{
					Name:  "@",
					Value: "fd00::1",
					TTL:   300,
				},
			},
			CNAME: []resolver.DnsRecord{
				{
					Name:  "www",
					Value: "@",
					TTL:   300,
				},
			},
			MX: []resolver.DnsRecord{
				{
					Name:  "@",
					Value: "10.100.0.3",
					TTL:   300,
					Pref:  0,
				},
			},
		},
		{
			Name: "skynet.lab.",
			A: []resolver.DnsRecord{
				{
					Name:  "@",
					Value: "192.168.110.250",
					TTL:   300,
				},
			},
		},
	},
	Forwarding: resolver.DnsForwarding{
		Enabled: true,
		Server:  "1.1.1.1",
	},
	RootServers: resolver.RootServers,
}

func main() {
	host := ""
	port := ""

	if len(os.Args) == 1 {
		host = "0.0.0.0"
		port = "53"
	} else if len(os.Args) == 2 {
		host = os.Args[1]
		port = "53"
	} else if len(os.Args) == 3 {
		host = os.Args[1]
		port = os.Args[2]
	} else {
		fmt.Println("Invalid arguments")
		return
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Invalid port:", err)
		return
	}

	resolver.Listen(host, portInt, registry)
}
