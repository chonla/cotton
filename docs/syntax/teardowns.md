---
title: Teardowns
layout: default
nav_order: 5
permalink: /syntax/teardowns
parent: Syntax
---

# Teardowns

Teardown is a request operation that will be executed after the testcase.

Like [setups](./setups.md), each teardown must be a separated markdown file. If you have several teardown steps, just separate them one step one file. This allows the teardown operation to be able to be reused in other testcases.

The teardown is written as a link in a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
* [teardown_title](link_to_teardown_file)
{% endhighlight %}

{: .info }
The teardown list can be anywhere in the testcase, but must be **AFTER** the request.

**Example**

{% highlight markdown %}
* [Session clean up](shared/authsession_clean_up.md)
{% endhighlight %}
