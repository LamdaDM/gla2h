# Description
Gla2h is a program for benchmarking argon2 hashes to help the user determine what arguments to pass. 
Uses [tvdburgt's Argon2 library](https://github.com/tvdburgt/go-argon2), which provides bindings for the
reference C implementation of [Argon2](https://github.com/P-H-C/phc-winner-argon2).

You may optionally configure gla2h through these environment variables:
- threads = Number of threads to pass to hash function.
- maxtime = Deadline given to create a hash (in milliseconds).
- runs = Number of runs to do of each setting before taking the average runtime.
- mode = Version of Argon2 to use (2i/2d/2id).

*None of the above can be zero*!

#### **Requires libargon2 and Go**
``` shell
$ go build
```

### Example output:
```
Threads = 8; Maximum Time = 250; Number of runs = 1; Mode = 2i;
MEMORY  PASSES  TIME    
64mb    3       60ms
64mb    4       62ms
64mb    5       78ms
64mb    6       90ms
64mb    7       96ms
64mb    8       109ms
64mb    9       121ms
64mb    10      131ms
64mb    11      143ms
64mb    12      152ms
64mb    13      166ms
64mb    14      180ms
64mb    15      203ms
64mb    16      206ms
64mb    17      212ms
64mb    18      221ms
64mb    19      235ms
128mb   3       93ms
128mb   4       110ms
128mb   5       135ms
128mb   6       171ms
128mb   7       180ms
128mb   8       216ms
128mb   9       225ms
256mb   3       171ms
256mb   4       218ms

Longest runs:
MEMORY  PASS    TIME
64mb    19      235ms
128mb   9       225ms
256mb   4       218ms

```