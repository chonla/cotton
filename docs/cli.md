---
layout: default
title: CLI
nav_order: 3
permalink: /cli
---

# Command Line Interface (CLI)

## Usage

{% highlight bash %}
  cotton [file|dir|glob ...] [-d] [-q] [-m] [-f] [-v] [--timeout <ms>] [--json-report <file>]
  cotton -v
  cotton --help

  -d, --debug
        enable debug logging
  -q, --quiet
        suppress output (exit code only)
  -m, --monochrome
        disable colored output
  -f, --fail-fast
        stop on first failing test
  -v, --version
        display cotton version
  --timeout <ms>
        set HTTP timeout (1–60000 ms; default 30000)
  --json-report <file>
        generate JSON report to specified file
{% endhighlight %}

## Options

By default, Cotton prints out test title, setup title, teardown title, result, and summary.

### `-d, --debug`

Debug mode enables verbose logging. Cotton prints out setup/teardown titles, request/response bodies, variable resolutions, and expectation evaluations.

### `-q, --quiet`

Quiet mode suppresses all output and prints only the exit code (0 for success, non-zero for failure). Useful for CI/CD pipelines.

### `-m, --monochrome`

Monochrome mode disables colored output. Useful for terminals that don't support colors.

### `-f, --fail-fast`

Fail-fast mode stops test execution on the first failing test. Exit code will be non-zero if any test failed.

### `-v, --version`

Displays the Cotton version and exits with code 0.

### `--timeout <ms>`

Sets the HTTP timeout in milliseconds. Valid range: 1–60000; default is 30000. Values exceeding 60000 are capped at 60000. If an invalid timeout is provided, Cotton exits with a ParseError.

**Example**

{% highlight bash %}
cotton ./specs --timeout 45000
{% endhighlight %}

### `--json-report <file>`

Generates a JSON report with test results, timing, and summary to the specified file.

**Example**

{% highlight bash %}
cotton ./specs --json-report results.json
{% endhighlight %}
