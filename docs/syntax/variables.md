---
title: Variables
layout: default
nav_order: 6
permalink: /syntax/variables
parent: Syntax
---

# Variables

A variable is named value which can be reused. There are 2 ways to define a variable, explicitly define a variable with value or capture a value from response into a variable.

The variables which defined in earlier step, no matter in setup or testcase or teardown, it will be reusable in later steps.

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

Defining a variable which has a duplicate name will replace the previous definition. The `keyword` in the following example yield the value of `"another value"`.

**Example**

{% highlight markdown %}
* id:3
* keyword:"value"
* keyword:"another value"
{% endhighlight %}

## Capturing a value into a variable

See [Captures](./captures.md).

## Using a variable

The variable can be used within the request code block by putting variable name enclosed with `{{` and `}}`.

**Example**

{% highlight markdown %}
* keyword:"value"

```http
POST https://somedomain.com/search HTTP/1.1
Content-Length: 17
Content-Type: application/json

{"keyword":"{{keyword}}"}
```
{% endhighlight %}

The actual request which sent to server will look like this:

{% highlight http %}
POST https://somedomain.com/search HTTP/1.1
Content-Length: 17
Content-Type: application/json

{"keyword":"value"}
{% endhighlight %}