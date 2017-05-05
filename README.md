# Timetravel  [![CircleCI](https://circleci.com/gh/radekdymacz/timetravel.svg?style=svg)](https://circleci.com/gh/radekdymacz/timetravel)

Simple tool to travel in time.

## What does it do

Batch change modified time of files in given volume/folder if it's in the future.

## Why

Its quite a common occurrence for files to have corrupted metadata timestamps like access, modified and created time in the future. This impacts availability to properly backup or archive these files hence the need for the tool.


## Usage

```
timetravel /mnt/volume/folder
```

## Status

WIP
