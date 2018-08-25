# Cotton

Markdown Test Specification Runner. This project is originally inspired by [silk test](https://github.com/matryer/silk).

* Executable RESTful API Markdown-based Document Specification.

[![Latest stable version](https://img.shields.io/badge/stable-0.2.1-green.svg)](https://github.com/chonla/cotton/releases)

## Installation

### Homebrew/Linuxbrew

```sh
brew tap chonla/universe
brew install cotton
```

### From source

```sh
go get github.com/chonla/cotton
```

## Upgrade

```sh
brew upgrade
```

## Usage

```sh
cotton [-u <base-api-url>] [-i] [-d] <directory or file>
```

```sh
$ cotton

Usage of cotton:

  cotton [-u <base-url>] [-i] [-d] [-p name1=value1] [-p name2=value2] ... <test-cases>

  test-cases can be a markdown file or a directory contain markdowns.

  -d	detail mode -- to dump test detail
  -h	show this help
  -i	insecure mode -- to disable certificate verification
  -p value
    	to inject predefined variable in variable-name=variable-value format
  -u string
    	set base url (default "http://localhost:8080")
  -v	show cotton version
```

## Executable markdown specification

See [Guide](./guide) for more information.

## Contributing

1. Fork it (<https://github.com/chonla/cotton/fork>).
1. Create your feature branch (`git checkout -b feature/fooBar`).
1. Commit your changes (`git commit -am 'Add some fooBar'`).
1. Push to the branch (`git push origin feature/fooBar`).
1. Create a new Pull Request.

### Testing

```sh
go test ./...
```

## License

[MIT](LICENSE.txt)
