---
title: Testcases
layout: default
nav_order: 1
permalink: /guide/testcases
parent: Guide
---

# Testcase

Testcase is a set of actions to test if a system meets the requirement.

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

**Example**

{% highlight markdown %}
# List product

Listing product should return a list of product.

```http
GET https://fakestoreapi.com/products HTTP/1.1
```

* `Body.0.id`==`1`
* `Body.1.id`==`2`
{% endhighlight %}

## An Internationalized Testcase

Since the testcase embraces the markdown syntax, you can also use any language to make the testcase more readable.

**Example**

{% highlight markdown %}
# รายชื่อผลิตภัณฑ์ทั้งหมด

ทดสอบ API สำหรับลิสต์รายชื่อผลิตภัณฑ์ทั้งหมด โดยจะต้องทำการล็อกอินเข้าใช้งานเพื่อเอา token ก่อน และส่ง token ไปใน header เพื่อเรียกดูรายชื่อผลิตภัณฑ์

* [ลงชื่อเข้าใช้ระบบ](../executables/auth.md)

```http
GET https://fakestoreapi.com/products HTTP/1.1
Authorization: Bearer {{access_token}}
```

* `Body.0.id`==`1`
{% endhighlight %}