# Example: JSONPlaceholder API

[Get all posts](example-setup.spec.md)

```
GET https://jsonplaceholder.typicode.com/posts/{{ post_id }} HTTP/1.1
Accept: application/json
```

- `$.id` > `0`
- `$.userId` == `1`
- `$.title` != `""`
