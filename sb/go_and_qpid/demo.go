package main

// #cgo LDFLAGS: -L . -lqpid -lqpidmessaging -lstdc++
// #include <stdlib.h>
// #include "qpid.h"
import "C"

import "fmt"
import "unsafe"

func main() {

    address := C.CString("localhost:5672")
    conn := C.new_qpid_connection(address)
    C.free(unsafe.Pointer(address))

    rec_queue := C.CString("receiver_queue")
    C.add_receiver(conn, rec_queue)
    C.free(unsafe.Pointer(address))

    mess := C.receive_message(conn)
    mess_bytes := C.GoBytes(mess.data, mess.length)

    fmt.Println("Received ", mess_bytes)
    C.delete_qpid_connection(conn)
}
