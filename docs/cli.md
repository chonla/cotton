---
layout: default
title: CLI
nav_order: 2
permalink: /cli
---

# Command Line Interface (CLI)

## Usage

{% highlight bash %}
  cotton [-d] [-c] <testpath|testdir>
  cotton -v
  cotton --help

  -c    compact mode
  -d    debug mode
  -h    display this help
  -i    disable certificate verification
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

### `-s` Panic mode

Cotton will terminate immidiately if a test failed.

### `-v` Version

Cotton prints out application version.
