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

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/mount.h>

#include "detach_mount.h"

int detach_mount(const char *path) {
  int err = umount2(path, MNT_DETACH);
  if (err) {
    if (fprintf(stderr, "failed to unmount %s: ", path) < 0) {
      perror("failed to print error about failing to unmount");
    } else {
      perror(NULL);
    }
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}
