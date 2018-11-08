#define _GNU_SOURCE
#include <sched.h>
#include <linux/sched.h>
#include <errno.h>

#include "pivot_root.c"
#include "config.h"

inline static int velen_run(char* argv[], char* envp[]) {
  int err;

  err = unshare(CLONE_NEWNS);
  if (err) {
    perror("failed to unshare");
    return err;
  }

  err = mount("none", "/", NULL, MS_PRIVATE | MS_REC, NULL);
  if (err) {
    perror("failed to make root mount recursively private");
    return err;
  }

  err = mount(velen_path "/overlay", velen_path "/overlay", "", MS_BIND, NULL);
  if (err) {
    perror("failed to rebind velen");
    return err;
  }

  err = chdir(velen_path "/overlay");
  if (err) {
    perror("failed to chdir to overlay");
    return err;
  }

  err = pivot_root(".", ".");
  if (err) {
    perror("failed to pivot_root to overlay");
    return err;
  }

  err = chroot(".");
  if (err) {
    perror("failed to chroot to overlay");
    return err;
  }

  err = umount2("/", MNT_DETACH);
  if (err) {
    perror("failed to detach old root");
    return err;
  }

  struct passwd *lord = getpwnam(velen_lord);
  if (lord == NULL) {
    perror("failed to get info on Velen lord user");
    return EXIT_FAILURE;
  }

  err = setgid(lord->pw_gid);
  if (err) {
    perror("failed to set group ID");
    return EXIT_FAILURE;
  }

  err = setuid(lord->pw_uid);
  if (err) {
    perror("failed to set user ID");
    return EXIT_FAILURE;
  }

  err = execve(argv[0], argv, envp);
  if (err) {
    perror("failed to execve to new process");
    return err;
  }
  __builtin_unreachable();
}
