# Appointments Service

This application is part of the Min Salus system and provides the functionality for managing appointments. It's part of the project built in the [Developer Experience on Kubernetes](https://www.manning.com/books/developer-experience-on-kubernetes) book by [Mauricio Salatino](https://salaboy.com) and [Thomas Vitale](https://www.thomasvitale.com).

## HTTP API

| Endpoint	      | Method   | Req. body   | Status | Resp. body     | Description    		   	             |
|:---------------:|:--------:|:-----------:|:------:|:--------------:|:--------------------------------------|
| `/`             | `GET`    |             | 200    | String         | Welcome message.                      |
| `/appointments` | `GET`    |             | 200    | Appointment[]  | Get all the booked appointments.      |
| `/appointments` | `POST`   | Appointment | 201    | Appointment    | Book a new appointment.               |
| `/appointments` | `DELETE` |             | 204    |                | Delete all appointments.              |
| `/chat`         | `POST`   | ChatRequest | 200    | ChatResponse   | Ask a question about the appointments.|

## Activate Environment

With Devbox:

```shell script
devbox shell
```

With Devcontainers, open the project in VS Code and select "Reopen in Container" from the Command Palette (`Ctrl+Shift+P` or `Cmd+Shift+P` on macOS).

## Run (Go)

First, ensure you have a PostgreSQL database running. Run this command:

```shell script
podman compose up -d db
```

Ensure you have Ollama running on your machine and then run the application as follows:

```shell script
go run appointments.go
```

You can now call the `/chat` endpoint and ask questions that are forward to the LLM served by Ollama.

```shell script
http :8081/chat question="What's the capital of Denmark?"
```

However, you'll notice that the LLM doesn't have any knowledge about the appointments booked in the system.

```shell script
http :8081/chat question="How many appointments are booked in the system?"
```

To stop the application, press `Ctrl+C`.

Finally, stop the PostgreSQL database:

```shell script
podman compose down
```

## Run (Kubernetes)

Run the application in development mode on Kubernetes:

```shell script
skaffold dev --port-forward
```

The application will start on port `8081` by default and the process will keep running, watching for changes in the source code.

You can now call the `/chat` endpoint and ask questions that are forward to the LLM served by Ollama.

```shell script
http :8081/chat question="What's the capital of Denmark?"
```

However, you'll notice that the LLM doesn't have any knowledge about the appointments booked in the system.

```shell script
http :8081/chat question="How many appointments are booked in the system?"
```

When you're done, stop the application process with `Ctrl+C`.
