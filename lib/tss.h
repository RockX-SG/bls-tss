#ifndef _BLS_H_
#define _BLS_H_

void* new_keygen(int i, int t, int n);
void free_keygen(void* state);
int keygen_current_round(const void* state);
int keygen_total_rounds(const void* state);
int keygen_party_ind(const void* state);
int keygen_parties(const void* state);
int keygen_wants_to_proceed(const void* state);
int keygen_proceed(void* state);
int keygen_has_outgoing(void* state);
int keygen_is_finished(void* state);
int keygen_pick_output(void* state, char* buf, int max_len);
int keygen_incoming(void* state, const char* msg);
int keygen_outgoing(void* state, char* buf, int max_len);

void* new_sign(const char* msg_hash, int i, int n, const char* local_key);
void free_sign(void* state);
int sign_current_round(const void* state);
int sign_total_rounds(const void* state);
int sign_party_ind(const void* state);
int sign_parties(const void* state);
int sign_wants_to_proceed(const void* state);
int sign_proceed(void* state);
int sign_has_outgoing(void* state);
int sign_is_finished(void* state);
int sign_pick_output(void* state, char* buf, int maxlen);
int sign_incoming(void* state, const char* msg);
int sign_outgoing(void* state, char* buf, int maxlen);

#endif