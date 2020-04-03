package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "gopkg.in/avro.v0"
    "log"
    "math/rand"
    "net"
    "os"
    "strings"
    "time"
)

var (
    red     = "\033[1;31m%s\033[0m"
    green   = "\033[1;32m%s\033[0m"
    yellow  = "\033[1;33m%s\033[0m"
    purple  = "\033[1;34m%s\033[0m"
    magenta = "\033[1;35m%s\033[0m"
    teal    = "\033[1;36m%s\033[0m"
)

var (
    colours = []string{
        red,
        green,
        yellow,
        purple,
        teal,
    }
)

var (
    writer = avro.NewSpecificDatumWriter()
    reader = avro.NewSpecificDatumReader()
)

type Client struct {
    socket net.Conn
    data   chan []byte
}

// This is a struct which supports the avro schema.
type Message struct {
    Author   string `avro:"author"` // the avro: tag specifies which avro field matches this field.
    Colour   string `avro:"colour"`
    Contents string `avro:"contents"`
}

func (client *Client) receive() {

    for {
        message := make([]byte, 4096)
        length, err := client.socket.Read(message)
        if err != nil {
            client.socket.Close()
            break
        }
        if length > 0 {
            decodedMessage := decode(message)
            fmt.Printf(decodedMessage.Colour, decodedMessage.Author + ": " + decodedMessage.Contents + "\n")
        }
    }
}

func encode(name string, colour string, contents string) (encodedMessage []byte) {

    // Write a message to a byte buffer as avro.
    message := &Message{
        Author: name,
        Colour: colour,
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

func decode(encodedAvro []byte) (decodedMessage Message) {

    var message Message
    decoder := avro.NewBinaryDecoder(encodedAvro)
    error := reader.Read(&message, decoder)
    if error != nil {
        log.Fatal(error)
    }
    
    return message
}

func start(port string, name string) {

    fmt.Println("Starting chat as " + name + " on port " + port)
    connection, error := net.Dial("tcp", "localhost:" + port)
    if error != nil {
        fmt.Println(error)
        log.Fatal(error)
    }
    client := &Client{socket: connection}

    // Pick a random colour to use for the client's messages
    rand.Seed(time.Now().Unix())
    colour := colours[rand.Int() % len(colours)]

    go client.receive()
    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
        encodedMessage := encode(name, colour, strings.TrimRight(message, "\n"))
        // Delete message from stdin
        fmt.Printf("\033[A\033[2K")
        connection.Write(encodedMessage)
    }
}

func main() {

    schema := `{
        "type": "record",
        "name": "Message",
        "fields": [
            {"name": "author", "type": "string"},
            {"name": "colour", "type": "string"},
            {"name": "contents", "type": "string"}
        ]
    }`

    // Parse a schema from JSON to get the schema object.
    parsedSchema, error := avro.ParseSchema(schema)
    if error != nil {
        fmt.Println(error)
        log.Fatal(error)
    }

    writer.SetSchema(parsedSchema)
    reader.SetSchema(parsedSchema)

    port := flag.String("port", "12345", "port to connect on")
    name := flag.String("name", "test", "name of client")
    flag.Parse()
    start(*port, *name)
}
