#define _GNU_SOURCE
#include <unistd.h>
#include <sched.h>
#include <linux/sched.h>
#include <errno.h>
#include <sys/types.h>
#include <grp.h>

#include "pivot_root.c"
#include "config.h"

extern char **environ;

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

  err = mount("/dev", velen_path "/overlay" "/dev", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind dev");
    return err;
  }

  err = mount("/dev/shm", velen_path "/overlay" "/dev/shm", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind devshm");
    return err;
  }

  err = mount("/proc", velen_path "/overlay" "/proc", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind proc");
    return err;
  }

  err = mount("/tmp", velen_path "/overlay" "/tmp", "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind tmp");
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

  err = execve(argv[0], argv, environ);
  if (err) {
    perror("failed to execve to new process");
    return err;
  }
  __builtin_unreachable();
}
