# 进程API

## 获取pid和父进程的pid（ppid）

**获取进程ID**

```c
#include <sys/types.h>
#include <unistd.h>

pid_t getpid(void);
```

**获取父进程ID**

```c
#include <sys/types.h>
#include <unistd.h>

pid_t getppid(void);
```

## fork系统调用

通过fork系统调用，可以创建一个和当前进程一样的进程。

```c
#include <sys/types.h>
#include <unistd.h>

pid_t fork(void);
```

fork调用成功，会创建一个新的进程。
调用者在创建了新进程之后依然会正常执行。

父进程和子进程在fork调用完成之后会执行同样的程序。
它们的不同在于在父进程中，fork会返回子进程的pid。
而子进程中的fork会返回0.

父进程和子进程的区别：
* pid不同
* 子进程的ppid是父进程的pid
* 子进程的资源统计信息被清零
* 父进程的信号不会被子进程继承
* 父进程的文件锁不会被子进程继承

出错返回-1，并设置errno的值：
* EAGAIN：内核申请资源失败
* ENOMEM：内核内存不足

```c
{{ #include ./fork.c }}
```

## wait系统调用
wait()系统调用用于等待进程和终止。

```c
#include <sys/types.h>
#include <sys/wait.h>

pid_t wait(int * status);
```

如果子进程没有终止，调用会阻塞；
调用成功，返回已终止的子进程的pid；
出错则返回-1，并设置errno的值：
* ECHILD：调用进程没有任何子进程
* EINTR：在等待子进程结束时收到信号，调用提前返回

**status**
返回子进程的附加信息。
这些信息用比特位来表示。可以用下面的这些宏来解释：

```c
#include <sys/wait.h>

int WIFEXITED(status);      // 程序正常结束（调用__exit()）
int WIFSIGNALED(status);    // 因为信号导致进程终止
int WIFSTOPPED(status);     // ptrace调试情况下，进程停止
int WIFCONTINUED(status);   // ptrace调试情况下，进程继续执行

int WEXITSTATUS(status);    // 正常退出情况下（WIFEXITED），进程的返回值
int WTERMSIG(status);       // 信号终止情况下（WIFSIGNALED），导致终止的信号编号
int WSTOPSIG(status);       // 导致信号停止的信号编号
int WCOREDUMP(status);      // 信号终止的情况下（WIFSIGNALED），生成coredump文件
```

**等待特定进程**

```c
#include <sys/types.h>
#include <sys/wait.h>

pid_t waitpid(pid_t pid, int * status, int options);
```
options是一个flag：
* WNOHANG：不要阻塞
* WUNTRACED：即使调用进程没有跟踪子进程，也会设置WIFSTOPPED位；对shell有帮助
* WCONTINUED：即使调用进程没有跟踪子进程，也会设置WIFCONTINUED位；对shell有帮助

errno多了一个`EINVAL`表示参数非法

**僵尸进程**

如果子进程在父进程之前终止，内核会把该进程设置成特殊的进程状态。
处于这种状态的进程称为僵尸进程。

这是为了让父进程获取子进程的状态而设计的，如果子进程终止，就消失。
那么，父进程无法获取到任何关于子进程的信息。

僵尸进程只会保存可能有用的信息。僵尸进程会等待父进程来查询状态，只
有当父进程查询到已终止的子进程的状态的时候，这个子进程才会消失，不
再处于僵尸状态。

```c
{{ #include ./wait.c }}
```

## exec系统调用

