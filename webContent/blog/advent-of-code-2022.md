---
title: Advent of code 2022 journal
date: 2022-12-01
tags:
  - adventofcode
difficulty: 2
title-image: aoc2022.png
---
It's that wonderful time of the year again when the European programmers wake up at 5 (or 6) am to help a little Elf to save Christmas by solving programming puzzles. **It's Advent of code time!**

This post will be my journal throughout the Advent of code 2022. I'll update this post every day, commenting the puzzle of that day (at least I hope :)).
I'm learning Rust this year, so some of days will be solved in Rust, and the rest in Python.

The reindeers are hungry for _**star fruit**_ found only in the depths of the jungle. We are joining the elves on their annual expedition to the grove where the fruit grows.

## Day 01

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
