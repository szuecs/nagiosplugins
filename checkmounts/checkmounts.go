package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var crit, warn int
var typ, regex, checkname, fstab, proc_mounts string
var debug, jsonflag bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
======================
Example:
  %% %s -type nfs -regex '^l.*[0-9]$'
  CRITICAL - %s|count=7;2;2

`, os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&warn, "warn", 2, "Warnning if count is greater than given value")
	flag.IntVar(&crit, "crit", 2, "Critical if count is greater than given value")
	flag.StringVar(&typ, "type", "", "Type of the mount, p.e. nfs")
	flag.StringVar(&regex, "regex", "", "Regular expression to filter mounts")
	flag.StringVar(&checkname, "checkname", "Check NFS Mounts", "Name to show in nagios message.")
	flag.StringVar(&fstab, "fstab", "/etc/fstab", "Path to fstab.")
	flag.StringVar(&proc_mounts, "proc_mounts", "/proc/mounts", "Path to /proc/mounts.")
	flag.BoolVar(&jsonflag, "json", false, "Enable json output, instead of nagios style.")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
}

type NagiosData struct {
	Msg      string
	Exitcode int
	Perfdata string
}

func check(count int) NagiosData {
	exitcode := 0
	msg := fmt.Sprintf("OK - %s", checkname)
	perfdata := fmt.Sprintf("count=%d;%d;%d", count, warn, crit)

	if count >= crit {
		msg = fmt.Sprintf("CRITICAL - %s", checkname)
	} else if count >= warn {
		msg = fmt.Sprintf("WARNING - %s", checkname)
	}

	return NagiosData{msg, exitcode, perfdata}
}

// filter by type implicitly
func get_devices(path string) []string {
	result := []string{}
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		if words[2] == typ {
			log.Printf("%s - words: %s %s\n", path, words[0], words[2])
			result = append(result, words[0])
		}
	}

	return result
}

func filter_by_name(strings []string, rs string) []string {
	result := []string{}
	r, err := regexp.Compile(rs)
	if err != nil {
		log.Panic("Can not compile given string to regexp: %s", rs)
	}
	for _, s := range strings {
		if s == "" {
			continue
		}
		if r.MatchString(s) {
			log.Printf("%s matches %s\n", rs, s)
			result = append(result, s)
		}
	}
	return result
}

type MergedEntry struct {
	Device  string
	Mounted bool
}

func run() NagiosData {
	fstab_entries := get_devices(fstab)
	procmount_entries := get_devices(proc_mounts)
	entries := []MergedEntry{}

	if regex != "" {
		fstab_entries = filter_by_name(fstab_entries, regex)
		procmount_entries = filter_by_name(procmount_entries, regex)
	}

	for _, fe := range fstab_entries {
		entry := MergedEntry{fe, false}
		for _, pe := range procmount_entries {
			if fe == pe {
				entry.Mounted = true
			}
		}
		entries = append(entries, entry)
	}
	log.Print(entries)

	// json output: do not care about exitcode
	if jsonflag {
		b, err := json.Marshal(entries)
		if err != nil {
			fmt.Println("Marshall err: ", err)
		}
		fmt.Printf("%s\n", b)
		os.Exit(0)
	}

	not_mounted := 0
	for _, e := range entries {
		if !e.Mounted {
			not_mounted += 1
		}
	}
	return check(not_mounted)
}

func main() {
	flag.Parse()
	log.SetPrefix(fmt.Sprintf("* %s - ", os.Args[0]))
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	// calls run() and make sure any log.Panic() calls will exit
	// with errorcode 3 as nagios defined it.
	func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("UNKNOWN %s - run function failed with panic: %s\n", os.Args[0], err)
				os.Exit(3)
			}
		}()

		data := run()
		fmt.Printf("%s|%s\n", data.Msg, data.Perfdata)
		os.Exit(data.Exitcode)
	}()

}
