#!/usr/bin/python3
k = 0
c = 1
while True:
    R = int(input())
    T = int(input())
    role = input()
    if role == "A":
        print(R * k // 100)
        ans = int(input())
        if ans == 1 and k > 0:
            k -= 5
        elif ans == 0 and k < 90:
            k += 10
    else:
        m = int(input())
        if m < R//2:
            k = 0
        print(1)
