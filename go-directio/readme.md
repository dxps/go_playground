## Using DirectIO with Go

The purpose of this playground is to evaluate the optios of Nick Craig-Wood's `directio` library.

### Notes

Producer is not appending to existing file, but writes from the start of it.
```shell
## Context: There are 4 items in the file, produced before the consumer started.

## Consumer starts, got these 4 items initially, then it sleeps (shown by the dots), waiting for new entries. Which eventually are consumed, after being produced, but already read blocks are ignored, of course.
$ ./run_consumer.sh
2020/12/20 16:09:56 Starting to consume ...
2020/12/20 16:09:56 Ready to read.
..2020/12/20 16:09:57 Got {Value:1}
2020/12/20 16:09:57 Got {Value:2}
2020/12/20 16:09:57 Got {Value:3}
2020/12/20 16:09:57 Got {Value:4}
............2020/12/20 16:10:09 Got {Value:5}
..2020/12/20 16:10:11 Got {Value:6}
2020/12/20 16:10:11 Got {Value:7}

## While the consumer is waiting for new entries in the file,
## producer gets started and it adds again items starting with value 0.
$ ./run_producer.sh
2020/12/20 16:10:03 Starting to produce ...
2020/12/20 16:10:03 Produced {Value:1}
2020/12/20 16:10:03 Ready to write.
..2020/12/20 16:10:04 Produced {Value:2}
2020/12/20 16:10:05 Produced {Value:3}
..2020/12/20 16:10:06 Produced {Value:4}
2020/12/20 16:10:07 Produced {Value:5}
..2020/12/20 16:10:08 Produced {Value:6}
2020/12/20 16:10:09 Produced {Value:7}

## And as you can see above, the consumer "detects" the new entries, meaning the new blocks added after the initial 4 ones already read before by the consumer.
```

### TODOs

- [ ] Externalize the config items: <br/>
    - `BlockSize` - used for writing and reading blocks of serialized data
    - `ExchangeFile` - full path to the file used for writing into and reading from
