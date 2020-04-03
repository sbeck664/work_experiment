package qpid

// #cgo LDFLAGS: -L .. -lqpid -lqpidmessaging -lstdc++
// #include <stdlib.h>
// #include "../qpid.h"
import "C"

import "unsafe"

type Connection struct {
    c C.QpidConnection
}

func NewQpidConnection(address string) Connection {

    var c Connection

    c_address := C.CString(address)
    c.c = C.new_qpid_connection(c_address)
    C.free(unsafe.Pointer(c_address))
    return c
}

func AddReceiver(conn Connection, queue_name string) {
    c_queue_name := C.CString(queue_name)
    C.add_receiver(conn.c, c_queue_name)
    C.free(unsafe.Pointer(c_queue_name))
}

func AddSender(conn Connection, sender_name string, queue_name string) {
    c_queue_name := C.CString(queue_name)
    c_sender_name := C.CString(sender_name)
    C.add_sender(conn.c, c_sender_name, c_queue_name)
    C.free(unsafe.Pointer(c_queue_name))
    C.free(unsafe.Pointer(c_sender_name))
}

func ReceiveMessage(conn Connection) []byte {
    mess := C.receive_message(conn.c)
    mess_bytes := C.GoBytes(mess.data, mess.length)
    return mess_bytes
}

func SendMessage(conn Connection, sender_name string, message []byte) {
    c_message := C.CBytes(message)
    c_sender_name := C.CString(sender_name)
    C.send_message(conn.c, c_sender_name, c_message, C.int(len(message)))
    C.free(unsafe.Pointer(c_message))
    C.free(unsafe.Pointer(c_sender_name))
}

func DeleteQpidConnection(conn Connection) {
    C.delete_qpid_connection(conn.c)
}
