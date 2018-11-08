#include <sys/stat.h>
#include <sys/mount.h>
#include <errno.h>
#include <pwd.h>
#include <unistd.h>
#include "config.h"

/**
 * @brief Prepares Velen for use.
 *
 * Creates overlayfs based on / in @p velen_path.
 * Velen must be destroyed first.
 *
 * @param paths – array of paths to be chown'd to @p velen_lord. Last item must be NULL.
 */
inline static int velen_prepare(char* paths[]) {
  int err;

  err = mkdir(velen_path, 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen");
    return err;
  }

  err = mkdir(velen_path "/top", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/top");
    return err;
  }

  err = mount("velentmp", velen_path "/top", "tmpfs", 0, NULL);
  if (err) {
    perror("failed to mount a tmpfs in velen/top");
    return err;
  }

  err = mkdir(velen_path "/top" "/workdir", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/top/workdir");
    return err;
  }

  err = mkdir(velen_path "/top" "/layer", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/top/layer");
    return err;
  }

  err = mkdir(velen_path "/overlay", 0755);
  if (err && errno != EEXIST) {
    perror("failed to mkdir velen/overlay");
    return err;
  }

  err = mount("velen", velen_path "/overlay", "overlay", 0,
      "lowerdir=/" ","
      "upperdir=" velen_path "/top/layer" ","
      "workdir=" velen_path "/top/workdir");
  if (err) {
    perror("failed to mount overlayfs");
    return err;
  }

  if (paths[0] != NULL) {
    struct passwd *lord = getpwnam(velen_lord);
    if (lord == NULL) {
      perror("failed to get info on Velen lord user");
      return EXIT_FAILURE;
    }

    err = chroot(velen_path "/overlay");
    if (err) {
      perror("failed to chroot to Velen for chowning");
      return err;
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
