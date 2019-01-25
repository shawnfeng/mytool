#! /usr/bin/env python
# -*- coding: utf-8 -*-


import random

def lowest_common_multiple(x, y):
    greater = x
    if x < y:
       greater = y
 
    lcm = greater
    while True:
        if greater % x == 0 and greater % y == 0:
            lcm = greater
            break
        greater += 1
 
    return lcm

def disp(days, user0, user1):
    lcm = lowest_common_multiple(len(user0), len(user1))
    job0 = []
    job1 = []


    for i in range(lcm/len(user0)):
        random.shuffle(user0)
        job0 += user0


    for i in range(lcm/len(user1)):
        random.shuffle(user1)
        job1 += user1

    print len(job0), len(job1)


    print "日期 第一值班人 第二值班人"
    for i, d in enumerate(days[:lcm]):
        print "%s %s %s" % (d, job0[i], job1[i])


def doit():
    user0 = [
        "邬勇",
        "何菱",
        "胡小平",
        "李习怀",
        "陶文亮",
        "卞绍雷",
        ]
    user1 = [
        "李桐",
        "王庆祝",
        "陈现麟",
        "冯广祥",
        ]

    days = [
        "10-01",
        "10-02",
        "10-03",
        "10-04",
        "10-05",
        "10-06",
        "10-07",
        "10-13",
        "10-14",
        "10-20",
        "10-21",
        "10-27",
        "10-28",
        ];



    disp(days, user0, user1)

if __name__ == '__main__':
    doit()

