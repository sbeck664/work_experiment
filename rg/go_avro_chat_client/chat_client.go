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
    "os/signal"
    "strings"
    "syscall"
    "time"
)

const (
    black   = "\033[1;30m%s\033[0m"
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

const (
    chatType = "chat"
    joinedChatType = "joined"
    leftChatType = "left"
    presenceRequestType = "presenceRequest"
    presenceResponseType = "presenceResponse"
)

var (
    writer = avro.NewSpecificDatumWriter()
    reader = avro.NewSpecificDatumReader()
)

type Client struct {
    socket net.Conn
    name   string
    colour string
    data   chan []byte
}

// This is a struct which supports the avro schema.
type Message struct {
    Author   string `avro:"author"` // the avro: tag specifies which avro field matches this field.
    Colour   string `avro:"colour"`
    Type     string `avro:"type"`
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
            switch decodedMessage.Type {
                case chatType:
                    fmt.Printf(decodedMessage.Colour, decodedMessage.Author + ": " + decodedMessage.Contents + "\n")
                case joinedChatType:
                    if (decodedMessage.Author != client.name) {
                        fmt.Printf(black, decodedMessage.Author + " has joined the chat\n")
                    } else {
                        client.send_notification(presenceRequestType)
                    }
                case leftChatType:
                    if (decodedMessage.Author != client.name) {
                        fmt.Printf(black, decodedMessage.Author + " has left the chat\n")
                    }
                case presenceRequestType:
                    if (decodedMessage.Author != client.name) {
                        client.send_response(decodedMessage.Author, presenceResponseType, client.name + " is also a participant in this chat")
                    }
                case presenceResponseType:
                    if (decodedMessage.Author == client.name) {
                        fmt.Printf(black, decodedMessage.Contents + "\n")
                    }
            }
        }
    }
}

func (client *Client) send_response(name string, messageType string, contents string) {
    message := encode(name, client.colour, messageType, contents)
    client.socket.Write(message)
}

func (client *Client) send_chat_message(contents string) {
    message := encode(client.name, client.colour, chatType, contents)
    client.socket.Write(message)
}

func (client *Client) send_notification(messageType string) {
    message := encode(client.name, client.colour, messageType, "")
    client.socket.Write(message)
}

func encode(name string, colour string, messageType string, contents string) (encodedMessage []byte) {

    // Write a message to a byte buffer as avro.
    message := &Message{
        Author:   name,
        Colour:   colour,
        Type:     messageType,
        Contents: contents,
    }

    var buf bytes.Buffer
    encoder := avro.NewBinaryEncoder(&buf)
    error := writer.Write(message, encoder)
    if error != nil {
        return []byte{}
    }
    
    return buf.Bytes()
}

func decode(encodedAvro []byte) (decodedMessage Message) {

    var message Message
    decoder := avro.NewBinaryDecoder(encodedAvro)
    error := reader.Read(&message, decoder)
    if error != nil {
        message.Type = "unknown"
    }
    
    return message
}

func start(port string, name string) {

    fmt.Printf(black, "Hello " + name + ", welcome to the chat\n")
    connection, error := net.Dial("tcp", "localhost:" + port)
    if error != nil {
        fmt.Println(error)
        log.Fatal(error)
    }

    // Pick a random colour to use for the client's messages
    rand.Seed(time.Now().Unix())
    colour := colours[rand.Int() % len(colours)]

    client := &Client{socket: connection, name: name, colour:colour}

    client.send_notification(joinedChatType)

    // Ctrl-C signal handler
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        client.send_notification(leftChatType)
        fmt.Printf(black, "\rYou have left the chat\n")
        os.Exit(0)
    }()

    go client.receive()

    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
        // Delete message from stdin
        fmt.Printf("\033[A\033[2K")
        client.send_chat_message(strings.TrimRight(message, "\n"))
    }
}

func main() {

    port := flag.String("port", "12345", "port to connect on")
    name := flag.String("name", "test", "name of client")
    flag.Parse()

    schema := `{
        "type": "record",
        "name": "Message",
        "fields": [
            {"name": "author", "type": "string"},
            {"name": "colour", "type": "string"},
            {"name": "type", "type": "string"},
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

    start(*port, *name)
}
