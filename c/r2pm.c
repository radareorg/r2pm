#include <stdio.h>
#include <string.h>

#include "libr2pm.h"

int main(int argc, char* argv[]) {
    if (argc < 2) {
        fprintf(stderr, "USAGE: %s <command>\n", argv[0]);
        return EXIT_FAILURE;
    }

    int ret;

    char* path = "r2pm_c-test-repo";
    char* command = argv[1];

    if (strcmp(command, "init") == 0) {
        if ((ret = r2pm_init(path)) != 0) {
            perror("init");
            return ret;
        }
    } else if (strcmp(command, "list-available") == 0) {
        struct r2pm_string_list* p;

        if ((ret = r2pm_list_available(path, &p)) != 0) {
            perror("list_available");
            return ret;
        }

        if (p == NULL) {
            fprintf(stderr, "p is null");
        } else {
            struct r2pm_string_list* current = p;

            while(current != NULL) {
                printf("%s\n", current->s);
                current = current->next;
            }
        }
    } else if (strcmp(command, "delete") == 0) {
        if ((ret = r2pm_delete(path)) != 0) {
            perror("delete");
            return ret;
        }
    } else {
        fprintf(stderr, "%s: unknown command\n", command);
        return EXIT_FAILURE;
    }
}
