# Test getting copyleft from opensource.org

## Setups

* [Setting up](<rootDir>/etc/examples/opensource.org/setup.md)

## API under test

```http
GET https://api.opensource.org/licenses/copyleft HTTP/1.1
```
## Setups

* [Tearing down](<rootDir>/etc/examples/opensource.org/teardown.md)

## Captures

* license:`Body.0.id`
* name2:`Body.0.name`

## Assertions

* `Body.1.id`==`"GPL-3.0"`
* `Body.1.id`!=`"GPL-4.0"`
* `Body.1.superseded_by`==`null`
* `Body.1.superseded_date` is undefined
* `Body.1.superseded_by` is defined
* `Body.1.id`==/^GPL/
* `Body.1.id`!=/^GpL/