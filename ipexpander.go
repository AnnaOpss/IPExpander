package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	ipranges := os.Args[1:]

	if len(ipranges) == 0 {
		printHelp()
		return
	}

	for _, iprange := range ipranges {
		fmt.Fprintln(os.Stderr, "Now printing "+iprange)
		ip, ipnet, err := net.ParseCIDR(iprange)
		if err != nil {
			log.Fatal(err)
		}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			fmt.Println(ip)
		}
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

func printHelp() {
	fmt.Println("Please provide at least an ip range in Classless inter-domain routing (CIDR) form")
}
