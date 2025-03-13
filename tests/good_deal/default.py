#!/usr/bin/python3
while True:
    R = int(input())
    T = int(input())
    role = input()
    if role == "A":
        m = R//3+1
        print(m)
        ans = int(input())
    else:
        m = int(input())
        print(1) if m > R//2 else print(0)