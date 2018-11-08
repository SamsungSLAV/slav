#include <assert.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

#include "velen-prepare.h"
#include "velen-run.h"
#include "velen-destroy.h"

int main(int argc, char* argv[]) {
  assert(argc > 0);
  if (strcmp(argv[0], "velen-prepare") == 0) {
    return velen_destroy() || velen_prepare(&argv[1]);
  }

  if (strcmp(argv[0], "velen-run") == 0) {
    return velen_run(&argv[1]);
  }

  if (strcmp(argv[0], "velen-destroy") == 0) {
    return velen_destroy();
  }

  if (fprintf(stderr, "no tool with such name found: %s\n", argv[0]) < 0) {
    perror("failed to print error about tool not found");
  }
  return EXIT_FAILURE;
}
