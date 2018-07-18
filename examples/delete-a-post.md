# Delete A Post

## Preconditions

* [Login With Valid Credential](common/login.md)
* Then [Create A New Post](common/create-a-new-post.md)

## DELETE {new-post-uri}

| Header | Value |
| - | - |
| Authorization | Bearer {token} |

## Expectation

New post should be deleted.

| Assert | Expected |
| - | - |
| StatusCode | 200 |
