---
title: Testcases
layout: default
nav_order: 1
permalink: /guide/testcases
parent: Guide
---

# Testcase

Testcase is a set of actions to test if a system meets the API requirement.

## A Minimal Testcase

A testcase requires two testcase elements, [request](../syntax/requests.md) and [assertions](../syntax/assertions.md). A testcase without title will be titled `Untitled`.

**Example**

{% highlight markdown %}
```http
GET https://fakestoreapi.com/products HTTP/1.1
```

* `Body.0.id`==`1`
{% endhighlight %}

## A Detailed Testcase

A testcase has 2 optional detail, title and description. The testcase title is the very first heading level 1 in the testcase file and the description is a paragraph right after the title. Other than that and not the testcase elements will be treated as documentation of the testcase.

{% highlight markdown %}
# List product

Listing product should return a list of product.

```http
GET https://fakestoreapi.com/products HTTP/1.1
```

* `Body.0.id`==`1`
{% endhighlight %}
