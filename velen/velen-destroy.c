#include <sys/mount.h>
#include <errno.h>

#include "config.h"

inline static int velen_destroy() {
  int err;

  err = umount2(velen_path "/overlay", MNT_DETACH);
  if (err && errno != EINVAL && errno != ENOENT) {
    perror("failed to unmount overlay");
    return err;
  }

  err = umount2(velen_path "/top", MNT_DETACH);
  if (err && errno != EINVAL && errno != ENOENT) {
    perror("failed to unmount tmpfs");
    return err;
  }

  return EXIT_SUCCESS;
}
