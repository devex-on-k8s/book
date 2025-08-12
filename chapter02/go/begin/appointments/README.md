# Appointments

This application is part of the Min Salus system and provides the functionality for managing appointments. It's part of the project built in the [Developer Experience on Kubernetes](https://www.manning.com/books/developer-experience-on-kubernetes) book by [Mauricio Salatino](https://salaboy.com) and [Thomas Vitale](https://www.thomasvitale.com).

## Requirements

- [httpie](https://httpie.io/cli), used for testing the HTTP API.
- [Air](https://github.com/air-verse/air) for hot reloading the application on changes.
- [golangci-lint](https://golangci-lint.run/) for linting the code.

## HTTP API

| Endpoint	      | Method   | Req. body   | Status | Resp. body     | Description    		   	              |
|:---------------:|:--------:|:-----------:|:------:|:--------------:|:-------------------------------------|
| `/`             | `GET`    |             | 200    | String         | Welcome message.                     |
| `/appointments` | `GET`    |             | 200    | Appointment[]  | Get all the booked appointments.     |
| `/appointments` | `POST`   | Appointment | 201    | Appointment    | Book a new appointment.              |
| `/appointments` | `DELETE` |             | 204    |                | Delete all appointments.             |

Get the welcome message:

```shell script
http :8081
```

Book an appointment:

```shell script
http :8081/appointments patientId=42 appointmentDate="2028-02-29T12:00:00Z"
```

Get all appointments:

```shell script
http :8081/appointments
```

Delete all appointments:

```shell script
http DELETE :8081/appointments
```

## Lint

Run the linter:

```shell script
golangci-lint run --fix -v
```

## Run

Thanks to `GoFiber`, `Testcontainers Go` and `air`, you can run the application in a containerized environment, so there is no need to install any database or other dependencies.

Just run the following command to start the application:

```shell script
air
```

The application will start on port `8081` by default and the process will keep running. When you're done, stop the application process with `Ctrl+C` and the dependent services will be automatically stopped for you.

## Test

Run all unit and integration tests:

```shell script
go test -tags dev ./...
```

The `-tags dev` flag is used to add to the build those files that are only used for development purposes. This application uses it to include the `config_dev.go` file, which is used to run the dependencies in a containerized environment.

Thanks to `Testcontainers Go`, the started services will be automatically stopped for you.

To know more about the different technologies used in this application, please refer to the following links:

- [GoFiber](https://gofiber.io/)
- [GoFiber's Services](https://docs.gofiber.io/next/api/services)
- [Testcontainers Go](https://golang.testcontainers.org/)
- [Testcontainers GoFiber's Contrib](https://docs.gofiber.io/contrib/testcontainers_v0.x.x/)
