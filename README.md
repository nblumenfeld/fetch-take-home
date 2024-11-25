# Fetch Take Home - Noah Blumenfeld

### About the project
This is my submission for the take home project for Fetch! Overall, I thought this exercise was a fun and intersting way to do a first technical challenge. While I've primarily worked with JVM/Spring services, I decided to take this opportunity to tackle a problem using Go and learn some of the basics around how it feels to work with the langauge and it's related dependencies.

## Setup Instructions
### Required Dependencies
Please ensure you have at least the following installed:
- [Go 1.23.3](https://go.dev/doc/install)
- [Docker](https://www.docker.com/get-started/)

All data is saved in memory, so there are no other requirements than the above (and an active internet connection if running via docker).

### Running the Server without Docker
To run the server locally, all you need to do is run the following command within the root of the project:

```
$ go run .
```

This will run the service on port 8080, so please ensure that the port is free before starting.

### Running the Server with Docker
To run the server using docker, please do the following within the root of the project:

```
$ docker build --tag fetch-take-home .
$ docker run --publish 8080:8080 fetch-take-home
```

This will expose the service on port 8080, so please ensure that the port is free before starting.

