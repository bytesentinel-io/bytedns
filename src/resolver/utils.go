package resolver

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

func NewRandom(len int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len)))
	return int(n.Int64())
}

func RandomRootServer(r *DnsRegistry) DnsRootServer {
	len := len(r.RootServers)
	// Random number generator (0 - len)
	random := NewRandom(len)
	return r.RootServers[random]
}

func ReverseIP(ip string) string {
	parts := strings.Split(ip, ".")
	reversedParts := make([]string, len(parts))
	for i := 0; i < len(parts); i++ {
		reversedParts[i] = parts[len(parts)-1-i]
	}
	return strings.Join(reversedParts, ".")
}

func ExtractDomain(q string) string {
	parts := strings.Split(q, ".")
	return strings.Join(parts[len(parts)-3:], ".")
}

func ExtractQuery(q string) (string, string) {
	if CheckIfRoot(q) {
		return "@", q
	}
	parts := strings.Split(q, ".")
	return parts[0], strings.Join(parts[1:], ".")
}

func SearchDomain(name string) DnsDomain {
	for _, domain := range Registry.Domains {
		if domain.Name == name {
			return domain
		}
	}
	return DnsDomain{}
}

func IPv4ToBytes(ip net.IPAddr) []byte {
	fmt.Println(ip)
	parts := strings.Split(ip.String(), ".")
	bytes := make([]byte, len(parts))
	for i, part := range parts {
		b, _ := strconv.ParseInt(part, 10, 8)
		bytes[i] = byte(b)
	}
	return bytes
}

func CheckIfRoot(q string) bool {
	parts := strings.Split(q, ".")
	// Remove the last empty string
	parts = parts[:len(parts)-1]
	return len(parts) == 2
}
