Language symentics
```go
    int a;
    int data = 342 // semicolon is optional
    bool b = true
    bad := false
    string str = "hehe"
    string longStr ="""
        This is Very long string
    """

    int arr[] = {1,2,3,4}
    string strs[] = {"la", "ba", "ka"}

    // struct student {
    //     name string
    //     age int
    // }
    // s1 = new student
    // s1.name = "yash"
    // s1.age = 2 

    if 2 < 10 {
        return true
    } else if a > b {
        return false
    } else {
        return true 
    }

    fn add (int a, int b) int {
        return a + b
    }

    // i probably won't implement match function and instead just go for switch case
    match <expression> with{
          pattern1 -> dofunc1 
        | pattern2 -> dofunc2 
        | pattern2 -> dofunc3
        | _        -> doDefualtFunc // this is defualt case
    }
    // and nah there won't be any for/while loops just use good old recursion 

    /* Also both single line and multi/inline line commenting is allowed *\

```