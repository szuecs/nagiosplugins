package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

var crit, warn int
var path, regex, checkname string
var mtime int64
var debug, jsonflag bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s
======================
Example:
  %% %s -path /tmp -mtime 3 -regex '^l.*[0-9]$'
  CRITICAL - %s|count=7;2;2

`, os.Args[0], os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&warn, "warn", 2, "Warnning if count is greater than given value")
	flag.IntVar(&crit, "crit", 2, "Critical if count is greater than given value")
	flag.Int64Var(&mtime, "mtime", -1, "Files with mtime in hours bigger than given")
	flag.StringVar(&path, "path", "", "path in which to find files")
	flag.StringVar(&regex, "regex", "", "Regular expression to filter files")
	flag.StringVar(&checkname, "checkname", "findbyname", "Name to show in nagios message.")
	flag.BoolVar(&jsonflag, "json", false, "Enable json output, instead of nagios style.")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
}

type NagiosData struct {
	Msg      string
	Exitcode int
	Perfdata string
}

func get_files(path string) []os.FileInfo {
	file_infos, err := ioutil.ReadDir(path)
	if err != nil {
		log.Panic("Can not get directory listing", err)
	}
	return file_infos
}

func filter_by_mtime(file_infos []os.FileInfo, mtime int64) {
	now := time.Now()
	tduration := time.Duration(mtime) * time.Hour
	for i := 0; i < len(file_infos); i++ {
		if file_infos[i] == nil {
			continue
		}
		fmtime := now.Sub(file_infos[i].ModTime())
		if fmtime > tduration {
			log.Printf("mtime filter %s\n", file_infos[i].Name())
		} else {
			file_infos[i] = nil
		}
	}
}

func filter_by_name(file_infos []os.FileInfo, rs string) {
	r, err := regexp.Compile(rs)
	if err != nil {
		log.Panic("Can not compile given string to regexp: %s", rs)
	}
	for i := 0; i < len(file_infos); i++ {
		if file_infos[i] == nil {
			continue
		}
		if r.MatchString(file_infos[i].Name()) {
			log.Printf("%s matches %s\n", rs, file_infos[i].Name())
		} else {
			file_infos[i] = nil
		}

	}
}

func filter_nil(file_infos []os.FileInfo) []os.FileInfo {
	result := []os.FileInfo{}
	for _, fi := range file_infos {
		if fi != nil {
			result = append(result, fi)
		}
	}
	return result
}

// Returns nagios Exitcodes for given list, checking thresholds
func check(file_infos []os.FileInfo) NagiosData {
	exitcode := 0
	msg := fmt.Sprintf("OK - %s", checkname)
	count := len(file_infos)
	perfdata := fmt.Sprintf("count=%d;%d;%d", count, warn, crit)

	if count >= crit {
		msg = fmt.Sprintf("CRITICAL - %s", checkname)
	} else if count >= warn {
		msg = fmt.Sprintf("WARNING - %s", checkname)
	}

	return NagiosData{msg, exitcode, perfdata}
}

func run() NagiosData {
	if path == "" {
		log.Panic("No path specified")
	}

	files := get_files(path)

	if mtime >= 0 {
		filter_by_mtime(files, mtime)
	}
	if regex != "" {
		filter_by_name(files, regex)
	}
	files = filter_nil(files)

	if debug {
		for _, fi := range files {
			log.Println("Alert:", fi.Name())
		}
	}

	// do not care about exitcode
	if jsonflag {
		fmt.Printf("[{\"count\": %d}]\n", len(files))
		os.Exit(0)
	}

	return check(files)
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
