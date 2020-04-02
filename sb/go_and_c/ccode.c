
#include "ccode.h"

#include <stdio.h>

int my_func(int x)
{
    int i;

    for(i = 0; i < x; i++)
    {
        printf("Hello world from the world of C\n");
    }
}
