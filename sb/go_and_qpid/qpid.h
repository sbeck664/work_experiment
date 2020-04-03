
#ifdef __cplusplus
extern "C" {
#endif

typedef void * QpidConnection;

struct Message
{
    const void *data;
    int length;
};

QpidConnection new_qpid_connection(char *address);

void add_receiver(QpidConnection conn, char *queue_name);

void add_sender(QpidConnection conn, char *sender_name, char *queue_name);

void send_message(QpidConnection conn, char *sender_name, void *message_data, int message_size);

struct Message receive_message(QpidConnection conn);

void delete_qpid_connection(QpidConnection conn);

#ifdef __cplusplus
}
#endif
