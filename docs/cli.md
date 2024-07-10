---
layout: default
title: CLI
nav_order: 3
permalink: /cli
---

# Command Line Interface (CLI)

## Usage

{% highlight bash %}
  cotton [-d] [-c] [-p] [-b <basedir>] <testpath|testdir>
  cotton -v
  cotton --help

  -b string
        set baseDir path
  -c    compact mode
  -d    debug mode
  -h    display this help
  -i    disable certificate verification
  -p    paranoid mode
  -s    stop when test failed
  -v    display cotton version
{% endhighlight %}

## Options

By default, Cotton prints out test title, setup title, teardown title, result, and summary.

### `-b` Set baseDir

All links in setups and teardows with relative path will be referenced from baseDir. Default is current directory.

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

### `-s` Panic mode

Cotton will terminate immidiately if a test failed.

### `-v` Version

Cotton prints out application version.
