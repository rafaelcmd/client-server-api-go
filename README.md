# Go Client-Server Application

This is a basic example of a client-server application written in Go. The application consists of two parts: the server and the client.

## Server

### Prerequisites

- Go (Golang) installed on your machine.

### How to Run

1. Open a terminal and navigate to the "server" folder:

    ```bash
    cd path/to/your/app/server
    ```

2. Run the server:

    ```bash
    go run main.go
    ```

3. The server will start and wait for incoming connections.

## Client

### Prerequisites

- Go (Golang) installed on your machine.

### How to Run

1. Open another terminal and navigate to the "client" folder:

    ```bash
    cd path/to/your/app/client
    ```

2. Run the client:

    ```bash
    go run main.go
    ```

3. The client will connect to the server and generate the dollar exchange rate.

## Notes

- The server listens on `127.0.0.1:8080` by default. You can modify the code to use a different address or port if needed.

- This is a basic example and serves as a starting point. Feel free to modify and extend the code based on your specific requirements.

## License

This project is licensed under the [MIT License](LICENSE).