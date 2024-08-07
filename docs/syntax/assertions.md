---
title: Assertions
layout: default
nav_order: 3
permalink: /syntax/assertions
parent: Syntax
---

# Assertions

Once a response has been returned from the request, an assertion will be evaluated to verify if value in the response satisfies the requirements.

The assertion is written as a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
* `selector` assertion_operator `expected_value`
{% endhighlight %}

* Selector is a value selector, locating a value in the response. See [Selectors](./selectors.md)
* Assertion operator is an assertion operation, which can be unary operator or binary operator. The binary operator requires an expected value as the operand. See [Assertion Operators](#assertion-operators) below.
* Expected value is an optional operand, depending on operator.

## Assertion Operators

### is defined

`is defined` is an operator to check if the value is defined in response.

**Example**

{% highlight markdown %}
* `Body.message.id` is defined
{% endhighlight %}

### is undefined

`is undefined` is an operator to check if the value is not defined in response.

**Example**

{% highlight markdown %}
* `Body.message.error` is undefined
{% endhighlight %}

### == (Equal)

`==` is an operator to check if the value is equal to the expected value.

Value types for this assertion can be a number, string, or [regular expression](#regular-expression-assertion).

**Example**

{% highlight markdown %}
* `StatusCode` == `200`
* `StatusText` == `"OK"`
* `Headers.Content-Type` == `"application/json"`
{% endhighlight %}

### != (Not Equal)

`!=` is an operator to check if the value is different from the expected value.

Value types for this assertion can be a number, string, or [regular expression](#regular-expression-assertion).

**Example**

{% highlight markdown %}
* `StatusCode` != `404`
* `StatusText` != `"NOT FOUND"`
* `Headers.Content-Type` != `"text/plain"`
{% endhighlight %}

### > (Greater Than)

`>` is an operator to check if the value is greater than the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `Body.id` > `0`
{% endhighlight %}

### >= (Greater Than or Equal To)

`>=` is an operator to check if the value is greater than or equal to the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `Body.id` >= `0`
{% endhighlight %}

### < (Less Than)

`<` is an operator to check if the value is less than the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `Body.id` < `30`
{% endhighlight %}

### <= (Less Than or Equal To)

`<=` is an operator to check if the value is less than or equal to the expected value.

Value types for this assertion must be a number.

**Example**

{% highlight markdown %}
* `Body.id` <= `30`
{% endhighlight %}

## Regular Expression Assertion

Regular expression is used in equality and difference assertion. The pattern is enclosed between slashes like in Javascript.

{: .info }
The value to be asserted with regular expression assertion **MUST** be string.

**Example**

{% highlight markdown %}
* `StatusText` == `/^NOT.*/`
{% endhighlight %}
