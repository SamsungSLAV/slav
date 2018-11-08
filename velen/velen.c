#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "velen-prepare.c"
#include "velen-run.c"
#include "velen-destroy.c"

int main(int argc, char* argv[], char* envp[]) {
  if (strcmp(argv[0], "velen-prepare") == 0) {
    return velen_destroy() || velen_prepare(&argv[1]);
  }

  if (strcmp(argv[0], "velen-run") == 0) {
    return velen_run(&argv[1], envp);
  }

  if (strcmp(argv[0], "velen-destroy") == 0) {
    return velen_destroy();
  }

  perror("no tool with such name found");
  return EXIT_FAILURE;
}
