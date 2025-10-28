# Chit-Chat-distributed-systems
Mandatory activity 3 in course distributed systems

## Running the program
You need to install [Go](https://go.dev/doc/install) to run the program.
To run the program you must first start the server. From the root directory, run:
```bash
go run server/server.go
```
The server will start, and output logs are provided in a file `events.log`.

Then you can start up the clients. Be sure to include a username as an argument.
```bash
go run client/client.go <username>
```
If everything is working correctly, the client will see the following:
```
[1] Server: Participant <username> joined Chit Chat at logical time 1
> 
```
You can now write messages in the chat room.