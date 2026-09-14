# Web Proxy with Content Control

Miguel Casarim da Silva  
Guilherme Estrella

This project consists of the implementation of an educational web proxy with content control, developed for the Internet Systems 2 course. The proxy acts as an intermediary between the client (browser) and the Internet, intercepting HTTP and HTTPS requests to forward, block, or filter content based on pre-configured rules.

## Why Did We Choose Go (Golang)?

The choice of technology was based on the following technical pillars of the Go language:

* **Absence of Heavy Frameworks:** Go's standard library (`net/http`) is extremely robust. It allowed us to build the HTTP server, manipulate requests, and manage TCP sockets natively, without needing to install third-party packages or frameworks.
* **Native Concurrency and Parallelism:** A proxy server needs to handle multiple simultaneous requests. Go solves this natively through **Goroutines**, which are functions executed concurrently. They consume a fraction of the memory compared to traditional operating system threads.
* **HTTPS Tunnel Efficiency:** For the implementation of the `CONNECT` method, we used Goroutines to read and copy bidirectional data traffic (client ↔ destination server) in real-time, ensuring high performance and preventing the server's main flow from blocking.

## Project Structure and File Logic

The project is modularized to clearly separate the responsibilities of configuration, logging, and traffic handling.

### Main Directory and Configurations
* **`config.go`**: Responsible for reading and loading rule JSON files into memory (`blocked.json` and `words.json`). It contains the `IsBlocked` logic to check the presence of a domain in the blocklist in a *case-insensitive* manner.
* **`logger.go`**: Records all requests in a structured history, storing the *Timestamp*, *URL*, and *Action* performed. Because the environment is concurrent (multiple Goroutines accessing the same resource), a `sync.Mutex` was implemented to guarantee mutual exclusion during write operations, avoiding *race conditions* and log file corruption.

### Request Handlers
* **`proxy_handler.go`**: The main proxy router. It intercepts requests, parses the destination URL, and evaluates what action to take based on the loaded configurations, directing the flow to the specific Handler.
* **`pass_handler.go` (Forwarding)**: Used when the requested site is unrestricted. It clones the original client request, clears incompatible encoding headers (such as `Accept-Encoding`) to allow subsequent analysis if necessary, dispatches the request to the destination server, and returns the full response to the client.
* **`block_handler.go` (Site Blocking)**: If the domain is listed in the blocklist file, the external request is aborted. The proxy responds directly to the client with a `403 Forbidden` status and renders a custom HTML template from `templates/blocked.html`.
* **`filter_handler.go` (Content Filter)**: Makes the request to the origin site and inspects whether the `Content-Type` is `text/html`. If positive, the page body is intercepted and processed by regular expressions, which locate forbidden words in a *case-insensitive* manner and replace them with the configured equivalent terms.
* **`connect_handler.go` (HTTPS Tunneling)**: Handles requests using the HTTP `CONNECT` method. It establishes a pure TCP connection (`net.Dial`) with the destination server, uses the `http.Hijacker` feature to take full control of the client's TCP socket, and starts two parallel Goroutines using `io.Copy` to transparently route encrypted data between endpoints.

## Prerequisites and Installation

1. **Install Go:** Make sure you have Go installed on your machine (version 1.25 or newer). Downloads available at [go.dev](https://go.dev/dl/).
2. **Clone the Repository:**
   ```bash
   git clone [https://github.com/Miguel-casarin/Web-Proxy-](https://github.com/Miguel-casarin/Web-Proxy-)
   cd Web-Proxy-
