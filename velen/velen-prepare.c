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

#define _XOPEN_SOURCE 500
#include <errno.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/mount.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#include "config.h"
#include "velen-prepare.h"

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
