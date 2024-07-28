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
* **Setups and Teardowns**: Cotton allows you to test setups and teardown.
* **Variables**: Cotton allows you to define initial variables and use them later.
* **Captures**: Cotton offers value capturing into variable, which help you to perform tests with dynamic data.
* **Regular Expression Assertions**: Not only simple assertions, but Cotton also provides built-in regular expression assertions.
* **Debug Logging**: Cotton CLI has a `-d` option to print out detailed information of test operations.
* **CI/CD Ready**: Cotton CLI returns error if any test fails.

## About Cotton

Cotton 0 was originally inspired by [silk test](https://github.com/matryer/silk) which is no longer maintained.

Testcase syntax in Cotton 0 is close to silk test, but in version 1, the syntax is completely redesigned to promote more readability.