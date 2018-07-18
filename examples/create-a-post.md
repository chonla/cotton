# Create A Post

## Preconditions

* [Login With Valid Credential](common/login.md)

## POST /posts

| Header | Value |
| - | - |
| Authorization | Bearer {token} |
| Content-Type | application/json |

```
{
	"title": "New Title 10",
	"body": "Text body",
	"featured_image": 0
}
```

## Expectation

New post should be created.

| Assert | Expected |
| - | - |
| StatusCode | 201 |
| Header.Location | /^/posts/\d+$/

## Capture

| Name | Value |
| - | - |
| new-post-uri | header.location |

## Finally

* [Delete A Newly Created Post](common/delete-the-new-post.md)