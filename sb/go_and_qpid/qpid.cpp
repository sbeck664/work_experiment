
#include "qpid.h"
#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Duration.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>

#include <iostream>
#include <map>

struct _QpidConnection
{
    qpid::messaging::Connection conn;
    qpid::messaging::Session session;
    qpid::messaging::Receiver receiver;
    bool connected;
    bool have_receiver;
    std::map<std::string, qpid::messaging::Sender> senders;
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
        conn_obj->receiver = conn_obj->session.createReceiver(queue_name);
        conn_obj->have_receiver = true;
        std::cout << "add_receiver" << std::endl;
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
    }
}

void add_sender(QpidConnection conn, char *sender_name, char *queue_name)
{
    _QpidConnection *conn_obj = (_QpidConnection *)conn;
    if (!conn_obj->connected)
    {
        std::cout << "not connected" << std::endl;
        return;
    }
    try
    {
        conn_obj->senders[sender_name] = conn_obj->session.createSender(queue_name);
        std::cout << "add_sender" << std::endl;
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
    }
}

void send_message(QpidConnection conn, char *sender_name, void *message_data, int message_size)
{
    _QpidConnection *conn_obj = (_QpidConnection *)conn;
    if (!conn_obj->connected)
    {
        std::cout << "not connected" << std::endl;
        return;
    }

    std::map<std::string, qpid::messaging::Sender>::iterator it;
    it = conn_obj->senders.find(sender_name);
    if (it == conn_obj->senders.end())
    {
        std::cout << "no such sender" << std::endl;
        return;
    }
    qpid::messaging::Message message;
    message.setContent((char *)message_data, message_size);
    it->second.send(message);
    std::cout << "send_message" << std::endl;
}

struct Message receive_message(QpidConnection conn)
{
    static qpid::messaging::Message message;

    _QpidConnection *conn_obj = (_QpidConnection *)conn;
    if (conn_obj->connected && conn_obj->have_receiver)
    {
        conn_obj->receiver.fetch(message, qpid::messaging::Duration::FOREVER);
        conn_obj->session.acknowledge(true);
        std::cout << "receive_message" << std::endl;
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
