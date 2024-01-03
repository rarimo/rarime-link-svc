# Rarime Link Service

Rarime Link Service stores users proofs and provide simple access to it via REST API in [JSON:API] format.

## Description

The service includes 1 main runner:<br>
```proofs-cleaner```: removes expired proofs from the database.

## Getting Started
### Prerequisites

Before you begin, ensure you have met the following requirements:

- Docker installed, see [Docker installation guide]
- Go 1.20 installed, see [Go installation guide]


### Building

#### Binary
To build the service binary file, follow these steps:

1. Clone the repository.

    ```bash
    git clone github.com/rarimo/rarime-link-svc
    cd rarime-link-svc
    ```

1. Install dependencies and build the service.

    ```bash
    go mod tidy
    go build main.go
    ```

#### Docker

To build the service Docker image, follow these steps:


1. Clone the repository.

    ```bash
    git clone github.com/rarimo/rarime-link-svc
    cd rarime-link-svc
    ```

1. Build the service image.

    ```bash
    sh ./build.sh
    ```

### Configuration

To properly configure the service, provide valid config file, see [config-example.yaml](config-example.yaml)
for example.

### Running with Docker

To run the service using Docker, follow these steps:

1. Build the service image, see [Building](#building).
1. Run the service image.

    ```bash
    docker-compose up -d
    ```
1. The service will be available on the `8000` port.

## Usage

To use the service, you could use the swagger documentation, see [API Documentation](#api-documentation), or
any other http client.

## API Documentation

We use [openapi:json] standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor], here is how you can start it
```bash
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.

## Contributing

We welcome contributions from the community! To contribute to this project, follow these steps:

1. Fork the repository.
1. Create a new branch with a descriptive name for your feature or bug fix.
1. Make your changes and commit them.
1. Push your changes to your branch on your GitHub fork.
1. Create a pull request from your branch to the `main` branch of this repository.

Please ensure your pull request adheres to the following guidelines:
- Add a clear pull request title;
- Add a comprehensive pull request description that includes the motivation behind the changes, steps needed to test them, etc;
- Update the [CHANGELOG.md](CHANGELOG) accordingly;
- Keep the codebase clean and well-documented;
- Make sure your code is properly tested;
- Reference any related issues in your pull request;

The maintainers will review your pull request and may request changes or provide feedback before merging. We appreciate your contributions!

## Changelog

For the changelog, see [CHANGELOG.md](CHANGELOG).

## License

This project is under the MIT License - see the [LICENSE](LICENSE) file for details.

[JSON:API]: https://jsonapi.org/
[Docker installation guide]: https://docs.docker.com/get-docker/
[Go installation guide]: https://golang.org/doc/install
[swagger editor]: http://localhost:8080/swagger-editor/
[openapi:json]: https://www.openapis.org/
