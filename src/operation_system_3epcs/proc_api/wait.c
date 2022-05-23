#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/wait.h>

int main()
{
    int status;

    printf("Parent process: %d\n", (int)getpid());

    pid_t pid = fork();
    if (pid == -1)
    {
        fprintf(stderr, "fork failed.\n");
        exit(1);
    }
    else if (pid == 0)
    {
        printf("Child process starts.\n");
        sleep(60);
        exit(123);
    }
    else
    {
        sleep(120); // for observing zombie proc status
        pid = wait(&status);
    }
    printf("Child process %d exit.\n", (int)pid);

    if (WIFEXITED(status))  // wait 120 seconds, exec `ps ax` after 60 seconds, the proc status turn into `Z+`.
        printf("Normal termination with exit status: %d\n",
            WEXITSTATUS(status));

    if (WIFSIGNALED(status))    // manually kill child process
        printf("Killed by signed: %d%s\n",
            WTERMSIG(status),
            WCOREDUMP(status) ? " (dumped core)" : "");

    if (WIFSTOPPED(status))
        printf("Stopped by signal=%d\n",
            WSTOPSIG(status));
    
    if (WIFCONTINUED(status))
        printf("Continued\n");
    
    return 0;
}
