# Using DirectIO with Go

The purpose of this playground is to evaluate the optios of Nick Craig-Wood's [`directio`](https://github.com/ncw/directio) library.

## Features

This section describes the features implemented by both parties: Producer and Consumer.

### Configuration

The configuration items used by both parties are stored in `.env` file. You'll find there further details about each item's purpose.

### Producer

Accoring to the max file size (defined in `IO_FILE_MAX_SIZE` config item), Producer is:
- appending to the latest written file, if that file's size didn't reach the max size
- writing to a new file, if there is no written file or if the size of the latest one exceeded the max size

### Consumer

Consumer:
- reads (_consumes_) files - one by one - from the same path (define in `IO_PATH` config item)
- saves the state (aka `ConsumerState` in the consumer's code), so that it can resume the work any time

## Todos

- If writer fails to init properly, main should get notified, stop the producer, and end itself.
    - That can happens when the file system where the file resides does not support O_DIRECT flag.<br/>
      See these [notes on linux kernel and O_DIRECT](https://lists.archive.carbon60.com/linux/kernel/720702).

