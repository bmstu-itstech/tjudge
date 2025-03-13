#!/usr/bin/python3
n = 100 
m = int(input()) 
spent = 0
per_game = n // m

while True: 
    k = int(input()) 
    if k < 0: 
        n -= spent
        m -= 1
        spent = 0
        per_game = n // m
        continue

    planed_to_spend = k + 1
    if (spent + planed_to_spend <= per_game): 
        print(planed_to_spend)
        spent += planed_to_spend
    else:
        n -= spent
        m -= 1
        spent = 0
        if m > 0:
            per_game = n // m
        else:
            per_game = n
        print(-1)