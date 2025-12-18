---
# Feel free to add content and custom Front Matter to this file.
# To modify the layout, see https://jekyllrb.com/docs/themes/#overriding-theme-defaults

layout: default
title: General
nav_order: 1
description: "Cotton is a markdown-formatted API specification runner."
permalink: /
---

# Cotton

![GitHub Repo stars](https://img.shields.io/github/stars/chonla/cotton) ![GitHub Release](https://img.shields.io/github/v/release/chonla/cotton)

Cotton is a markdown-formatted API specification runner. Cotton promotes the readability and understandability of API specification.

## Key Features

* **Customizable Documentation**: Cotton offers intuitive syntax. You can make your API document more readable.
* **Setups and Teardowns**: Cotton allows you to test setups and teardown with proper variable sharing across contexts.
* **Variables**: Cotton allows you to define initial variables and use them later, with support for capturing values from responses.
* **Captures**: Cotton offers value capturing into variables using JSON paths, which help you to perform tests with dynamic data.
* **Assertions**: Cotton provides multiple assertion operators including equality, comparison, and regular expression matching.
* **Debug Logging**: Cotton CLI has a `-d` option to print out detailed information of test operations.
* **Report Generation**: Cotton can generate JSON reports with test results and timing information.
* **CI/CD Ready**: Cotton CLI returns appropriate exit codes and supports multiple input formats (file, directory, or glob patterns).