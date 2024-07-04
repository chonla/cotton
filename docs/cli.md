---
layout: default
title: CLI
nav_order: 2
permalink: /cli
---

# Command Line Interface (CLI)

## Usage

{% highlight bash %}
  cotton [-d] [-c] [-p] [-r <rootdir>] <testpath|testdir>
  cotton -v
  cotton --help

  -c    compact mode
  -d    debug mode
  -h    display this help
  -i    disable certificate verification
  -p    paranoid mode
  -r string
        set rootDir path
  -s    stop when test failed
  -v    display cotton version
{% endhighlight %}

## Options

By default, Cotton prints out test title, setup title, teardown title, result, and summary.

### `-c` Compact mode

Compact mode is minimal logging execution. Cotton prints out only test title, result, and summary.

### `-d` Debug mode

Debug mode is verbose logging execution. Cotton prints out everything for debugging purpose.

### `-h` Display help

Cotton prints out help information.

### `-i` Insecure mode

Insecure mode will disable certificate verification.

### `-p` Paranoid mode

Paranoid mode is detailed verbose logging execution. Cotton prints out very detailed information for further debugging purpose.

### `-r` Set rootDir

`rootDir` is a shortcut variable which can be used in setups and teardows link. Default is current directory.

### `-s` Panic mode

Cotton will terminate immidiately if a test failed.

### `-v` Version

Cotton prints out application version.
