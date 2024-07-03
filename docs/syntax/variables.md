---
title: Variables
layout: default
nav_order: 6
permalink: /syntax/variables
parent: Syntax
---

# Variables

A variable is named value which can be reused. There are 2 ways to define a variable, explicitly define a variable with value or capture a value from response into a variable.

## Defining a variable

A variable definition is written in a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
* variable_name:value
{% endhighlight %}

If the value is string, enclose it with double quote. Otherwise, it will be treated as a number. Anyway, if the value cannot be parsed into a number, it will be assumed to be a string.

**Example**

{% highlight markdown %}
* id:3
* keyword:"value"
{% endhighlight %}
