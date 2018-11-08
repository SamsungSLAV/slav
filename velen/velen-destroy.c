#include <stdlib.h>

#include "velen-destroy.h"
#include "detach_mount.h"
#include "config.h"

int velen_destroy() {
  return detach_mount(VELEN_PATH "/overlay")
    || detach_mount(VELEN_PATH "/top");
}
