# nagiosplugins
[![Build Status](https://secure.travis-ci.org/szuecs/nagiosplugins.png?branch=master)](http://travis-ci.org/szuecs/nagiosplugins)

All nagiosplugins are able to output nagiosplugin style output with
perfdata or with --json parameter you can get json written to STDOUT.

## Install

Install preconditions, tested on Ubuntu trusty (ubuntu14.04):

    % apt-get install golang-go git-core

Setup GOPATH:

    % export GOPATH=$HOME/go
    % export GOBIN=$HOME/go/bin
    % mkdir -p $GOPATH

Install all plugins into your $GOPATH/bin with:

    % go get -u github.com/szuecs/nagiosplugins/...

## check_gomelon
This checks the exposed /metrics url of
[gomelon](https://github.com/goburrow/gomelon), a go web framework.

### Usage of check_gomelon

    % ./check_gomelon -h
    Check metrics url http://localhost:8081/metrics

    Usage of check_gomelon
    ======================
    Example:
      % check_gomelon -url http://localhost:8081/metrics
      OK - Check Gomelon|Alloc=1441944 TotalAlloc=2144632 Sys=5114104 Lookups=40 Mallocs=4056 Frees=3197 HeapAlloc=1441944 HeapSys=2899968 HeapIdle=1138688 HeapInuse=1761280 HeapReleased=1064960 HeapObjects=859 StackInuse=245760 StackSys=245760 MSpanInuse=9152 MSpanSys=16384 MCacheInuse=1200 MCacheSys=16384 BuckHashSys=1440592 GCSys=202793 OtherSys=292223 NextGC=2853216 LastGC=1425328529426483999 PauseTotalNs=132651838 NumGC=31 EnableGC=true DebugGC=false

      -checkname="Check Gomelon": Name to show in nagios message.
      -crit=2: Critical if count is greater than given value
      -debug=false: Enable debug output
      -json=false: Enable json output, instead of nagios style.
      -memkey="": Which item to check from runtime.Memstats.
      -metricskey="": Which item to check from runtime.Metrics.
      -metricstype="": Which item to check from metrics Counters or Gauges.
      -url="http://localhost:8081/metrics": Metrics URL.
      -warn=2: Warnning if count is greater than given value

## checkdns
This nagiosplugin lets you check if the target is within the result
set of the given tocheck request for a given DNS Record. You can check
the request time in nanoseconds.

### Usage of checkdns

    % checkdns -h

    Usage of checkdns
    ======================
    Example:
    % checkdns -type SRV -tocheck _xmpp-client._tcp.google.com -target xmpp.l.google.com. -crit 200000000
    WARNING response took too long - checkdns: SRV xmpp.l.google.com. |time=33657371;1000000;200000000

    -checkname="Check DNS": Name to show in nagios message.
    -crit=2000000: Critical if request time in nano seconds is greater than given value
    -debug=false: Enable debug output
    -json=false: Enable json output, instead of nagios style.
    -target="173.194.72.125": String target as result to check as dig would return.
    -tocheck="alt3.xmpp.l.google.com.": String to check as you would use with dig.
    -type="A": Type of the DNS Record, A, SRV, CNAME, MX, NS, TXT, PTR
    -warn=1000000: Warnning if request time in nano seconds is greater than given value

## checkmounts
This nagiosplugin lets you check the state of /proc/mounts vs. the
definition of /etc/fstab.

### Usage of checkmounts

    % checkmounts -h

    Usage of checkmounts
    ======================
    Example:
      % checkmounts -type nfs -regex '^l.*[0-9]$'
      CRITICAL - checkmounts|count=7;2;2

      -checkname="Check NFS Mounts": Name to show in nagios message.
      -crit=2: Critical if count is greater than given value
      -debug=false: Enable debug output
      -fstab="/etc/fstab": Path to fstab.
      -json=false: Enable json output, instead of nagios style.
      -proc_mounts="/proc/mounts": Path to /proc/mounts.
      -regex="": Regular expression to filter mounts
      -type="": Type of the mount, p.e. nfs
      -warn=2: Warnning if count is greater than given value

## findbyname
This nagiosplugin lets you count files by name and or mtime in a given
target directory. Optionally you can get json output, if you don't
care about nagios style output and exitcode.

### Usage of findbyname

    % findbyname -h

    Usage of findbyname
    ======================
    Example:
      % findbyname -path /tmp -mtime 3 -regex '^l.*[0-9]$'
      CRITICAL - findbyname|count=7;2;2

      -checkname="findbyname": Name to show in nagios message.
      -crit=2: Critical if count is greater than given value
      -debug=false: Enable debug output
      -json=false: Enable json output, instead of nagios style.
      -mtime=-1: Files with mtime in hours bigger than given
      -path="": path in which to find files
      -regex="": Regular expression to filter files
      -warn=2: Warnning if count is greater than given value
