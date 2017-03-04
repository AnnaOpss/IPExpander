package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	os.Exit(nixMain(os.Args[1:]))
}

func nixMain(ipranges []string) (code int) {
	if len(ipranges) == 0 {
		fmt.Println("Please provide at least an ip range in Classless inter-domain routing (CIDR) form")
		//Signal that not enough args were provided
		return 1
	}

	for _, iprange := range ipranges {
		fmt.Fprintln(os.Stderr, "Now printing "+iprange)
		if strings.Contains(iprange, "-") {
			ips, err := parseDashed(iprange)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				//Signal that an error happened
				code = 2
				break
			}
			for _, ip := range ips {
				fmt.Println(ip)
			}
		} else {
			ip, ipnet, err := net.ParseCIDR(iprange)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				//Signal that an error happened
				code = 2
				break
			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				fmt.Println(ip)
			}
		}
	}
	return
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

//This is the single fucking most unexpectedly complex piece of garbage I have ever written
func parseDashed(in string) ([]net.IP, error) {
	sections := strings.Split(in, ".")
	//An array of slices!
	var ranges [4][]byte
	//This label allows the flow to jump out of the switch below
sections:
	for i, s := range sections {
		var err error
		limits := strings.Split(s, "-")
		//Make sure there was one and only one dash (weird switch syntax here)
		switch l := len(limits); {
		case l < 2:
			//No dashes, skip this section
			value, err := sanitize(limits[0])
			if err != nil {
				return nil, err
			}
			//This is totally not a cast
			ranges[i] = []byte{value}
			continue sections
		case l > 2:
			//More than one dash, don't understand it
			return nil, errors.New("Ranges should contain a single dash per ip byte.")
		}
		var start byte
		if limits[0] == "" {
			//User did not specify a start, let's use 0
			start = byte(0)
		} else {
			start, err = sanitize(limits[0])
			if err != nil {
				return nil, err
			}
		}
		var end byte
		if limits[1] == "" {
			//User did not specify an end, let's use 255
			end = byte(255)
		} else {
			end, err = sanitize(limits[1])
			if err != nil {
				return nil, err
			}
		}
		//User specified a reversed range, let's straighten it up a little
		if end < start {
			start, end = end, start
		}
		//Avoid reallocations by append, I already know how long the slice is going to be
		ranges[i] = make([]byte, 0, int(end-start)+1)
		//Generate the bytes for the IPs for this range
		for b := int(start); b <= int(end); b++ {
			ranges[i] = append(ranges[i], byte(b))
		}
	}
	//Let's allocate all the memory in one call
	ips := make([]net.IP, 0, len(ranges[0])*len(ranges[1])*len(ranges[2])*len(ranges[3]))
	//Now let's do the job we have to do
	for _, i := range ranges[0] {
		for _, j := range ranges[1] {
			for _, k := range ranges[2] {
				for _, l := range ranges[3] {
					//It would be better to just output the IPs here, saving all of them does look like a waste of memory
					ips = append(ips, net.IPv4(i, j, k, l))
				}
			}
		}
	}
	return ips, nil
}

func sanitize(in string) (out byte, err error) {
	tmp, err := strconv.Atoi(in)
	if err != nil {
		return
	}
	if tmp > 255 || tmp < 0 {
		return byte(0), errors.New("Index should be greater than 0 and less than 255")
	}
	return byte(tmp), nil
}
