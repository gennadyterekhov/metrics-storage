package ipaddr

import (
	"fmt"
	"net"
)

func GetHostIPAsString() string {
	ip, err := getFirstNonPrivateIP()
	if err != nil {
		return ""
	}
	return ip.String()
}

func getFirstNonPrivateIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if !isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find public ip")
}

func isPrivateIP(ip net.IP) bool {
	for _, cidr := range []string{
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
