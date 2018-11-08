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
#include <string.h>

#include "velen-run.h"
#include "pivot_root.h"
#include "config.h"

int bind_to_overlay(const char *path) {
  int len = strlen(VELEN_PATH) + strlen("/overlay") + strlen(path) + 1;
  char new_path[len];
  char *path_ptr = stpcpy(stpcpy(stpcpy(new_path, VELEN_PATH), "/overlay"), path);
  int err = mount(path, new_path, "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind mount");
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}

int switch_user(const char *username) {
  int err;
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

  err = setenv("USER", lord->pw_name, 1)
     || setenv("LOGNAME", lord->pw_name, 1)
     || setenv("HOME", lord->pw_dir, 1)
     || setenv("SHELL", "/bin/sh", 1)
     || setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin", 1);
  if (err) {
    perror("failed to set environment");
    return EXIT_FAILURE;
  }
}

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

  // TODO: is this necessary?
  err = mount(NULL, "/", NULL, MS_PRIVATE | MS_REC, NULL);
  if (err) {
    perror("failed to make root mount recursively private");
    return EXIT_FAILURE;
  }

  // TODO: also this
  err = mount(VELEN_PATH "/overlay", VELEN_PATH "/overlay", "", MS_BIND, NULL);
  if (err) {
    perror("failed to rebind velen");
    return EXIT_FAILURE;
  }

  err = bind_to_overlay("/dev")
    || bind_to_overlay("/dev/shm")
    || bind_to_overlay("/proc")
    || bind_to_overlay("/tmp");
  if (err != EXIT_SUCCESS) {
    return EXIT_FAILURE;
  }

  err = chroot(VELEN_PATH "/overlay");
  if (err) {
    perror("failed to chroot to overlay");
    return EXIT_FAILURE;
  }

  err = chdir("/");
  if (err) {
    perror("failed to chdir to root in overlay");
    return EXIT_FAILURE;
  }

  err = switch_user(VELEN_LORD);
  if (err) {
    return EXIT_FAILURE;
  }

  err = execvp(argv[0], argv);
  if (err) {
    perror("failed to execve to new process");
    return EXIT_FAILURE;
  }
  __builtin_unreachable();
}
