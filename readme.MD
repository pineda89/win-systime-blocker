simple example of how modify the call to kernel32 api on any process

in this case, we are modifying the opcodes. Instead calling the kernel32 function, is returning the desired value.

![Alt text](original.png?raw=true "original value")

![Alt text](modified.png?raw=true "modified value")


```
output:

current time: 1922231491 30840826 132460340973857987
current time parsed 2020-10-01 13:54:57.3857987 +0000 UTC
modified time: 1173307423 47618039 204517921377959967
modified time parsed 2249-02-03 13:22:17.7959967 +0000 UTC
```