# Slice Queue

A slice queue uses an internal array and channels to keep track of data and order of insert

```
Queue is structure of
    data is type array
    push-channel is type channel
    pop-channel is type channel
    next-push is type integer
    next-pop is type integer
    size is type integer
```

We keep two buffered channels that are the same size as the array. One for pop indexes and one for push indexes.

Each time we push a value to the end of the queue, we are actually inserting it at whatever index is available first:

```
function getPushIndex()
    let index = next-push

    if push-channel is empty
        next-push = next-push + 1
    else
        index = push-channel.next()

    return index
```

Then we can push add the value to the array in our queue and add the index to our pop-channel

```
function push(value)
    let index = getPushIndex()
    data[index] = value
    pop-channel.add(index)
```

When we peek at or pop the front of the queue, we do so by getting the index for the value from the pop-channel. To support peek-ing multiple times and pop-ing after a peek, we have to keep track of what index is removed from the pop-channel:

```
function getPopIndex(peeking)
    let index = next-pop

    if index < 0
        if pop-channel is empty
            throw error(empty queue)
        else
            index = pop-channel.next()

    if peeking
        next-pop = index
    else
        next-pop = -1

    return index
```

Then we can return our value and if we are pop-ing instead of peek-ing, add the index to the push-channel for reuse

```
function pop()
    let index = getPopIndex(false)
    push-channel.add(index)
    return data[index]
```

Peek-ing is the same as pop-ing, however we don't add the index to our push-channel

```
function peek()
    let index = getPopIndex(true)
    return data[index]
```
