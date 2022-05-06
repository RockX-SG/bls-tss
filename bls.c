#include <stdio.h>
#include <stdlib.h>
#include "bls.h"

#define BUFFER_SIZE 2048

void proceed_if_needed(void *keygen){
    int res = keygen_wants_to_proceed(keygen);
    printf("keygen_wants_to_proceed: %d\n", res);
    if (res == 1) {
        res = keygen_proceed(keygen);
        printf("keygen_proceed: %d\n", res);
    }
}

void send_outgoing_if_there_is(void *keygen, char *buffer){
    int res = keygen_has_outgoing(keygen);
    printf("keygen_has_outgoing: %d\n", res);
    while (res > 0) {
        int outgoing_bytes_size = keygen_outgoing(keygen, buffer, BUFFER_SIZE);
        printf("outgoing bytes size: %d\n", outgoing_bytes_size);
        printf("outgoing is:\n");
        printf("\033[0;32m");
        printf("%s\n", buffer);
        printf("\033[0m");
        res = keygen_has_outgoing(keygen);
    }
}

void wait_for_incoming(void *keygen, char *buffer){
    printf("incoming > ");
    fgets(buffer, BUFFER_SIZE, stdin);
    keygen_incoming(keygen, buffer);
}

void finish_if_possible(void *keygen, char *buffer) {
    int finished = keygen_is_finished(keygen);
    if (finished != 1)
        return;
    int res = keygen_pick_output(keygen, buffer, BUFFER_SIZE);
    if (res > 0) {
        printf("Output is:\n%s\n", buffer);
    }
}

void interpret_loop(void *keygen, char *buffer){
    int finished = keygen_is_finished(keygen);
    while (finished != 1) {
        wait_for_incoming(keygen, buffer);
        send_outgoing_if_there_is(keygen, buffer);
        proceed_if_needed(keygen);
        send_outgoing_if_there_is(keygen, buffer);
        finished = keygen_is_finished(keygen);
    }
    finish_if_possible(keygen, buffer);
}

void other_check(void *keygen, char *buffer) {
    int n = keygen_parties(keygen);
    printf("parties 0: %d\n", n);
    n = keygen_parties(keygen);
    printf("parties 1: %d\n", n);
//    char msg[] = "{\"sender\":2,\"receiver\":null,\"body\":{\"Round1\":{\"com\":\"3083c12e285b7e388728f00270664e6fedd3c079fed0ac28595148082af629fc\"}}}";

    int res = keygen_current_round(keygen);
    printf("keygen_current_round: %d\n", res);
    res = keygen_total_rounds(keygen);
    printf("keygen_total_rounds: %d\n", res);
    res = keygen_party_ind(keygen);
    printf("keygen_party_ind: %d\n", res);
    res = keygen_parties(keygen);
    printf("keygen_parties: %d\n", res);
    proceed_if_needed(keygen);
}

void init(void *keygen, char *buffer) {
    proceed_if_needed(keygen);
    send_outgoing_if_there_is(keygen, buffer);
    finish_if_possible(keygen, buffer);
}

int main(int argc, char *argv[]) {
    if(argc < 4) {
        printf("Error: Insufficient arguments.\n");
        printf("Usage:\n");
        printf("  bls <i> <t> <n>\n");
        return -1;
    }
    int i = atoi(argv[1]);
    int t = atoi(argv[2]);
    int n = atoi(argv[3]);
    char *buffer = (char*)malloc(BUFFER_SIZE*sizeof(unsigned char));
    void* keygen = new_keygen(i, t, n);
    init(keygen, buffer);
    interpret_loop(keygen, buffer);
    free_keygen(keygen);
    free(buffer);
    return 0;
}

// gcc -Wall bls.c -o bls -L/home/gldeng/local/gldeng/multi-party-bls-wrapper/target/release -lmulti_party_bls_wrapper
// export LD_LIBRARY_PATH=/home/gldeng/local/gldeng/multi-party-bls-wrapper/target/release
