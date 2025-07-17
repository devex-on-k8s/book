# Appointments

This application is part of the Min Salus system and provides the functionality for managing appointments. It's part of the project built in the [Developer Experience on Kubernetes](https://www.manning.com/books/developer-experience-on-kubernetes) book by [Mauricio Salatino](https://salaboy.com) and [Thomas Vitale](https://www.thomasvitale.com).

## Requirements

- [httpie](https://httpie.io/cli), used for testing the HTTP API.

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

## Run

First, run the services the application depends on:

```shell script
podman compose up -d
```

Then, run the application:

```shell script
go run appointments.go
```

The application will start on port `8081` by default and the process will keep running. When you're done, stop the application process with `Ctrl+C` and then stop the dependent services:

```shell script
podman compose down
```

## Test

First, run the services the application depends on:

```shell script
podman compose up -d
```

Run all unit and integration tests:

```shell script
go test
```

When you're done, stop the dependent services:

```shell script
podman compose down
```
