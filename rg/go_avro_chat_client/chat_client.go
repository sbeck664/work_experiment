package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "gopkg.in/avro.v0"
)

type Client struct {
    socket net.Conn
    data   chan []byte
}

// This is a struct which supports the avro schema.
type Message struct {
    Author   string `avro:"author"` // the avro: tag specifies which avro field matches this field.
    Contents string `avro:"contents"`
}

func (client *Client) receive(parsedSchema avro.Schema, name string) {
    for {
        message := make([]byte, 4096)
        length, err := client.socket.Read(message)
        if err != nil {
            client.socket.Close()
            break
        }
        if length > 0 {
            decodedMessage := decode(parsedSchema, message)
            // Don't print message if the client is the author
            if (decodedMessage.Author != name) {
                fmt.Println("RECEIVED: Message from " + decodedMessage.Author + ": " + decodedMessage.Contents)
            }
        }
    }
}

func encode(parsedSchema avro.Schema, name string, contents string) (encodedMessage []byte) {

    // Create a SpecificDatumWriter, which you can re-use multiple times.
    writer := avro.NewSpecificDatumWriter()
    writer.SetSchema(parsedSchema)
    
    // Write a message to a byte buffer as avro.
    message := &Message{
        Author:   name,
        Contents: contents,
    }
    
    var buf bytes.Buffer
    encoder := avro.NewBinaryEncoder(&buf)
    err := writer.Write(message, encoder)
    if err != nil {
        log.Fatal(err)
    }
    
    return buf.Bytes()
}

func decode(parsedSchema avro.Schema, encodedAvro []byte) (decodedMessage Message) {
   
    // Create a SpecificDatumReader, which you can re-use multiple times.
    reader := avro.NewSpecificDatumReader()
    reader.SetSchema(parsedSchema)
    
    var message Message
    decoder := avro.NewBinaryDecoder(encodedAvro)
    err := reader.Read(&message, decoder)
    if err != nil {
        log.Fatal(err)
    }
    
    return message
}

func start(parsedSchema avro.Schema, port string, name string) {
    fmt.Println("Starting client " + name + " on port " + port)
    connection, error := net.Dial("tcp", "localhost:" + port)
    if error != nil {
        fmt.Println(error)
        return
    }
    client := &Client{socket: connection}
    go client.receive(parsedSchema, name)
    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
        encodedMessage := encode(parsedSchema, name, strings.TrimRight(message, "\n"))
        connection.Write(encodedMessage)
    }
}

func main() {
    
    schema := `{
        "type": "record",
        "name": "Message",
        "fields": [
            {"name": "author", "type": "string"},
            {"name": "contents", "type": "string"}
        ]
    }`

    // Parse a schema from JSON to get the schema object.
    parsedSchema, err := avro.ParseSchema(schema)
    if err != nil {
        log.Fatal(err)
        return
    }

    port := flag.String("port", "12345", "port to connect on")
    name := flag.String("name", "test", "name of client")
    flag.Parse()
    start(parsedSchema, *port, *name)
}
