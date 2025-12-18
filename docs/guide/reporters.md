---
title: Reporters
layout: default
nav_order: 4
permalink: /guide/reporters
parent: Guide
---

# Report Generation

Cotton can generate test reports in JSON format by specifying the output file with the `--json-report` flag. See [CLI](../cli.md) for details.

## JSON Report

JSON reports contain test results, timing information, and a summary. Use the `--json-report` flag to generate a report:

{% highlight bash %}
cotton ./specs --json-report results.json
{% endhighlight %}

The report will be saved to the specified file in JSON format.
