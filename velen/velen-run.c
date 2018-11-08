// Copyright (c) 2018 Samsung Electronics Co., Ltd All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License
//

#define _GNU_SOURCE
#include <errno.h>
#include <grp.h>
#include <linux/sched.h>
#include <pwd.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mount.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include "config.h"
#include "detach_mount.h"
#include "pivot_root.h"
#include "velen-run.h"

int bind_to_overlay(const char *path) {
  int len = strlen(VELEN_PATH) + strlen("/overlay") + strlen(path) + 1;
  char new_path[len];
  int written = snprintf(new_path, len, "%s%s%s", VELEN_PATH, "/overlay", path);
  if (written < 0) {
    perror("failed to build mount path");
    return EXIT_FAILURE;
  } else if (written != len - 1) {
    fputs("failed to build mount path: too few chars written\n", stderr);
    return EXIT_FAILURE;
  }
  int err = mount(path, new_path, "", MS_BIND, NULL);
  if (err) {
    perror("failed to bind mount");
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}

int switch_user(const char *username) {
  int err;
  struct passwd *lord = getpwnam(username);
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

  return EXIT_SUCCESS;
}

/**
 * @brief Runs a program in Velen.
 *
 * @param argc – length of argv,
 * @param argv – null-terminated array of argv. TODO
 */
int velen_run(int argc, char* argv[]) {
  int err;

  if (argc < 1) {
    fputs("command to sandbox not provided\n"
          "usage: velen-run COMMAND [ARGUMENTS]...\n",
          stderr);
    return EXIT_FAILURE;
  }

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

  // TODO: is this necessary? check for behavioral difference
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

  err = chdir(VELEN_PATH "/overlay");
  if (err) {
    perror("failed to chdir to overlay");
    return EXIT_FAILURE;
  }

  err = pivot_root(".", ".");
  if (err) {
    perror("failed to pivot to overlay");
    return EXIT_FAILURE;
  }

  err = detach_mount("."); // to get rid of old root
  if (err) {
    perror("failed to unmount old root");
    return EXIT_FAILURE;
  }

  err = chdir("/"); // as per man pivot_root(2)
  if (err) {
    perror("failed to chdir to overlay");
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
