package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var inip = flag.String("ip", "192.168.1.0/24", "The IP address and submask, CIDR notation")

func main() {
	flag.Parse()
	//TODO implementare controllo IP formato correttamente.

	ip, ipnet, err := net.ParseCIDR(*inip)
	if err != nil {
		log.Fatal(err)
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		fmt.Println(ip)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
