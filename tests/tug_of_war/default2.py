#!/usr/bin/python3
n = 100 # кол-во сил
m = int(input()) # кол-во раундов

while True: 
    k = int(input()) 
    if k >= 0:
        planed_to_spend = k + 1 # например, хотим тратить на 1 ед. больше, чем соперник всегда
        print(planed_to_spend)
    else:
        m -= 1 # раунд закончился, проверяющая программа прислала -1