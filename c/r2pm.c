#include "libr2pm.h"

int main() {
    R2pmSetDebug(1);

    char* path = "r2pm_c-test-repo";

    R2pmInit(path);
    R2pmDelete(path);
}
