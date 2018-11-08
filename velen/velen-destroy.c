#include <sys/mount.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>

#include "velen-destroy.h"
#include "config.h"

int velen_destroy() {
  int err;

  err = umount2(VELEN_PATH "/overlay", MNT_DETACH);
  if (err && errno != EINVAL && errno != ENOENT) {
    perror("failed to unmount overlay");
    return EXIT_FAILURE;
  }

  err = umount2(VELEN_PATH "/top", MNT_DETACH);
  if (err && errno != EINVAL && errno != ENOENT) {
    perror("failed to unmount tmpfs");
    return EXIT_FAILURE;
  }

  return EXIT_SUCCESS;
}
