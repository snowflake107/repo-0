#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <sys/wait.h>

int main(int argc, char * argv[])
{
    int rc = fork();
    if (rc < 0)
    {
        fprintf(stderr, "fork failed.\n");
        exit(1);
    }
    else if (rc == 0)
    {
        close(STDOUT_FILENO);   
        // 关掉文件描述符STDOUT_FILENO，
        // 下一个打开的文件会从小的文件描述符开始找
        // 给open的文件的描述符就会是STDOUT_FILENO。相当于将STDOUT重定向到文件中。
        open("./wc.ouput", O_CREAT | O_WRONLY | O_TRUNC, S_IRWXU);

        char * myargs[3];
        myargs[0] = strdup("wc");
        myargs[1] = strdup("wc.c");
        myargs[2] = NULL;
        execvp(myargs[0], myargs);
    }
    else
    {
        wait(NULL);
    }

    return 0;
}
