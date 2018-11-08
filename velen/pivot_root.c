#define _GNU_SOURCE
#include <unistd.h>
#include <sys/syscall.h>

inline static int pivot_root(const char* new_root, const char* put_old) {
  return syscall(SYS_pivot_root, new_root, put_old);
}
