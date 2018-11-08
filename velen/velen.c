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

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "velen-destroy.h"
#include "velen-prepare.h"
#include "velen-run.h"

int main(int argc, char* argv[]) {
  if (strcmp(argv[0], "velen-prepare") == 0) {
    return velen_prepare(&argv[1]);
  }

  if (strcmp(argv[0], "velen-run") == 0) {
    return velen_run(argc - 1, &argv[1]);
  }

  if (strcmp(argv[0], "velen-destroy") == 0) {
    return velen_destroy();
  }

  if (fprintf(stderr, "no tool with such name found: %s\n", argv[0]) < 0) {
    perror("failed to print error about tool not found");
  }
  return EXIT_FAILURE;
}
