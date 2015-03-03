package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

var crit, warn int64
var typ, checkname, target, tocheck string
var debug, jsonflag bool

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
======================
Example:
  %% %s -type SRV -tocheck _xmpp-client._tcp.google.com -target xmpp.l.google.com. -crit 200000000
  WARNING response took too long - %s: SRV xmpp.l.google.com. |time=33657371;1000000;200000000

`, bin, bin, bin)
		flag.PrintDefaults()
	}

	flag.Int64Var(&warn, "warn", 1000000, "Warnning if request time in nano seconds is greater than given value")
	flag.Int64Var(&crit, "crit", 2000000, "Critical if request time in nano seconds is greater than given value")
	flag.StringVar(&typ, "type", "A", "Type of the DNS Record, A, SRV, CNAME, MX, NS, TXT, PTR")
	flag.StringVar(&checkname, "checkname", "Check DNS", "Name to show in nagios message.")
	flag.StringVar(&target, "target", "173.194.72.125", "String target as result to check as dig would return.")
	flag.StringVar(&tocheck, "tocheck", "alt3.xmpp.l.google.com.", "String to check as you would use with dig.")
	flag.BoolVar(&jsonflag, "json", false, "Enable json output, instead of nagios style.")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
}

// A-Record
func reqARecord(name string) ([]string, error) {
	if names, err := net.LookupHost(name); err != nil {
		log.Fatal("ERR:", err)
		return nil, err
	} else {
		return names, nil
	}

}

// CNAME-Record
func reqCNAME(name string) (string, error) {
	if names, err := net.LookupCNAME(name); err != nil {
		return "", err
	} else {
		return names, nil
	}
}

// SRV-Record  _xmpp-client._tcp.google.com
func reqSRV(name string) (string, []*net.SRV, error) {
	name = strings.Replace(name, "_", "", -1)
	l := strings.Split(name, ".")
	service, prot, domain := l[0], l[1], strings.Join(l[2:], ".")
	log.Println(service, prot, domain)
	if cname, addrs, err := net.LookupSRV(service, prot, domain); err != nil {
		return "", nil, err
	} else {
		return cname, addrs, nil
	}
}

// PTR-Record - reverse lookup
func reqPTRRecord(ip string) ([]string, error) {
	if names, err := net.LookupAddr(ip); err != nil {
		log.Fatal("ERR:", err)
		return nil, err
	} else {
		return names, nil
	}
}

// TXT-Record
func reqTXTRecord(name string) ([]string, error) {
	if names, err := net.LookupTXT(name); err != nil {
		log.Fatal("ERR:", err)
		return nil, err
	} else {
		return names, nil
	}
}

// MX-Record
func reqMXRecord(name string) ([]*net.MX, error) {
	if mx, err := net.LookupMX(name); err != nil {
		log.Fatal("ERR:", err)
		return nil, err
	} else {
		return mx, nil
	}
}

// NS-Record
func reqNSRecord(name string) (ns []*net.NS, err error) {
	if ns, err := net.LookupNS(name); err != nil {
		log.Fatal("ERR:", err)
		return nil, err
	} else {
		return ns, nil
	}
}

func main() {
	var exitcode int
	exitcode = 2
	state := "CRITICAL - record not found"

	bin := path.Base(os.Args[0])

	log.SetPrefix(fmt.Sprintf("* %s - ", bin))
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	flag.Parse()
	msg := fmt.Sprintf("- %s: %s %s ", bin, typ, target)

	// start the timer
	start := time.Now()

	switch typ {
	case "A":
		records, _ := reqARecord(tocheck)
		for _, record := range records {
			log.Printf("Found A Record: %s", record)
			if record == target {
				state = "OK"
				exitcode = 0
			}
		}

	case "PTR":
		records, _ := reqPTRRecord(tocheck)
		for _, record := range records {
			log.Printf("Found PTR Record: %s", record)
			if record == target {
				state = "OK"
				exitcode = 0
			}
		}

	case "TXT":
		records, _ := reqTXTRecord(tocheck)
		for _, record := range records {
			log.Printf("Found TXT Record: %s", record)
			if record == target {
				state = "OK"
				exitcode = 0
			}
		}

	case "MX":
		records, _ := reqMXRecord(tocheck)
		for _, record := range records {
			log.Printf("Found MX Record: %s", record)
			if record.Host == target {
				state = "OK"
				exitcode = 0
			}
		}

	case "NS":
		records, _ := reqNSRecord(tocheck)
		for _, record := range records {
			log.Printf("Found NS Record: %s", record)
			if record.Host == target {
				state = "OK"
				exitcode = 0
			}
		}

	case "CNAME":
		record, _ := reqCNAME(tocheck)
		if record == target {
			state = "OK"
			exitcode = 0
		}

	case "SRV":
		if cname, addrs, err := reqSRV(tocheck); err != nil {
			fmt.Println("UNKNOWN:", err)
			os.Exit(3)
		} else {
			log.Printf("SRV:\n\tcname: %s\n", cname)
			for _, addr := range addrs {
				log.Printf("\taddr.Target:\t%v\n", addr.Target)
				if addr.Target == target {
					state = "OK"
					exitcode = 0
				}
			}
		}

	default:
		log.Fatal("Unknown record type %s\n", typ)
	}

	now := time.Now()
	ts := now.Sub(start)

	if jsonflag {
		fmt.Printf("{\"time\": %d}\n", ts.Nanoseconds())
		os.Exit(0)
	}

	if ts.Nanoseconds() > crit {
		state = "CRITICAL response took too long"
		exitcode = 2
	} else if ts.Nanoseconds() > warn {
		state = "WARNING response took too long"
		exitcode = 1
	}

	fmt.Printf("%s %s|time=%d;%d;%d\n", state, msg, ts.Nanoseconds(), warn, crit)
	os.Exit(exitcode)
}
