# lim-lang

```c
    int data1 = 3;
    pub int data2 = 4;
    int dataExp = 2 * 4 + 123;
    bool databool = true;
    float data = 3.14;
    string dataStr = "data";
    dataStr = "newData";

    if 5 < 10 {
        print("worked if");
    } else if data1 != data2 {
        print("worked else if data1 != data2 ")
    } else if data1 >= data2 {
        print("worked else if data1 >= data2 ")
    }else {
        print("worked else");
    }

    fn printNumber(int num){
        print("This is the number",num);
    }

    pub fn printNextNumber(int num){
        print("This is the next number",num + 1);
    }

    int arr[5] = {1,2,3,4,5};
    string arrStr[3] = {"abc", "123", "aaa"};

    int arrAdr = &arr;
    print(arrAdr);
    print(*arrAdr[0])

```