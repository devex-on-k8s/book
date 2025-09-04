# Appointments

This application is part of the Min Salus system and provides the functionality for managing appointments. It's part of the project built in the [Developer Experience on Kubernetes](https://www.manning.com/books/developer-experience-on-kubernetes) book by [Mauricio Salatino](https://salaboy.com) and [Thomas Vitale](https://www.thomasvitale.com).

## HTTP API

| Endpoint	      | Method   | Req. body   | Status | Resp. body     | Description    		   	             |
|:---------------:|:--------:|:-----------:|:------:|:--------------:|:--------------------------------------|
| `/`             | `GET`    |             | 200    | String         | Welcome message.                      |
| `/appointments` | `GET`    |             | 200    | Appointment[]  | Get all the booked appointments.      |
| `/appointments` | `POST`   | Appointment | 201    | Appointment    | Book a new appointment.               |
| `/appointments` | `DELETE` |             | 204    |                | Delete all appointments.              |
| `/chat`         | `POST`   | Question    | 200    | String         | Ask a question about the appointments.|

## Activate Environment

With Devbox:

```shell script
devbox shell
```

With Devcontainers, open the project in VS Code and select "Reopen in Container" from the Command Palette (`Ctrl+Shift+P` or `Cmd+Shift+P` on macOS).

## Run (Java)

If you have Ollama up and running natively on your machine, run the application in development mode, with live reload:

```shell script
./gradlew bootRun
```

If you'd like the application to take care of provisioning an Ollama instance, you can run it as follows:

```shell script
./gradlew bootRun -Darconia.dev.services.ollama.enabled=true
```

The application will start on port `8081` by default and the process will keep running, watching for changes in the source code.

Book a series of appointments:

```shell script
http :8081/appointments patientId=42 category="cardiology" appointmentDate="2028-02-28T12:00:00Z"

http :8081/appointments patientId=21 category="reumatology" appointmentDate="2029-02-28T11:00:00Z"

http :8081/appointments patientId=73 category="cardiology" appointmentDate="2039-02-27T15:00:00Z"
```

Call the chat endpoint to ask a question about the appointments:

```shell script
http :8081/chat question="How many patients have booked cardiology appointments? Just tell me the number and their IDs."
```

When you're done, stop the application process with `Ctrl+C`.

## Run (Kubernetes)

Run the application in development mode on Kubernetes, with live reload:

```shell script
skaffold dev --port-forward
```

The application will start on port `8081` by default and the process will keep running, watching for changes in the source code.

Book a series of appointments:

```shell script
http :8081/appointments patientId=42 category="cardiology" appointmentDate="2028-02-28T12:00:00Z"

http :8081/appointments patientId=21 category="reumatology" appointmentDate="2029-02-28T11:00:00Z"

http :8081/appointments patientId=73 category="cardiology" appointmentDate="2039-02-27T15:00:00Z"
```

Call the chat endpoint to ask a question about the appointments:

```shell script
http :8081/chat question="How many patients have booked cardiology appointments? Just tell me the number and their IDs."
```

When you're done, stop the application process with `Ctrl+C`.

## Clean

Clean the build directory:

```shell script
./gradlew clean
```
