# Micro-Gopia Golang API

This is a Golang API project that uses the Gin framework to handle HTTP requests and responses. The project reads a configuration file and uses it to set various settings. The configurations are loaded using the Viper library. If the project is run with the command "migrate", it executes a database migration.

## Installation

To use this project, you'll need Golang installed on your computer. Follow these steps:

1. Clone the repository from GitHub: `git clone https://github.com/ali-shokoohi/micro-gopia.git`
2. Navigate to the cloned repository: `cd micro-gopia`
3. Install the dependencies: `go get ./...`

## Configuration

The configuration settings for this project are stored in a YAML file located at `config/config-debug.yaml`. If the environment variable `GIN_MODE` is set to "release", it loads the configuration file at `config/config.yaml` and sets the debug flag to false. Below is a sample configuration file:

```yaml
service:

  http:

    host: "localhost"

    port: "8080"

  db:

    host: "localhost"

    port: "5432"

    dbname: "micro-gopia"

    user: "user"

    password: "password"
```

Make sure to replace the values with your own settings.

## Usage

To start the server, run the following command:

`go run cmd/app/main.go`


You can access the API at `http://localhost:8080/api/v1`.

## API Routes

The following routes are available for this API:

- `/api/v1`: The API's home page.
	- `/`: The home page's route.
	- `/users`: Routes for interacting with users.

These routes are defined in the `internal/api/routes/routes.go` file. The routes are created by nesting groups of routes using the Gin framework's `Group` method.

## Database Migrations

To run a database migration, use the following command:

`go run cmd/app/main.go migrate`


This will execute the database migration code defined in the `pkg/migrations` package.

## Contributing

If you wish to contribute to this project, feel free to submit a pull request or open an issue.
