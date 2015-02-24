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
