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
