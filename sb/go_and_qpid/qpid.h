
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

struct Message receive_message(QpidConnection conn);

void delete_qpid_connection(QpidConnection conn);

#ifdef __cplusplus
}
#endif
