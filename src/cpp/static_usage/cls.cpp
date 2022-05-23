class A
{
public:
    static int i;
    int j;

    static int func_static() { return i; }

    int func() { return j; }
};

int A::i = 20;

int bar()
{
    return A::func_static();
}

int foo()
{
    A a;
    a.j = 10;
    return a.func();
}
