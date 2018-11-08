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

int velen_shell(int argc, char* argv[]) {
  argv[0] = VELEN_FORCED_SHELL;
  return velen_run(argc, argv);
}

#ifdef VELEN_SHELL_BUILD
int main(int argc, char* argv[]) {
  return velen_shell(argc, argv);
}
#endif // VELEN_SHELL_BUILD
