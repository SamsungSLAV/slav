#include <sys/mount.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>

#include "velen-destroy.h"
#include "config.h"

int detach_mount(const char *path) {
  int err = umount2(path, MNT_DETACH);
  if (err && errno != EINVAL && errno != ENOENT) {
    perror("failed to unmount");
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}

int velen_destroy() {
  return detach_mount(VELEN_PATH "/overlay")
    || detach_mount(VELEN_PATH "/top")
    || EXIT_SUCCESS;
}
