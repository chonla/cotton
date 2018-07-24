# Cotton

Markdown Test Specification Runner. This project is originally inspired by [silk test](https://github.com/matryer/silk).

* Executable RESTful API Markdown-based Document Specification.

## Installation

### From source

```
go get github.com/chonla/cotton
```

## Usage

```
cotton [-u <base-api-url>] [-i] [-d] <directory or file>
```

```
$ cotton

Usage of cotton:

  cotton [-u <base-url>] <test-cases>

  test-cases can be a markdown file or a directory contain markdowns.

  -h  show this help
  -i  insecure mode -- to disable certificate verification
  -d  detail mode -- to dump test detail
  -u string
      set base url (default "http://localhost:8080")
  -v  show cotton version
```

## Executable markdown specfication

See [Guide](./guide) for more information.

## License

[License](LICENSE.txt)