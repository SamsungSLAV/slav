// Copyright (c) 2019 Samsung Electronics Co., Ltd All Rights Reserved
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
#include <string.h>
#include <sys/types.h>
#include <pwd.h>

#include "config.h"
#include "velen-shell.h"
#include "velen-run.h"

int velen_shell() {
  struct passwd *user = getpwnam(VELEN_LORD);
  if (user == NULL) {
    perror("failed to obtain user information");
    return EXIT_FAILURE;
  }

  // velen_run invokes getpwnam, which according to its manpage
  // may or may not overwrite the contents of a previously returned
  // struct.
  char shell[strlen(user->pw_shell)];
  strcpy(shell, user->pw_shell);

  char* argv[] = {shell[0] != '\0' ? shell : VELEN_DEFAULT_SHELL, NULL};

  return velen_run(1, argv);
}

