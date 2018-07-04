retry
====

### Usage

```
$ retry -h

Usage: retry [OPTIONS] COMMAND [ARG...]

Repeat command

Options:
  -d int
    	set delay seconds (default 10)
  -h	show help
  -n int
    	set number of max retry count (default 3)
  -s	silent mode
  -v	show version

```

### Example

```
$ retry -n 3 -d 1 false
exit status 1, retry after 1 seconds.
exit status 1, retry after 1 seconds.
exit status 1, retry after 1 seconds.
$ echo $?
1

$ retry -n 3 -d 1 true
$ echo $?
0

$ echo hoge | retry -d 1 bash -c 'cat && false'
hoge
exit status 1, retry after 1 seconds.
hoge
exit status 1, retry after 1 seconds.
hoge
exit status 1, retry after 1 seconds.
```
