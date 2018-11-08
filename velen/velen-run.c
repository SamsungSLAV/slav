#define _GNU_SOURCE
#include <unistd.h>
#include <sched.h>
#include <linux/sched.h>
// fun fact: neither sched.h nor linux/sched.h declare unshare on my glibc
#include <errno.h>
#include <stdio.h>
#include <sys/mount.h>
#include <sys/types.h>
#include <pwd.h>
#include <grp.h>
#include <stdlib.h>

#include "velen-run.h"
#include "pivot_root.h"
#include "config.h"

/**
 * @brief Runs a program in Velen.
 *
 * @param argv – null-terminated array of argv. TODO
 */
int velen_run(char* argv[]) {
  int err;

  err = unshare(CLONE_NEWNS);
  if (err) {
    perror("failed to unshare");
    return EXIT_FAILURE;
  }

  err = mount(NULL, "/", NULL, MS_PRIVATE | MS_REC, NULL);
  if (err) {
    perror("failed to make root mount recursively private");
    return EXIT_FAILURE;
  }

  err = mount(VELEN_PATH "/overlay", VELEN_PATH "/overlay", "", MS_BIND, NULL);
  if (err) {
    perror("failed to rebind velen");
    return EXIT_FAILURE;
  }

  err = mount("/dev", VELEN_PATH "/overlay" "/dev", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind dev");
    return EXIT_FAILURE;
  }

  err = mount("/dev/shm", VELEN_PATH "/overlay" "/dev/shm", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind devshm");
    return EXIT_FAILURE;
  }

  err = mount("/proc", VELEN_PATH "/overlay" "/proc", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind proc");
    return EXIT_FAILURE;
  }

  err = mount("/tmp", VELEN_PATH "/overlay" "/tmp", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind tmp");
    return EXIT_FAILURE;
  }

  err = chdir(VELEN_PATH "/overlay");
  if (err) {
    perror("failed to chdir to overlay");
    return EXIT_FAILURE;
  }

  err = pivot_root(".", ".");
  if (err) {
    perror("failed to pivot_root to overlay");
    return EXIT_FAILURE;
  }

  err = chroot(".");
  if (err) {
    perror("failed to chroot to overlay");
    return EXIT_FAILURE;
  }

  err = umount2("/", MNT_DETACH);
  if (err) {
    perror("failed to detach old root");
    return EXIT_FAILURE;
  }

  struct passwd *lord = getpwnam(VELEN_LORD);
  if (lord == NULL) {
    perror("failed to get info on Velen lord user");
    return EXIT_FAILURE;
  }

  err = initgroups(lord->pw_name, lord->pw_gid);
  if (err) {
    perror("failed to change groups");
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

  err = clearenv();
  if (err) {
    perror("failed to clear env");
    return EXIT_FAILURE;
  }

  err = setenv("USER", lord->pw_name, 1);
  if (err) {
    perror("failed to set $USER");
    return EXIT_FAILURE;
  }

  err = setenv("LOGNAME", lord->pw_name, 1);
  if (err) {
    perror("failed to set $LOGNAME");
    return EXIT_FAILURE;
  }

  err = setenv("HOME", lord->pw_dir, 1);
  if (err) {
    perror("failed to set $HOME");
    return EXIT_FAILURE;
  }

  err = setenv("SHELL", "/bin/sh", 1);
  if (err) {
    perror("failed to set $SHELL");
    return EXIT_FAILURE;
  }

  err = setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin", 1);
  if (err) {
    perror("failed to set $USER");
    return EXIT_FAILURE;
  }

  err = execvp(argv[0], argv);
  if (err) {
    perror("failed to execve to new process");
    return EXIT_FAILURE;
  }
  __builtin_unreachable();
}
