# memStress

This is a tool to simulate memory allocation.

## Usage:

```
>memStress -h

Usage of ./memStress:
  -client
        the process runs as a client
  -size string
        size of memory you want to allocate (default "0KB")
  -time string
        time to reach the size of memory you allocated (default "0s")
  -workers int
        number of workers allocating memory (default 1)
```

You can generate a model that simulates a memory usage like `memStress --size 1GiB --time 1m	--workers 2`. This command will generate two workers, each of which will allocate 1GiB of memory and the **allocation process** will last 1 minute.

(notice: The allocation process will last 1minute, not the program runs for 1 minute.)

## Build:

`make build`





