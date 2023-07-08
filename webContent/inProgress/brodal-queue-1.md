---
title: Brodal queue - how and why does it work
date: 2023-07-09
tags: 
  - compsci
  - math
  - brodalQ
difficulty: 4
latex: true
---

This is a topic I'm passionate about[^1] and I think it would be fun to write about and introduce it to others.

[^1]: Probably because it was a part of my master thesis :smile:

As the Wiki says: 
> Brodal queue is a heap-like [priority queue](https://en.wikipedia.org/wiki/Priority_queue) with very low worst case time bounds: $O(1)$ for insertion, 
> find-min, meld and decrease-key, and $O(log n)$ for deletion.

So, as you can see, one of the cool things about a Brodal queue is the low time-bound[^2].
But, the implementation of a Brodal queue and the role it plays in achieving that time-bound is way more interesting.

[^2]: It's the lowest possible time bound for heap-like priority queues.

I'll write a couple posts about this topic because I want to try to explain as much as possible, and this post would be way to long if I tried to cram it all in here.
This post will take a high level look at this data structure. It'll also try to explain how all priority queue methods should work.
I won't get into any implementation specifics here. I'll cover that in a different post.

## The structure

Brodal queue is defined as a pair of root nodes $(t_{1},t_{2})$.
