# PackageSpy

PackageSpy is a versatile command-line tool designed to simplify the process of searching for secrets within packages on popular package managers using Gitleaks. It provides a convenient interface for security researchers, developers and system administrators to identify and manage sensitive information leaks across different environments.

## Installation

Before you start using PackageSpy, make sure you have Go (Golang) installed on your system. You can download and install Go from the official website: [Go Downloads](https://golang.org/dl/)

Once you have Go installed, you can install PackageSpy using the following command:

```shell
go install github.com/aydinnyunus/PackageSpy@latest
```

## Usage

PackageSpy supports four different search options, combining keyword and package manager:

1. Search for packages using a keyword on npm:
   ```shell
   go run . scan --search keyword --npm
   ```

2. Search for packages using a keyword on PyPI:
   ```shell
   go run . scan --search keyword --pypi
   ```

3. Search for packages by a user's username on npm:
   ```shell
   go run . scan --username username --npm
   ```

4. Search for packages by a user's username on PyPI:
   ```shell
   go run . scan --username username --pypi
   ```

Replace `keyword` with your desired search term and `username` with the username you want to search for.

## Example

Here's an example of using PackageSpy to search for Python packages related to data science on PyPI:

```shell
go run . scan --search datascience --pypi
```

## Features

- Cross-platform compatibility: PackageSpy is written in Go, making it compatible with Windows, macOS, and Linux.
- Seamless integration: Easily incorporate PackageSpy into your development workflow by using the provided CLI commands.
- Efficient searches: Quickly find packages related to your specific needs using either keywords or usernames on npm and PyPI.

## Contributing

PackageSpy is an open-source project, and we welcome contributions from the community. If you have ideas for improvements or would like to report issues, please visit our GitHub repository: [PackageSpy](https://github.com/aydinnyunus/PackageSpy)


## Contact

[<img target="_blank" src="https://img.icons8.com/bubbles/100/000000/linkedin.png" title="LinkedIn">](https://linkedin.com/in/yunus-ayd%C4%B1n-b9b01a18a/) [<img target="_blank" src="https://img.icons8.com/bubbles/100/000000/github.png" title="Github">](https://github.com/aydinnyunus/WhatsappBOT) [<img target="_blank" src="https://img.icons8.com/bubbles/100/000000/instagram-new.png" title="Instagram">](https://instagram.com/aydinyunus_/) [<img target="_blank" src="https://img.icons8.com/bubbles/100/000000/twitter-squared.png" title="LinkedIn">](https://twitter.com/aydinnyunuss)
