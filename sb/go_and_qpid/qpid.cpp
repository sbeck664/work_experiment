
#include "qpid.h"
#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Duration.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Session.h>

#include <iostream>

struct _QpidConnection
{
    qpid::messaging::Connection conn;
    qpid::messaging::Session session;
    qpid::messaging::Receiver receiver;
    bool connected;
    bool have_receiver;
};

QpidConnection new_qpid_connection(char *address)
{
    std::cout << "new_qpid_connection" << std::endl;
    _QpidConnection *conn = new _QpidConnection;

    conn->conn = qpid::messaging::Connection(address, "");
    conn->connected = false;
    conn->have_receiver = false;
    try
    {
        conn->conn.open();
        conn->connected = true;
        conn->session = conn->conn.createSession();
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
    }
    return (void *)conn;
}

void add_receiver(QpidConnection conn, char *queue_name)
{
    _QpidConnection *conn_obj = (_QpidConnection *)conn;
    if (!conn_obj->connected)
    {
        std::cout << "not connected" << std::endl;
        return;
    }
    try
    {
        std::cout << "add_receiver" << std::endl;
        conn_obj->receiver = conn_obj->session.createReceiver(queue_name);
        conn_obj->have_receiver = true;
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
    }
}

struct Message receive_message(QpidConnection conn)
{
    static qpid::messaging::Message message;

    _QpidConnection *conn_obj = (_QpidConnection *)conn;
    if (conn_obj->connected && conn_obj->have_receiver)
    {
        conn_obj->receiver.fetch(message, qpid::messaging::Duration::FOREVER);
        conn_obj->session.acknowledge(true);
        std::cout << "receive_message and ack" << std::endl;
    }
    else
    {
        std::cout << "receiver not connected" << std::endl;
    }
    struct Message ret;
    ret.data = message.getContentPtr();
    ret.length = message.getContentSize();
    return ret;
}

void delete_qpid_connection(QpidConnection conn)
{
    std::cout << "delete_qpid_connection" << std::endl;
    _QpidConnection *conn_obj = (_QpidConnection *)conn;
    delete conn_obj;
}
