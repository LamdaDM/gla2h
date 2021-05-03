# gla2h
Takes 6 arguments: 
1) input -> string that will be hashed.
2) timer -> flag to time how long the hash takes, leave as 'y' if you want the execution time outputted, else leave as 'x' or 'n'.
3) benchmark -> 'y' to run through 268 testing iterations of different memory/passes combinations, else leave as 'x' or 'n'.
4) memory cost -> memory to give to the threadpool in megabytes, e.g. '256' (256*1024 kb).
5) passes -> how many passes or iterations to do, e.g. '10' (2^10 passes).
6) parallelism -> how many threads to give argon2, usually double the cores the host's CPU has.
