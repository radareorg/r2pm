#include <stdio.h>

#include "libr2pm.h"

int main() {
    r2pm_set_debug(1);

    char* path = "r2pm_c-test-repo";

    int ret;

    if ((ret = r2pm_init(path)) != 0) {
        perror("init");
        return ret;
    }

//    if ((ret = r2pm_install(path, "r2dec")) != 0) {
//        perror("install");
//        return ret;
//    }
//
//    struct r2pm_string_list* p;
//
//    if ((ret = r2pm_list_available(path, &p)) != 0) {
//        perror("list_available");
//        return ret;
//    }
//
//    if (p == NULL) {
//        fprintf(stderr, "p is null");
//    }

    if ((ret = r2pm_delete(path)) != 0) {
        perror("delete");
        return ret;
    }
}
