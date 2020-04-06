
#include "qpid.h"
#include <qpid/messaging/Connection.h>
#include <qpid/messaging/Duration.h>
#include <qpid/messaging/Message.h>
#include <qpid/messaging/Receiver.h>
#include <qpid/messaging/Sender.h>
#include <qpid/messaging/Session.h>

#include <iostream>
#include <map>
#include <optional>

struct _QpidConnection
{
    std::optional<qpid::messaging::Connection> conn {};
    std::optional<qpid::messaging::Session> session {};
    std::optional<qpid::messaging::Receiver> receiver {};
    std::map<std::string, qpid::messaging::Sender> senders {};
};

QpidConnection new_qpid_connection(char *address)
{
    std::cout << "new_qpid_connection" << std::endl;
    auto conn = new _QpidConnection;

    conn->conn = qpid::messaging::Connection(address, "");
    try
    {
        conn->conn->open();
        conn->session = conn->conn->createSession();
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
        conn->conn = {};
        conn->session = {};
    }
    return reinterpret_cast<void *>(conn);
}

void add_receiver(QpidConnection conn, char *queue_name)
{
    auto conn_obj = reinterpret_cast<_QpidConnection *>(conn);
    if (!(conn_obj->conn && conn_obj->session))
    {
        std::cout << "not connected" << std::endl;
        return;
    }
    try
    {
        conn_obj->receiver = conn_obj->session->createReceiver(queue_name);
        std::cout << "add_receiver" << std::endl;
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
    }
}

void add_sender(QpidConnection conn, char *sender_name, char *queue_name)
{
    auto conn_obj = reinterpret_cast<_QpidConnection *>(conn);
    if (!(conn_obj->conn && conn_obj->session))
    {
        std::cout << "not connected" << std::endl;
        return;
    }
    try
    {
        conn_obj->senders[sender_name] = conn_obj->session->createSender(queue_name);
        std::cout << "add_sender" << std::endl;
    }
    catch(const std::exception &e)
    {
        std::cout << e.what() << std::endl;
    }
}

void send_message(QpidConnection conn, char *sender_name, void *message_data, int message_size)
{
    auto conn_obj = reinterpret_cast<_QpidConnection *>(conn);
    if (!(conn_obj->conn && conn_obj->session))
    {
        std::cout << "not connected" << std::endl;
        return;
    }

    auto it = conn_obj->senders.find(sender_name);
    if (it == conn_obj->senders.end())
    {
        std::cout << "no such sender" << std::endl;
        return;
    }
    qpid::messaging::Message message;
    message.setContent(reinterpret_cast<char *>(message_data), message_size);
    it->second.send(message);
    std::cout << "send_message" << std::endl;
}

struct Message receive_message(QpidConnection conn)
{
    static qpid::messaging::Message message;

    auto conn_obj = reinterpret_cast<_QpidConnection *>(conn);
    if (!(conn_obj->conn && conn_obj->session && conn_obj->receiver))
    {
        conn_obj->receiver->fetch(message, qpid::messaging::Duration::FOREVER);
        conn_obj->session->acknowledge(true);
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
    auto conn_obj = reinterpret_cast<_QpidConnection *>(conn);
    delete conn_obj;
}
