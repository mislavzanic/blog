---
title: Advent of code 2022 journal
date: 2023-04-10
tags:
  - adventofcode
  - math
difficulty: 2
title-image: aoc2022.png
---
It's that wonderful time of the year again when the European programmers wake up at 5 (or 6) am to help a little Elf to save Christmas by solving programming puzzles. **It's Advent of code time!**

This post will be my journal throughout the Advent of code 2022. I'll update this post every day, commenting the puzzle of that day (at least I hope :)).
I'm learning Rust this year, so some of days will be solved in Rust, and the rest in Python.

The reindeers are hungry for _**star fruit**_ found only in the depths of the jungle. We are joining the elves on their annual expedition to the grove where the fruit grows.

![AOC](/post/images/aoc2022_1.png)

Since the first couple of puzzles will be nothing to write home to, I'll update this post with more interesting puzzles.

## Table of contents
1. [Day 01](#d01)
2. [Day 02](#d02)
3. [Day 03](#d03)
4. [Rest of the days](#rest)

## <a class="inpost" name="d01">Day 01</a>

Pretty easy day (as per usual).
The elves are carrying food, and each food item has a number assigned to it (puzzle input).
Input looks like this:
```
1000
2000
3000

4000

5000
6000

7000
8000
9000

10000
```
Basically, each Elf separates their own inventory from the previous Elf's inventory (if any) by a blank line.
We need to find the maximal sum of items for part one, and sum of three largest sums of items for part two.

```rust
use std::fs;

fn day01() {
    let input: String = fs::read_to_string("d01.in")
        .unwrap()
        .trim_end_matches(&['\r', '\n'])
        .to_string();

    let mut input: Vec<i32> = input
        .split("\n\n")
        .map(|s| s
             .split("\n")
             .map(|x| x
                  .parse::<i32>()
                  .unwrap())
             .sum())
        .collect();

    input.sort();
    input.reverse();

    println!("{}, {}", input[0], input[0..3].to_vec().into_iter().sum::<i32>());
}
```

## <a class="inpost" name="d02">Day 02</a>

Now for a bit of rock, paper, scissors.
Pretty easy day, lots of typing and HashMaps :).

```rust
use std::collections::HashMap;
use itertools::{Itertools};
use std::fs;

pub fn solve() {
    let binding = fs::read_to_string("d02.in")
        .unwrap()
        .trim_end_matches(&['\r', '\n'])
        .to_string();

    let guide: Vec<(&str, &str)> = binding.split("\n")
        .map(|s| s.split(" ").next_tuple().unwrap())
        .collect();

    let part1: HashMap<(&str, &str), i32> = HashMap::from([
        (("X", "C"), 7),(("Y", "A"), 8),(("Z", "B"), 9),
        (("X", "A"), 4),(("Y", "B"), 5),(("Z", "C"), 6),
        (("X", "B"), 1),(("Y", "C"), 2),(("Z", "A"), 3)
    ]);

    let part2: HashMap<(&str, &str), i32> = HashMap::from([
        (("X", "C"), 2),(("Y", "A"), 4),(("Z", "B"), 9),
        (("X", "A"), 3),(("Y", "B"), 5),(("Z", "C"), 7),
        (("X", "B"), 1),(("Y", "C"), 6),(("Z", "A"), 8)
    ]);

    let mut p1 = 0;
    let mut p2 = 0;
    for (other, me) in guide.iter() {
        p1 += part1[&(*me, *other)];
        p2 += part2[&(*me, *other)];
    };

    println!("{}, {}", p1, p2)
}
```

## <a class="inpost" name="d03">Day 03</a>

A fun day.
We needed to find a couple of intersections.
For part one, we needed to find the intersection between the left and right-hand sides of strings, and, for part two, between batches of three strings.
Excuse the poor rust code, I'm still learning :).

```rust
use std::collections::HashSet;
use std::fs;

use itertools::Itertools;

fn prio(rune: &char) -> u32 {
    if rune.is_uppercase() {
        (*rune as u32) - ('A' as u32) + 27
    } else {
        (*rune as u32) - ('a' as u32) + 1
    }
}

pub fn solve() {
    let binding = fs::read_to_string("d03.in").unwrap();

    let mut part1 = 0;
    let mut part2 = 0;
    let mut batch: HashSet<char> = HashSet::new();

    for (idx, item) in binding.trim_end_matches(&['\r', '\n']).split('\n').enumerate() {
        if idx % 3 != 0 {
            let s: HashSet<char> = item.chars().collect();
            batch = batch.into_iter().filter(|c| s.contains(c)).collect();

            if idx % 3 == 2 {
                part2 += batch.iter().map(|c| prio(c)).sum::<u32>();
            }

        } else {
            batch.clear();
            batch.extend(&HashSet::from(item.chars().collect::<HashSet<char>>()));
        }

        let len: usize = item.chars().count();
        let (p1, p2) = item.split_at(len / 2);
        let (s1, s2): (HashSet<char>, HashSet<char>) = (
            p1.chars().collect(), p2.chars().collect()
        );
        part1 += s1.intersection(&s2).map(|c| prio(c)).sum::<u32>()
    }

    println!("{}, {}", part1, part2);
}
```

## <a class="inpost" name="rest">Rest of the days</a>

Sooo.....
The journal went as planned (obviously :), since it's the day after Easter)...

I struggled with Rust for another day or two, and then I switched over to Python :).
The days weren't anything special this year. The most notable ones are [#16](https://adventofcode.com/2022/day/16) and [#22](https://adventofcode.com/2022/day/22).

### Day 16
This one was legit fun. Like, really fun!
Also, this one holds my best rank (for part one only (it's 164)).

Part one was a depth-first search through a graph. 
The catch was that the graph had lots of nodes that had the weight 0, so we could use the [Floyd-Warshall algorithm](https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm) to create a **complete** graph.
Also, the time variable $t$ speeds up the search.


Part two was not so easy at first glance, but, after some thinking, you come to realize that it's just part one ran twice for two parts of the **complete** graph from part one.
Time variable $t$ does some heavy lifting in this part!


### Day 22
This one is plain evil.
I'm not even going to explain it (since there is no explaining needed!). Just hardcode the cube and be done with it...



## Final thoughts
This year had the best puzzle I ever solved ([#16](https://adventofcode.com/2022/day/16)), and the worst ([#22](https://adventofcode.com/2022/day/22)).
All in all, it was really fun, and I hope it will be as fun in 2023.

