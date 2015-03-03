package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

var crit, warn uint64
var target, checkname, memkey, metricskey, metricstype string
var debug, jsonflag bool

func init() {
	bin := path.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Check metrics url http://localhost:8081/metrics

Usage of %s
======================
Example:
  %% %s -url http://localhost:8081/metrics
  OK - Check Gomelon|Alloc=1441944 TotalAlloc=2144632 Sys=5114104 Lookups=40 Mallocs=4056 Frees=3197 HeapAlloc=1441944 HeapSys=2899968 HeapIdle=1138688 HeapInuse=1761280 HeapReleased=1064960 HeapObjects=859 StackInuse=245760 StackSys=245760 MSpanInuse=9152 MSpanSys=16384 MCacheInuse=1200 MCacheSys=16384 BuckHashSys=1440592 GCSys=202793 OtherSys=292223 NextGC=2853216 LastGC=1425328529426483999 PauseTotalNs=132651838 NumGC=31 EnableGC=true DebugGC=false

`, bin, bin)
		flag.PrintDefaults()
	}

	flag.Uint64Var(&warn, "warn", 2, "Warnning if count is greater than given value")
	flag.Uint64Var(&crit, "crit", 2, "Critical if count is greater than given value")
	flag.StringVar(&memkey, "memkey", "", "Which item to check from runtime.Memstats.")
	flag.StringVar(&metricskey, "metricskey", "", "Which item to check from runtime.Metrics.")
	flag.StringVar(&metricstype, "metricstype", "", "Which item to check from metrics Counters or Gauges.")
	flag.StringVar(&checkname, "checkname", "Check Gomelon", "Name to show in nagios message.")
	flag.StringVar(&target, "url", "http://localhost:8081/metrics", "Metrics URL.")
	flag.BoolVar(&jsonflag, "json", false, "Enable json output, instead of nagios style.")
	flag.BoolVar(&debug, "debug", false, "Enable debug output")
}

type NagiosData struct {
	Msg      string
	Exitcode int
	Perfdata string
}

func (nagios_data *NagiosData) check(count uint64) {
	if count >= crit {
		nagios_data.Msg = fmt.Sprintf("CRITICAL - %s", checkname)
		nagios_data.Exitcode = 2
	} else if count >= warn {
		nagios_data.Msg = fmt.Sprintf("WARNING - %s", checkname)
		nagios_data.Exitcode = 1
	} else {
		nagios_data.Msg = fmt.Sprintf("OK - %s", checkname)
		nagios_data.Exitcode = 0
	}
}

type Stats struct {
	cmdline  map[string]string
	Memstats runtime.MemStats
	Metrics  map[string]map[string]int
}

func get_stats() Stats {
	resp, err := http.Get(target)
	if err != nil {
		log.Fatal("Could not GET %s, cause %s", target, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Could not read body, caused by: %s", err)
	}

	if jsonflag {
		fmt.Printf("%s\n", body)
		os.Exit(0)
	}

	var v Stats
	if err := json.Unmarshal(body, &v); err != nil {
		log.Fatal("Could not unmarshal JSON data, caused by: %s", err)
	}

	return v
}
func parseStats(stat Stats) NagiosData {
	mem := stat.Memstats
	msg := checkname
	exitcode := 0
	perfdata := []string{
		fmt.Sprintf("Alloc=%d", mem.Alloc),
		fmt.Sprintf("TotalAlloc=%d", mem.TotalAlloc),
		fmt.Sprintf("Sys=%d", mem.Sys),
		fmt.Sprintf("Lookups=%d", mem.Lookups),
		fmt.Sprintf("Mallocs=%d", mem.Mallocs),
		fmt.Sprintf("Frees=%d", mem.Frees),
		fmt.Sprintf("HeapAlloc=%d", mem.HeapAlloc),
		fmt.Sprintf("HeapSys=%d", mem.HeapSys),
		fmt.Sprintf("HeapIdle=%d", mem.HeapIdle),
		fmt.Sprintf("HeapInuse=%d", mem.HeapInuse),
		fmt.Sprintf("HeapReleased=%d", mem.HeapReleased),
		fmt.Sprintf("HeapObjects=%d", mem.HeapObjects),
		fmt.Sprintf("StackInuse=%d", mem.StackInuse),
		fmt.Sprintf("StackSys=%d", mem.StackSys),
		fmt.Sprintf("MSpanInuse=%d", mem.MSpanInuse),
		fmt.Sprintf("MSpanSys=%d", mem.MSpanSys),
		fmt.Sprintf("MCacheInuse=%d", mem.MCacheInuse),
		fmt.Sprintf("MCacheSys=%d", mem.MCacheSys),
		fmt.Sprintf("BuckHashSys=%d", mem.BuckHashSys),
		fmt.Sprintf("GCSys=%d", mem.GCSys),
		fmt.Sprintf("OtherSys=%d", mem.OtherSys),
		fmt.Sprintf("NextGC=%d", mem.NextGC),
		fmt.Sprintf("LastGC=%d", mem.LastGC),
		fmt.Sprintf("PauseTotalNs=%d", mem.PauseTotalNs),
		//fmt.Sprintf("PauseNs=%d", mem.PauseNs),
		//fmt.Sprintf("PauseEnd=%d", mem.PauseEnd),
		fmt.Sprintf("NumGC=%d", mem.NumGC),
		fmt.Sprintf("EnableGC=%v", mem.EnableGC),
		fmt.Sprintf("DebugGC=%v", mem.DebugGC),
	}
	return NagiosData{msg, exitcode, strings.Join(perfdata, " ")}
}

func run() NagiosData {
	v := get_stats()

	// json output: do not care about exitcode
	if jsonflag {
		b, err := json.Marshal(v.Memstats)
		if err != nil {
			fmt.Println("Marshall err: ", err)
		}
		fmt.Printf("%s\n", b)
		os.Exit(0)
	}

	nagios_data := parseStats(v)
	if memkey != "" {
		// use reflection to get the field by name
		r := reflect.ValueOf(v.Memstats)
		f := reflect.Indirect(r).FieldByName(memkey)
		nagios_data.check(uint64(f.Uint()))
	}

	if metricstype != "" && metricskey != "" {
		val := v.Metrics[metricstype][metricskey]
		nagios_data.check(uint64(val))
	}
	return nagios_data
}

func main() {
	flag.Parse()
	log.SetPrefix(fmt.Sprintf("* %s - ", path.Base(os.Args[0])))
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
