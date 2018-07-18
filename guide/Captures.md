# Captures

Captures is about to collect some data from the action and pass them to the future actions in the same markdown.

Captures section is declared by using `## Capture`.

The list of the captures are defined in the 2-columned table with title `Name` and `Value`.

Name of the capture is the variable name which can be referred in future actions.

Values is the name of item in response like declared in [Assertions](Assertions.md).

## Example

```
## Capture

| Name | Value |
| - | - |
| new-todo-uri | Header.Location |
```

| Previous | Index | Next |
| :-: | :-: | :-: |
| [Test Teardowns](TestTeardowns.md) | [Top](README.md) | [Comments](Comments.md) |