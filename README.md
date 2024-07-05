# Cotton

Cotton is a markdown-formatted API specification runner. Cotton promotes the readability and understandability of API specification.

## Syntax, Guide and Examples

See [Guidebook](https://chonla.github.io/cotton) for detail.

## Key Features

* **Customizable Documentation**: Cotton offers intuitive syntax. You can make your API document more readable.
* **Setups and Teardowns**: Cotton allows you to to test setups and teardown.
* **Variables**: Cotton allows you to define initial variables and use them later.
* **Captures**: Cotton offers value capturing into variable, which help you to perform tests with dynamic data.
* **Regular Expression Assertions**: Not only simple assertions, but Cotton also provides built-in regular expression assertions.
* **Debug Logging**: Cotton CLI has a `-d` option to print out detailed information of test operations.
* **CI/CD Ready**: Cotton CLI returns error if any test fails.

## Backward Compatibility

Cotton 1 contains several breaking changes and is not compatible with Cotton 0.

## To Do List

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
* ~~Regular expression assertion operands~~
  * ~~== (Match)~~
  * ~~!= (Not match)~~
* ~~Captures integration~~
* ~~Debug mode~~
* ~~CLI~~
* ~~Cookbook~~
* ~~Colorized logging~~
* HTML reports
* Upload request
* ~~Insecure mode~~
* ~~Track time taken~~
* ~~Stop tests execution when failed~~
* ~~Variables~~
* ~~Ordered list support for Setups,Teardowns,Assertions,Captures~~
* ~~Customized rootDir~~

## Contributing

1. Fork it (https://github.com/chonla/cotton/fork).
2. Create your feature branch (git checkout -b feature/fooBar).
3. Commit your changes (git commit -am 'Add some fooBar').
4. Push to the branch (git push origin feature/fooBar).
5. Create a new Pull Request.

## License

[MIT](./LICENSE)