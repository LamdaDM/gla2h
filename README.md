# gla2h
Takes 6 arguments:
1) string to be hashed
2) add 'timed' if you want the execution time outputted, else leave as 'x'
3) add 'benchmark' to run through 268 testing iterations of different memory/passes combinations, else leave as 'x'
4) add memory cost in kilobytes, e.g. 262144 (256mb)
5) add passes, e.g. 5 passes (2^5 rounds)
6) add thread count for argon2