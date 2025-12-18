---
title: Variables
layout: default
nav_order: 6
permalink: /syntax/variables
parent: Syntax
---

# Variables

A variable is named value which can be reused. There are 2 ways to define a variable, explicitly define a variable with value or capture a value from response into a variable.

{: .info }
The variables which defined in earlier step, no matter in setup or testcase or teardown, it will be reusable in later steps.

## Defining a variable

A variable definition is written in a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
* `variable_name` = `value`
{% endhighlight %}

If the value is a string, enclose it with double quotes. Otherwise, it will be treated as a number. If the value cannot be parsed as a number, it will be assumed to be a string.

**Example**

{% highlight markdown %}
* `id` = `3`
* `keyword` = `"value"`
{% endhighlight %}

{: .warning }
Missing variables in requests (e.g., using `{{ undefined_var }}`) will cause a ParseError and abort the spec execution.

## Capturing a value into a variable

Declare the variable before the request with a JSON path (starting with `$`). After the request completes, the value will be captured from the JSON response body.

**Example**

{% highlight markdown %}
* `token` = `$.data.access_token`

```http
POST https://api.example.com/auth HTTP/1.1
Content-Type: application/json

{"username":"user","password":"pass"}
```
{% endhighlight %}

See [Captures](./captures.md) for more details.

## Using a variable

The variable can be used within the request code block by putting variable name enclosed with {% raw %}`{{`{% endraw %} and {% raw %}`}}`{% endraw %} (double curly braces).

**Example**

{% highlight markdown %}
* `username` = `"peter"`
* `password` = `"1234test"`

```http
POST https://fakestoreapi.com/auth/login HTTP/1.1
Content-Type: application/json
Content-Length: 42

{% raw %}{"username":"{{username}}","password":"{{password}}"}{% endraw %}
```
{% endhighlight %}

The actual request which sent to server will look like this:

{% highlight http %}
POST https://fakestoreapi.com/auth/login HTTP/1.1
Content-Type: application/json
Content-Length: 42

{"username":"peter","password":"1234test"}
{% endhighlight %}