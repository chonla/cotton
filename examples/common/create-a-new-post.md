# Create A New Post

## POST /posts

| Header | Value |
| - | - |
| Content-Type | application/json |
| Authorization | Bearer {token} |

```
{
	"title": "New Post",
	"body": "Text body",
	"featured_image": 0
}
```

## Capture

| Name | Value |
| - | - |
| new-post-uri | header.location |