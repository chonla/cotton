---
title: Setups
layout: default
nav_order: 4
permalink: /syntax/setups
parent: Syntax
---

# Setups

Setup is a request operation that will be executed before the testcase.

The setup must be a separated markdown file. If you have several setup steps, just separate them one step one file. This allows the setup operation to be able to be reused in other testcases.

The setup is written as a link in a list item (ordered or unordered list) in the following syntax:

{% highlight markdown %}
```markdown
* [setup_title](link_to_setup_file)
```
{% endhighlight %}

**Note**: A built-in variable `<rootDir>` can be used in link-to-setup-file to make link more readable.

**Example**

{% highlight markdown %}
```markdown
* [Authenticate with a support credential](<rootDir>/shared/auth_support_cred.md)
```
{% endhighlight %}
