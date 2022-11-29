---
title: Advent Of Code, 2019, day 22
date: 2022-11-13
tags: 
  - adventofcode
  - math
title-image: aoc.jpg
latex: true
difficulty: 2
---
I love Advent of code.
It's an annual event for everyone who likes coding and problem-solving or wants to learn those two things.
A two-part coding puzzle is posted every day at midnight EST (UTC-5) on [adventofcode.com](https://adventofcode.com "adventofcode.com").
Some of the puzzles can be very interesting.
Such is the puzzle I'm going to be talking about today.

You can try to solve the puzzle [here](https://adventofcode.com/2019/day/22 "https://adventofcode.com/2019/day/22").

## Solution
I'll suppose that you read the puzzle.
This solution is written in python.

You're given a deck of $p \in \mathbb{P}$ cards, where $\mathbb{P}$ is the set of all prime numbers (the prime number stuff is important for later).
You're also given three types of shuffles with which you shuffle the deck.
Since we only need the new position of the card denoted by $x$ we can try to implement all of the shuffles as numeric functions.

*Deal into new stack* can be implemented as $$f(x) = p - 1 - x.$$

*Cut N* can be implemented as $$f(x, N) = (x - N)\; (mod\; p).$$

*Deal with increment N* can be implemented as $$f(x, N) = (x \cdot N)\; (mod\; p).$$

Or, in python:
```python
def cut(N, x, p):
    return (x - N) % p

def deal_with_increment(N, x, p):
    return (x * N) % p

def deal_new_stack(x, p):
    return p - 1 - x
```

### Part 1

This part is pretty straightforward.
We get a deck with $10007$ cards and a card in position $2019$. We need to calculate the position of our card after the shuffle process (the input).

We just apply our shuffle process on the card and, we get the solution.

Code:
```python
def main():
    shuffle = [x.strip() for x in open('input').readlines()]

    def shuffle_cards(card_pos, num_of_cards):
        for s in shuffle:
            if s == "deal into new stack":
                card_pos = deal_new_stack(card_pos, num_of_cards)
                continue

            N = int(s.split(" ")[-1])

            if s.split(" ")[0] == "cut":
                card_pos = cut(N, card_pos, num_of_cards)
                continue

            if s.split(" ")[0] == "deal":
                card_pos = deal_with_increment(N, card_pos, num_of_cards)
                continue

            assert(False)
        return card_pos

    print(shuffle_cards(2019, 10007))
    
if ___name___ == '___main___':
    main()
```
And we get the solution ($4684$ in my case).

### Part 2
This is where the puzzle gets interesting.

You get a new deck with $119315717514047$ cards (you can confirm that it's a prime on wolfram alpha :D).
You also need to apply the shuffle process (same one as in the first part) $101741582076661$ *times in a row* (also a prime number). 

Also, we need to calculate _**the number of the card at the position**_ $2020$ (_**the inverse**_ of the function in the first part).

It's clear by this point that we can't rely on algorithms. We need math!

Since we use linear functions to calculate, calculating one shuffle process is pretty quick, but calculating it $101741582076661$ times is a lot.

Since all of the shuffles are linear functions, and the composition of linear functions is a linear function, our whole shuffle process is a linear function.
Suppose that our shuffle process is a function of shape $$f(x) = ax + b.$$
Then, applying that function $m$ times gives us a linear function of shape 
$$
\begin{align}
f^{(m)}(x) & = a^{m}x + b\sum_{n=0}^{m-1}a^{n} \\\\
           & = a^{m}x + b\frac{a^{m} - 1}{a - 1},
\end{align}
$$
in which we used the formula for a finite geometric series.

Since the whole function is calculated modulo $119315717514047$ denoted as $p$, we can write our function like this:
$$
f^{(m)}(x) = (a^{m} \; (mod \; p))x + b((a^{m}\; (mod\; p)) - 1) \cdot ((a - 1)^{-1}\; (mod\; p)).
$$

Great! Now we know how to calculate the new position of a card with a number $n \in \\{0,\ldots,119315717514046\\}$.
Now we just need a way to find the inverse of our linear function (the one we get by applying one shuffle process) and we are done!

Since our function is of shape $f(x) = ax + b$ because all of the shuffles are linear functions, we can calculate the $a$ and $b$ constants by continuously adding and multiplying those constants during our shuffle process. 

Here is the modified code for the shuffle process.

```python
def shuffle_cards(card_pos, num_of_cards):
    a, b = 1, 0
    for s in shuffle:
        if s == "deal into new stack":
            card_pos = deal_new_stack(card_pos, num_of_cards)
            a *= -1
            b = num_of_cards - 1 - b
            continue

        N = int(s.split(" ")[-1])

        if s.split(" ")[0] == "cut":
            card_pos = cut(N, card_pos, num_of_cards)
            b -= N
            continue

        if s.split(" ")[0] == "deal":
            card_pos = deal_with_increment(N, card_pos, num_of_cards)
            a *= N
            b *= N
            continue

        assert(False)
    return card_pos, a % num_of_cards, b % num_of_cards
```

As you can see, we start with constants $0$ and $1$ and we add/multiply them each turn.

Now, if our shuffle process is a function of shape $f(x) = ax + b$, then the inverse if of a shape:
$$
f^{-1}(x) = (a^{-1}\; (mod \; p))x + (-b \cdot a^{-1}\; (mod\; p)),
$$
or in python:
```python
_, c, d = shuffle_cards(start, num_of_cards)
a, b = pow(c, -1, num_of_cards), (-d * pow(c, -1, num_of_cards)) % num_of_cards
```

Now, we are finally ready to solve part 2 of the puzzle. Here is the code for that:
```python
def big_shuffle(start):
    num_of_cards, times = 119315717514047, 101741582076661
    _, c, d = shuffle_cards(start, num_of_cards)
    a, b = pow(c, -1, num_of_cards), (-d * pow(c, -1, num_of_cards)) % num_of_cards
    p1 = pow(a, times, num_of_cards)
    p2 = (p1 - 1) * pow(a - 1, num_of_cards - 2, num_of_cards)
    return (( p1 * start ) + ( b * p2 )) % num_of_cards
```

And our solution is $452290953297$.

## Final remarks

This was a really interesting programming puzzle. Instead of using some clever algorithm, we needed some knowledge of modular arithmetic and linear functions.
The tricky part was recognizing the whole linear function part and taking that approach instead of the algorithm one.
Rest of the puzzle was easy after figuring that out.


This is probably my favorite puzzle of all advent of code puzzles. I'll probably analyze a couple more interesting puzzles from AoC and, since December is close, some of the 2022 puzzles.
If anyone is interested in discussing advent of code puzzles, feel free to contact me at [mislavzanic3@gmail.com](mailto:mislavzanic3@gmail.com).
