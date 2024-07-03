---
title: Assertions
layout: default
nav_order: 3
permalink: /syntax/assertions
parent: Syntax
---

# Assertions

An assertion is a test to verify if requirements are satisfied.

The values in response are addressable by [selector](./selectors.md).

## Assertion operators

### is defined

`is defined` is an operator to check if the value is defined in response.

**Example**

{% highlight markdown %}
```markdown
* `Body.message.id` is defined
```
{% endhighlight %}

### is undefined

`is undefined` is an operator to check if the value is not defined in response.

**Example**

{% highlight markdown %}
```markdown
* `Body.message.error` is undefined
```
{% endhighlight %}

### == (Equality)

`==` is an operator to check if the value is equal to the expected value.

**Example**

{% highlight markdown %}
```markdown
* `StatusCode` == `200`
* `StatusText` == `"OK"`
* `Headers.Content-Type` == `"application/json"`
```
{% endhighlight %}