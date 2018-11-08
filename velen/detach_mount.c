#include <sys/mount.h>
#include <stdio.h>
#include <errno.h>
#include <stdlib.h>
#include "detach_mount.h"

int detach_mount(const char *path) {
  int err = umount2(path, MNT_DETACH);
  if (err && errno != EINVAL && errno != ENOENT) {
    perror("failed to unmount");
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}

