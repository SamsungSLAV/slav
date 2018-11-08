#define _XOPEN_SOURCE 500
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/mount.h>
#include <errno.h>
#include <pwd.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>

#include "velen-prepare.h"
#include "config.h"

/**
 * @brief Prepares Velen for use.
 *
 * Creates overlayfs based on / in @p VELEN_PATH.
 * Velen must be destroyed first if it exists.
 *
 * @param paths – null-terminated array of paths to be chown'd to @p VELEN_LORD.
 */
int velen_prepare(char* paths[]) {
  int err;

  err = mkdir(VELEN_PATH, 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen");
    return EXIT_FAILURE;
  }

  err = mkdir(VELEN_PATH "/top", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/top");
    return EXIT_FAILURE;
  }

  err = mount("velentmp", VELEN_PATH "/top", "tmpfs", 0, NULL);
  if (err) {
    perror("failed to mount a tmpfs in velen/top");
    return EXIT_FAILURE;
  }

  err = mkdir(VELEN_PATH "/top" "/workdir", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/top/workdir");
    return EXIT_FAILURE;
  }

  err = mkdir(VELEN_PATH "/top" "/layer", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/top/layer");
    return EXIT_FAILURE;
  }

  err = mkdir(VELEN_PATH "/overlay", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/overlay");
    return EXIT_FAILURE;
  }

  err = mount("velen", VELEN_PATH "/overlay", "overlay", 0,
      "lowerdir=/" ","
      "upperdir=" VELEN_PATH "/top/layer" ","
      "workdir=" VELEN_PATH "/top/workdir");
  if (err) {
    perror("failed to mount overlayfs");
    return EXIT_FAILURE;
  }

  if (paths[0] != NULL) {
    struct passwd *lord = getpwnam(VELEN_LORD);
    if (lord == NULL) {
      perror("failed to get info on Velen lord user");
      return EXIT_FAILURE;
    }

    err = chroot(VELEN_PATH "/overlay");
    if (err) {
      perror("failed to chroot to Velen for chowning");
      return EXIT_FAILURE;
    }

    for (int i = 0; paths[i] != NULL; i++) {
      err = chown(paths[i], lord->pw_uid, lord->pw_gid);
      if (err) {
        fputs(paths[i], stderr);
        perror(": failed to chown");
      }
    }
  }

  return EXIT_SUCCESS;
}
