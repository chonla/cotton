---
title: Testcases
layout: default
nav_order: 1
permalink: /guide/testcases
parent: Guide
---

# Testcase

Testcase is a set of actions to test if a system meets the API requirement.

## Minimal Testcase

A testcase requires two testcase elements, [request](../syntax/requests.md) and [assertions](../syntax/assertions.md).

**Example**

{% highlight markdown %}
```http
POST https://somedomain.com/resource HTTP/1.1
Content-Length: 17
Content-Type: application/json

{"resource":8839}
```

* `StatusText` == `/^NOT.*/`
{% endhighlight %}
