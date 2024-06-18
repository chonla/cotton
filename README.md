# Cotton

Cotton is a markdown-formatted API specification runner.

## Version 0.x compatibility

This version contains several breaking changes and is not compatible with previous version.

## Markdown as an API specification

Cotton allows users to write tests in more intuitive and readable way.

## Specification structure

* <u>**Test case title**</u> will be picked from <u>**the first heading level 1**</u>.
* <u>**Test case description**</u> is right after <u>**the test case title**</u>.
* <u>**Test case request**</u> is the request in <u>**the first http-language code block**</u>.
* An <u>**unordered list**</u> with links to executable files -- if it is written <u>**before the test case request**</u>, it is test setups. if it is written <u>**after the test case request**</u>, it is test teardowns.
* <u>**Assertions**</u> is also written in an <u>**unordered list**</u> in the format `` `selector` <operator> `expected-value` ``.
* <u>**Captures**</u> is also written in an <u>**unordered list**</u> in the format `` `name` : `selector` ``.
* Other than things listed above will be considered as documentation. This embraces readability to everyone.

## Examples

See the files in `etc/examples` directory.

## TO DO

* ~~Three-tilded code block parsing~~
* ~~More assertion operators~~
  * ~~== (Equal to)~~
  * ~~&gt; (Greater than)~~
  * ~~&lt; (Less than)~~
  * ~~!=~~ (Not equal to)
  * ~~&gt;= (Greater than or equal to)~~
  * ~~&lt;= (Less than or equal to)~~
  * ~~is undefined (or not present)~~
  * ~~is defined (or present)~~
* ~~More readable logging~~
* Regular expression assertion operands
  * == (Match)
  * != (Not match)
* Captures integration
* ~~Debug mode~~
* CLI
* Cookbook
* ~~Colorized logging~~
* HTML reports