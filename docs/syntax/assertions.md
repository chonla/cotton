---
title: Assertions
layout: default
nav_order: 3
permalink: /syntax/assertions
parent: Syntax
---

# Assertions

Once a response has been returned from the request, assertions are evaluated against the JSON response body to verify expected values.

The assertion is written as a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
* `$.json.path` operator `expected_value`
{% endhighlight %}

* Selector: a JSON path starting with `$`. See [Selectors](./selectors.md).
* Operator: comparison or regex operator. See [Assertion Operators](#assertion-operators).
* Expected value: number, string (quoted), regex (`/pattern/`), `null`, or `undefined`.

## Assertion Operators

> Note: Unary operators like "is defined"/"is undefined" are not supported. Use comparisons with `undefined`/`null` as shown below.

### == (Equal)

`==` is an operator to check if the value is equal to the expected value.

Value types for this assertion can be a number, string, or [regular expression](#regular-expression-assertion).

**Example**

{% highlight markdown %}
* `$.status` == `"ok"`
* `$.data.id` == `42`
{% endhighlight %}

### != (Not Equal)

`!=` is an operator to check if the value is different from the expected value.

Value types for this assertion can be a number, string, or [regular expression](#regular-expression-assertion).

**Example**

{% highlight markdown %}
* `$.data.name` != `"Guest"`
{% endhighlight %}

### > (Greater Than)

`>` is an operator to check if the value is greater than the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `$.metrics.duration_ms` > `0`
{% endhighlight %}

### >= (Greater Than or Equal To)

`>=` is an operator to check if the value is greater than or equal to the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `$.metrics.duration_ms` >= `100`
{% endhighlight %}

### < (Less Than)

`<` is an operator to check if the value is less than the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `$.data.count` < `30`
{% endhighlight %}

### <= (Less Than or Equal To)

`<=` is an operator to check if the value is less than or equal to the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `$.data.count` <= `30`
{% endhighlight %}

### =~ (Regex Match)

`=~` is an operator to check if the value matches the regex pattern.

The pattern is enclosed between forward slashes like `/pattern/`.

**Example**

{% highlight markdown %}
* `$.title` =~ `/^Executive/`
* `$.department` =~ `/[A-Z]{2,}/`
{% endhighlight %}

### !~ (Regex Not Match)

`!~` is an operator to check if the value does **not** match the regex pattern.

**Example**

{% highlight markdown %}
* `$.title` !~ `/^Sales/`
{% endhighlight %}

## Regular Expression Assertion

Regular expressions are used in regex match/not-match assertions and can also be used in equality/difference assertions. The pattern is delimited by forward slashes.

{: .info }
Regex assertions operate on string values. Convert numbers in your API if needed.

Checking existence:

{% highlight markdown %}
* `$.data.error` == `undefined`   # not present
* `$.data.deleted_at` == `null`    # null or not present
* `$.data.name` != `undefined`     # present
{% endhighlight %}
