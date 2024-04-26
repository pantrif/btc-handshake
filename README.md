
## Overview of the Solution
The solution involves two main components:

1. Bitcoin Node

Running in a Docker container configured for regtest mode, which is used for local development and testing. 
This mode allows us to control block generation and thus enables rapid testing without real-world blockchain 
constraints like actual transaction confirmation times.

2. Go Client Application 

Also running in a Docker container, designed to connect to the Bitcoin node and perform a network handshake using 
the Bitcoin protocol. This includes sending a version message and receiving a verack message.

Key features of the Go client:
- **Network Communication:** Uses standard Go libraries (net, encoding/binary) to manage TCP connections and binary 
data manipulation.
- **Concurrency and Error Handling:** Employs goroutines to handle asynchronous I/O operations efficiently, 
with timeouts to avoid hanging indefinitely.
- **Logging:** Extensive logging to facilitate debugging and ensure that each step of the handshake is traceable.


## Building and Running
1. Set Up the Bitcoin Node Configuration: Create bitcoin.conf in the bitcoin-data directory with appropriate settings for regtest.
2. Build and Start with Docker Compose:
```shell
docker-compose build
docker-compose run client
```
Expected output:
```shell
 ✔ Container bitcoin-node-handshake-node-1  Running                                                                                                                                                                                                                                                                0.0s 
2024/04/26 13:51:04 Attempting to connect to the Bitcoin node at node:18444...
2024/04/26 13:51:04 Connected successfully to the Bitcoin node.
2024/04/26 13:51:04 Sending the version message...
2024/04/26 13:51:04 Message sent successfully.
2024/04/26 13:51:04 received version message, sending verack.
2024/04/26 13:51:04 sending the verack message...
2024/04/26 13:51:04 Message sent successfully.
2024/04/26 13:51:04 received verack message, handshake completed successfully.
```

## Future improvements

- Validation of Incoming Messages: Adding checks to validate the contents of the incoming version message from the peer could prevent potential security issues or miscommunications.
- Dynamic Port Configuration: Adaptability to different network setups, such as testnet or regtest, by dynamically adjusting network magic values and port numbers.
- Use of Secure and Reliable Libraries: While the current implementation avoids using external libraries for the handshake, ensuring that all cryptographic and network operations are secure and efficient is crucial.