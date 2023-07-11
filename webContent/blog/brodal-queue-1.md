---
title: Brodal queue - a high-level overview
date: 2023-07-12
tags: 
  - compsci
  - math
  - brodalQ
difficulty: 2
latex: true
---

This is a topic I'm passionate about[^1] and I think it would be fun to write about and introduce it to others.

[^1]: Probably because it was a part of my master's thesis :smile:

As the Wiki says: 
> Brodal queue is a heap-like [priority queue](https://en.wikipedia.org/wiki/Priority_queue) with very low worst case time bounds: $O(1)$ for insertion, 
> find-min, meld and decrease-key, and $O(log\\;n)$ for deletion.

So, as you can see, one of the cool things about the Brodal queue is the low time-bound[^2].
But, the implementation of the Brodal queue and the role it plays in achieving that time-bound is way more interesting.

[^2]: I think it's the lowest possible time bound for heap-like priority queues.

I'll write a couple of posts about this topic because I want to try to explain as much as possible, and this post would be way too long if I tried to cram it all in here.
This post will take a high-level look at this data structure. It'll also try to explain how all priority queue methods should work.
I won't get into any implementation specifics here. I'll cover that in a different post.

## The structure
For all of this to make some sense, I need to list out some rules and definitions.
This structure is designed with three sets of rules. 
I'll list a couple of them here so that you can follow along and I'll list all three sets in a later post.

### Definitions
The Brodal queue is defined as a pair of root nodes $(t_{1},t_{2})$, each representing a tree $T_{i}, i \in \\{1,2\\}$[^3]. 

[^3]: Getting ahead of myself... We need this so that we can have the min element in a root and be able to merge two queues in $O(1)$ time.
This is an oversimplification and will be properly explained in a post focusing on the implementation details.

Nodes will be allowed to break the heap property. 
Those kinds of nodes are called **violating** nodes and are tracked.
There will always be $O(log\\;n)$ of violating nodes and we'll have a method of getting rid of them.

The rank of a node is a number we relate to each node. It is defined in the [rules](#rules) section.

### Rules
1. $T_{1}$ root node must always be the "minimal" element of the queue.
2. The rank of a child node is always smaller than the rank of its parent. The rank of a leaf node is 0.
3. A parent of rank $r$ must have at least two children of rank $r-1$.
4. Number of children of rank $i$ must be less than or equal to $7$.
5. The rank of the $t_{2}$ node is greater than or equal to the rank of the $t_{1}$ node or $T_{2}$ is an empty tree.

Key points to take away from this set of rules are: we know that each node of rank $r$ has at least two child nodes of rank $r-1$ and will always have less than $8$ children of ranks $\\{0, \ldots, r-1\\}$.
This is crucial for achieving the constant-time balancing of the trees, making it crucial for achieving low worst-case time bounds.
A node of rank $r$ would look something like this:

![An example of a node of rank r in a Brodal queue](/static/img/brodal-01-white.svg)

With all that in mind, we can start describing how this implementation works.
	
## Merging queues
This is a simple method. 
We need to find the node with the smallest value and the node with the highest rank.
If the node with the smallest value is also the node with the highest rank, then the new priority queue will have only one tree, the $T_1$ tree.
If the node with the smallest value isn't the node with the highest rank, then the node with the smallest value becomes the new $t_1$ node, and the node with the highest rank becomes the new $t_2$ node.

Merging two trees is done by inserting the tree with the lower rank as a child of the root node of the other tree.
This can be done in constant time and will be explained as a part of the implementation.

## Deleting the $t_1$ node
This is the only method that doesn't have the $O(1)$ time bound[^4].
One of the reasons for this is the search for the new "minimal" element.
There are a couple of search spaces for this new minimum. 
One of them is the **direct** children of the old minimum; the other is the set of violating nodes.
Both of these search spaces are $O(log\\;n)$ in size.
This is fairly obvious for the children[^5], but not so obvious for the violating nodes.
We'll explain why there are $O(log\\;n)$ violating nodes in a later post.

[^4]: Technically, delete is too, but since delete is implemented as a composition of `decreaseKey` and `deleteMin`, I don't count it as a separate method.
[^5]: This is a direct consequence of the definition of the rank of a node and rule number 4.

We start by cutting all children off of $T_2$ and making them and $t_2$ direct children of $t_1$.
We search for the new $t_1$ node in the search space comprised of the direct children of $t_1$ and the violating nodes[^6].
After finding the new $t_1$ node, we need to make it the new $t_1$.
We also need to convert all violating nodes into good nodes by making them children of $t_1$.

[^6]: Notice that the intersection of these two sets of nodes is empty.

## Balancing a tree
You should have noticed that the merging and deleting methods require the ability to add/remove a child to a node in constant time.
This is a nontrivial task because we could break some rules listed in the [rules](#rules) section, like rules 3 and 4.

We solve this by implementing a helper data structure called **the guide**.
Think of it as a black box that takes an array of numbers representing the numbers of children of rank $i$ that a node has, and outputs a list of actions we need to do to rebalance a tree in the queue.
This list of actions is guaranteed to be of a constant length, and the actions can be performed in constant time.

## Resources
Here is the link to the paper on Brodal queues: [https://cs.au.dk/~gerth/papers/soda96.pdf](https://cs.au.dk/~gerth/papers/soda96.pdf).

If you have any questions or comments on this post, please send them to me at [mislavzanic3@gmail.com](mailto:mislavzanic3@gmail.com). 
Thank you for reading!
