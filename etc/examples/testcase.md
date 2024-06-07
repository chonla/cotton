# This is title of test case written with ATX Heading 1

The test case is described by providing paragraphs right after the test case title.

The description of test case can be single or multiple lines.

Cotton will consider only the first ATX Heading 1 as the test title.

## Before Test

All links defined as list before code block ` ```http` will be executed before test request.

* [Link before the test will be executed before executing test](../../etc/examples/executable_before.md)

## Request

HTTP request is described by code block with `http` annotation right after the open code block ` ```http `.

Cotton will consider only the first code block as the testing request.

```http
POST /some-path HTTP/1.1
Host: localhost

{
    "login": "login_name"
}
```

## After Test

All links defined as list after code block ` ```http` will be executed after test request.

* [Link after the test will be executed after executing test](../../etc/examples/executable_after.md)