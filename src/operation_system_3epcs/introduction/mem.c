#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <assert.h>

int main()
{
    int * p = (int*)malloc(sizeof(int));

    assert(p != NULL);

    printf("(%d) memory address of p:0x%p\n", getpid(), p);
    return 0;
}
