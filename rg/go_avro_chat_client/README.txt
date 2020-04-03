Test using Go with avro

Get the Go avro package:

  go get gopkg.in/avro.v0

Firstly follow the instructions in sb/helm_stuff/README.txt to get a chat server service set up (only follow the server related steps)

Use 'kubectl get services' to find out the external TCP port for the server service then run the following to start up a client:

  go run chat_client.go --port <external TCP port> --name <unique name>

This will start up a client which encodes messages as avro binary before sending them and decodes received messages before printing them. It also sends messages with the client's name and a random colour for the messages.

You can then start up many more clients with different names using the same port and chat between them.
