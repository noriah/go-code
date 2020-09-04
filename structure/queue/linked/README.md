# Linked List Queue

A linked list queue uses a series of nodes linked together. Each node holds the value it represents, and a reference to the next node in the queue.

```golang
type Node struct {
  next  *Node
  value interface{}
}
```

Since each node points to the next, you only need to keep track of the top-most node to keep a queue.

In addition, we don't want to lose track of how many nodes we have. That would be embarrasing... So we throw in a `count` value that we update every time the queue changes.

```golang
type Queue struct {
  count int
  root  *Node
  tail  *Node
}
```

For each element added to the queue, we create a new node, holding a value of our element. We point the new node to the queue root node. Then add it to the queue by updating the final node in the queue to point to this new one.

```golang
var newNode = &{
  next: root,
  value: "foobar",
}

var iterNode = root.next

for iterNode.next != root  {
  iterNode = node.next
}

iterNode.next = newNode
```

Some queue implementations (like this one) also hold reference to the last node in the queue. This allows fast enqueue of items, as you do not have to traverse the entire list in order to update the final node.

```golang
var newNode = &Node{
  next: root,
  value: "fast foobar",
}
tail.next = newNode
tail = newNode
```

This has huge performance impact, changing the cost a normal linked list insert from `O(n)` to `O(1)`, where `n` is the size of the queue at time of insert.

Finally, we update the count. If we don't, then we have to traverse the entire list to count the nodes inbetween the first and the last.

```golang
count++
```


The Dequeue cost of a linked list queue is also `O(1)`. We do not have to traverse at all. We simply copy the pointer held by our root, and then update its next value. We decrease count by 1. Finally returning the `value` property of our temproary node.

```golang
var temp = root.next
root.next = temp.next
count--
return temp.value
```

Peeking on a linked list is `O(1)` as the only difference between Dequeue and Peek is that we change a value with Dequeue. With peek, we only return the value held by the top node in our queue.

```golang
return root.next.value
```
