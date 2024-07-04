---
title: Captures
layout: default
nav_order: 7
permalink: /syntax/captures
parent: Syntax
---

# Captures

Capturing a value is one of two ways to define a [variable](./variables.md).

A capture is written in a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
* variable_name:`selector`
{% endhighlight %}

You can notice, the difference between capturing a value into a variable and variable definition is the selector for capture is enclosed with backticks, but the value is not.

{: .info }
The capture will capture the response within the same markdown file. You cannot capture the value from the other markdown files.