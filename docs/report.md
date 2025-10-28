# Mandatory Activity 3 - Report

## Our streaming 
We are using server-side streaming. This is because it is not important if a client recieves an event or not. The server already has to send out a lot of messages to every client on every new event, and if the server then also has to recieve confirmations from all the clients, then that would be a lot of unnecessary packages.

## System architecture
We use a client-server architecture where the clients send messages to the server, and the server relays messages to other clients. Considering the service is structured as client/server already, it is the obvious choice.

## RPC methods
We have three (3) RPC methods defined:

1. `rpc Join (User) returns (stream Chat)`: This is the method used by clients to open a stream where messages can be recieved. The username specified in the User will also be broadcasted to any already connected clients.

2. `rpc SendChat (Chat) returns (Empty) {};` Clients can send a chat to the server using this method. The chat contains timestamp, username and message. Nothing is returned; instead the message is recieved from the stream that was opened by `Join`.

3. `rpc Leave (User) returns (Empty) {};` This method closes the stream associated with the username specified in the User. A leave message will also be broadcasted.

## Calculation of timestamps
We implemented the Lamport timestamp algorithm. Before every messages sent by the client the local timestamp gets incremented and sent with the message, the server then updates its local timestamp based on the max of the recived timestamp and its own. The client does the same when receiving a timestamp. Server increments local timestamp when user joining or leaving.

## The diagram of RPC calls
![](RPC%20calls.drawio.png)

## Appendices

### Link
[https://github.com/LucasITU/Chit-Chat-distributed-systems]()

### Logs
```
2025/10/28 17:26:26 [Timestamp:0] [Component:Server] [EventType:Start] [Identifier:] 
2025/10/28 17:26:29 [Timestamp:1] [Component:Client] [EventType:Join] [Identifier:client-1] 
2025/10/28 17:26:29 [Timestamp:1] [Component:Client] [EventType:Recieve] [Identifier:Server] message:"Participant client-1 joined Chit Chat at logical time 1"  author:"Server"  timestamp:1
2025/10/28 17:26:29 [Timestamp:2] [Component:Server] [EventType:Delivery] [Identifier:client-1] 
2025/10/28 17:26:38 [Timestamp:3] [Component:Client] [EventType:Join] [Identifier:client-2] 
2025/10/28 17:26:38 [Timestamp:3] [Component:Client] [EventType:Recieve] [Identifier:Server] message:"Participant client-2 joined Chit Chat at logical time 3"  author:"Server"  timestamp:3
2025/10/28 17:26:38 [Timestamp:4] [Component:Server] [EventType:Delivery] [Identifier:client-2] 
2025/10/28 17:26:38 [Timestamp:4] [Component:Server] [EventType:Delivery] [Identifier:client-1] 
2025/10/28 17:26:46 [Timestamp:7] [Component:Client] [EventType:Recieve] [Identifier:client-1] message:"hello\r"  author:"client-1"  timestamp:6
2025/10/28 17:26:46 [Timestamp:8] [Component:Server] [EventType:Delivery] [Identifier:client-1] 
2025/10/28 17:26:46 [Timestamp:8] [Component:Server] [EventType:Delivery] [Identifier:client-2] 
2025/10/28 17:28:06 [Timestamp:9] [Component:Server] [EventType:Stop] [Identifier:] 
```