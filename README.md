# nagiosplugins
[![Build Status](https://secure.travis-ci.org/szuecs/nagiosplugins.png?branch=master)](http://travis-ci.org/szuecs/nagiosplugins)

## findbyname
This nagiosplugin lets you count files by name and or mtime in a given
target directory. Optionally you can get json output, if you don't
care about nagios style output and exitcode.

### Usage of findbyname

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
