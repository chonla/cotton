# Test executable

Only ` ```http` block and Capture in this file are used. Others will be considered as documentation.

```http
GET /get-info HTTP/1.1
Host: localhost
```

## Capture

Capturing data returned from recent request can be done by define them in `name:selector` syntax as a list like below. Backticking in selector is optional.

* readiness:`$.readiness`
* version:`$.version`