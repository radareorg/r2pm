#include "libr2pm.h"

int main() {
    R2pmSetDebug(1);

    char* path = "/home/quentin/.local/share/radare2/r2pm/";

    R2pmInit(path);
    R2pmDelete(path);
}
