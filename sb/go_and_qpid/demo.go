package main

import "fmt"
import "../go_and_qpid/qpid"

func main() {

    conn := qpid.NewQpidConnection("localhost:5672")

    qpid.AddReceiver(conn, "receiver_queue")

    qpid.AddSender(conn, "sender1", "receiver_queue")

    out_mess := []byte("a message")
    qpid.SendMessage(conn, "sender1", out_mess)

    mess := qpid.ReceiveMessage(conn)

    fmt.Println("Received ", string(mess))
    qpid.DeleteQpidConnection(conn)
}
