---
title: TEST
tags:
  - nixos
date: 2023-07-03
difficulty: 3
---

Ehh... Seems like there's^1^ something
<input id="_1" type="checkbox"><label class="drop" for="_1">Collapse 1 </label> 

<div>Content 1</div>

wrong with our `installPhase` script.
The best way to fix this (in my opinion) is to find a simple enough script in the nixpkgs repo and copy it ([this](https://github.com/NixOS/nixpkgs/blob/27343d6e6b710f386aa5df63bdeb16866a782b74/pkgs/tools/misc/pws/default.nix#L2) should do it, it uses the `makeWrapper`function to do the thing that we are about to).
lsdkjf;alkdf;alksdfj;klasjd;fljas;dlfkja ;lsdkjf;laskdjf;lasjdf; laskjdf;lajdflask jdfasldjfaldjfa;lsdfja;lsfj.
alkdfjalsjdf;alskjf;alksjf ;alskjdf;aslkjdf;aslkjfa;slfkj;alsfja;slkfja;lsjfa;lsdfj.
lsdkjf;alkdf;alksdfj;klasjd;fljas;dlfkja ;lsdkjf;laskdjf;lasjdf; laskjdf;lajdflask jdfasldjfaldjfa;lsdfja;lsfj.
alkdfjalsjdf;alskjf;alksjf ;alskjdf;aslkjdf;aslkjfa;slfkj;alsfja;slkfja;lsjfa;lsdfj.
lsdkjf;alkdf;alksdfj;klasjd;fljas;dlfkja ;lsdkjf;laskdjf;lasjdf; laskjdf;lajdflask jdfasldjfaldjfa;lsdfja;lsfj.
alkdfjalsjdf;alskjf;alksjf ;alskjdf;aslkjdf;aslkjfa;slfkj;alsfja;slkfja;lsjfa;lsdfj.
But since this is an analysis of a packaging process, we'll try to fix this manually.

<sidenote id="sn:1">
<sup>1 </sup>laskdjflaskdjflkasdflkasjdlkajflkjfk:
1. lkdsfj
2. alkdfj
</sidenote>

<sidenote id="sn:3">
laskdjflaskdjflkasdflkasjdlkajflkjfk:
1. lkdsfj
2. alkdfj
</sidenote>

<sidenote id="sn:4">
nekaj drugo
</sidenote>
