# Cotton

Markdown Test Specification Runner. This project is originally inspired by [silk test](https://github.com/matryer/silk).

* Executable RESTful API Markdown-based Document Specification.

## Installation

```
go get github.com/chonla/cotton
```

## Usage

```
cotton -u <base-api-url> [directory or file]
```

```
$ cotton

Usage of cotton:

  cotton [-u <base-url>] <test-cases>

  test-cases can be a markdown file or a directory contain markdowns.

  -h	show this help
  -u string
    	set base url (default "http://localhost:8080")
  -v	show cotton version
```

## Markdown specfication

See [Guide](./guide) for more information.

## License

[License](LICENSE.txt)