# Test Teardowns

Test teardown is a list of actions like test setup under `## Finally` section. The actions in test teardown also do not require assertions.

To declare an action to be done in test teardown, use anchor `[Title](Link)`.

## Example

```
## Finally

* [Remove a todo item](RemoveTodo.md)
* [Logout](Logout.md)
```

| Previous | Index | Next |
| :-: | :-: | :-: |
| [Test Setups](TestSetups.md) | [Top](README.md) | [Captures](Captures.md) |