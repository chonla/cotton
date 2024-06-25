# Creating a testcase

Testcase requires a request and at least 1 assertion.

## Request

Making a request is very simple. Just wrap your HTTP request with a code block annotated with `http`. The code block can be triple backticks or triple tildes.

You can place the request anywhere in the testcase file.

### Example

#### GET request

~~~markdown
```http
GET https://httpbin.org/get HTTP/1.1
```
~~~

#### POST request

~~~markdown
```http
POST https://httpbin.org/post HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 9

key=value
```
~~~

## Assertions

Assertion is an item look like these.

```markdown
* `Body.form.license`==`"GPL-3.0"`
* `Body.form.word`==`"world"`
```
