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

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include "config.h"
#include "detach_mount.h"
#include "velen-destroy.h"

static int remove_directory(const char* path) {
  int err = rmdir(path);
  if (err) {
    fprintf(stderr, "failed to remove directory %s: %m\n", path);
    return EXIT_FAILURE;
  }
  return EXIT_SUCCESS;
}

int velen_destroy() {
  int status = 0;
  status |= detach_mount(VELEN_PATH "/overlay");
  status |= detach_mount(VELEN_PATH "/top");
  status |= remove_directory(VELEN_PATH "/overlay");
  status |= remove_directory(VELEN_PATH "/top");
  status |= remove_directory(VELEN_PATH);
  return status;
}
