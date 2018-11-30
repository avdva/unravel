# unravel
An HTTP server that accepts form input events as POST requests (JSON).

## Client
Is a pure-JS page, that contains a form and sends data to `localhost:5000/api/card`

## Server
Is a golang app, that handles requests from clients.

### Command line arguments:
- addr: an address to serve on.
- hash: hash algorithm. Currently only `pjw` is supported.
