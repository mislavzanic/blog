---
title: First NixOS package
tags:
  - nixos
date: 2023-07-03
difficulty: 3
---

I created my first NixOS package a couple of months ago.
The package is called [Terraspace](https://terraspace.cloud/). 
I needed it for some work-related stuff, and I didn't find it in the nixpkgs repo, so I took the opportunity to become a nixpkgs contributor.
I had some problems while creating it that I had to debug alone, so I'll recreate the process of packaging it to help anyone stuck on a similar set of problems.
I'll also write a cookbook that I use when creating a nix package.

The great thing about nix and NixOS is the ability to reproducibly package niche programs for your system, and have that config stored in a VCS, and I think that all nix users should be familiar with it without needing to learn how to breathe nix.
It took me more than a year to grasp how powerful this is. I'm writing this post in the hope that it will help a novice nix user with things I had to figure out through troubleshooting.

## The cookbook
1. *Look up a similar package in the nixpkgs repo!* This will save you lots of time, and will give you a great starting point.
2. Search the NixOS wiki. This can give you some useful points, but it can also fail you epically (as you will see later).

Most of the time, the source of the package you want to build will be on GitHub.
These kind of packages are built by fetching a wanted release fron GitHub and then building it (probably by following the instructions found in the README.md).
The exeptions to the above instruction are the `pip`, `gem`, etc. packages. 
You can build those following the above instructions, or by using nix built-in modules for those kinds of packages.

## Building Terraspace
Terraspace is a ruby package, so googling "nixos create ruby package" yields [this](https://nixos.wiki/wiki/Packaging/Ruby), so I decided to follow it.

### Searching the NixOS wiki
The wiki gives us a `shell.nix`:
```nix
with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "env";
  buildInputs = [
    ruby.devEnv
    git
    sqlite
    libpcap
    postgresql
    libxml2
    libxslt
    pkg-config
    bundix
    gnumake
  ];
}
```
and a couble of commands:
```sh
$ nix-shell
$ bundle install      # generates Gemfile.lock
$ bundix              # generates gemset.nix
```
so let's create a `Gemfile`:
```ruby
source "https://rubygems.org"
gem "terraspace", '~> 2.2.7'
```
and do what the wiki tells us to.

Running `bundle install` and `bundix` generates a `Gemfile.lock` and `gemset.nix` files respectively.
```sh
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ ls
Gemfile  Gemfile.lock  gemset.nix  shell.nix
```

Next, we create the `default.nix` file.
Lucky for us, the wiki gives us a pretty good start, and, with some adjustments we get:
```nix
{ stdenv, bundlerEnv, ruby }:
let
  gems = bundlerEnv {
    name = "terraspace-env";
    inherit ruby;
    gemdir  = ./.;
  };
in stdenv.mkDerivation {
  name = "terraspace";
  src = ./.;
  buildInputs = [gems ruby];
  installPhase = ''
    mkdir -p $out/{bin,share/terraspace}
    cp -r * $out/share/terraspace
    bin=$out/bin/terraspace
    cat > $bin <<EOF
#!/bin/sh -e
exec ${gems}/bin/bundle exec ${ruby}/bin/ruby $out/share/terraspace/terraspace "\$@"
EOF
    chmod +x $bin
  '';
}
```

### NixOS wiki failing us
Running the `nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'` we get this:
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'
...
building '/nix/store/zw4qcxiiz4zx35w877mm3pnl93fknrks-nokogiri-1.15.2.gem.drv'...

trying https://rubygems.org/gems/nokogiri-1.15.2.gem
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 4501k  100 4501k    0     0  1522k      0  0:00:02  0:00:02 --:--:-- 1522k
error: hash mismatch in fixed-output derivation '/nix/store/zw4qcxiiz4zx35w877mm3pnl93fknrks-nokogiri-1.15.2.gem.drv':
         specified: sha256-W8ppa5KDrXzpe5wN/fAppiwm6S859ECmV5Xjd9RPEZo=
            got:    sha256-INyAC4++TE9LWxZOaqOrgqNxvLJ+toXBZpYcNN2KItc=
error: 1 dependencies of derivation '/nix/store/i507v7y7qfn5bsvw7l29zq99z7c4mhz9-ruby2.7.7-nokogiri-1.15.2.drv' failed to build
error: 1 dependencies of derivation '/nix/store/fl4db3lr6v8vks68j599nqfhd3cyhhlk-terraspace-env.drv' failed to build
error: 1 dependencies of derivation '/nix/store/lkfx6q1jyv6wzsp24hxdyrcchf25bvcj-terraspace.drv' failed to build
```
This errors tend to show up. 
Nothing big, we just need to replace the `sha` in the `gemset.nix` file.
After replacing it we get:
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'
...
error: builder for '/nix/store/3acagl668bb5bg8dpmam63xkvm782g61-ruby2.7.7-nokogiri-1.15.2.drv' failed with exit code 1;
       last 10 log lines:
       >      from extconf.rb:1034:in `<main>'
       >
       > To see why this extension failed to compile, please check the mkmf.log which can be found here:
       >
       >   /nix/store/ygjqya5igj2d77ik6p6bsfdw2asprr6b-ruby2.7.7-nokogiri-1.15.2/lib/ruby/gems/2.7.0/extensions/x86_64-linux/2.7.0/nokogiri-1.15.2/mkmf.log
       >
       > extconf failed, exit code 1
       >
       > Gem files will remain installed in /nix/store/ygjqya5igj2d77ik6p6bsfdw2asprr6b-ruby2.7.7-nokogiri-1.15.2/lib/ruby/gems/2.7.0/gems/nokogiri-1.15.2 for inspection.
       > Results logged to /nix/store/ygjqya5igj2d77ik6p6bsfdw2asprr6b-ruby2.7.7-nokogiri-1.15.2/lib/ruby/gems/2.7.0/extensions/x86_64-linux/2.7.0/nokogiri-1.15.2/gem_make.out
       For full logs, run 'nix log /nix/store/3acagl668bb5bg8dpmam63xkvm782g61-ruby2.7.7-nokogiri-1.15.2.drv'.
error: 1 dependencies of derivation '/nix/store/d1clnrg5mfzsi79v9kl264cxg0a5qcmg-terraspace-env.drv' failed to build
error: 1 dependencies of derivation '/nix/store/d7snzgwfizvix2y5ql4pygccjm4n5nr4-terraspace.drv' failed to build
```
Seems like we have a problem with nokogiri...
Looking at `Gemfile.lock` we see:
```
nokogiri (1.15.2-x86_64-linux)
  racc (~> 1.4)
```
Googling "nokogiri" gets us to the [rubygems site](https://rubygems.org/gems/nokogiri/versions/1.15.2).
Here we find an error. 
It seems like `bundle install` didn't add `mini_portile2` as a nokogiri dependency.
Also, it has `1.15.2-x86_64-linux` as a version, which seems weird, let's fix this.
This is what the nokogiri section of the `Gemfile.lock` should look like:
```
nokogiri (1.15.2)
  racc (~> 1.4)
  mini_portile2 (~> 2.8.2)
racc (1.7.1)
mini_portile2 (2.8.2)
```
We remove the old `gemset.nix` file and run `bundix` and try to build terraspace again and get:
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'
...
/nix/store/vlpslz0zpqgdn9yp019vva1jgr0rlky8-terraspace
```
Success! Let's run it!
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ result/bin/terraspace -v
Traceback (most recent call last):
/nix/store/yy9sbr2sd4qfn5fdygcqkmibscbcknhq-ruby-2.7.7/bin/ruby: No such file or directory -- /nix/store/vlpslz0zpqgdn9yp019vva1jgr0rlky8-terraspace/share/terraspace/terraspace (LoadError)
```
Ehh... Seems like there's something wrong with our `installPhase` script.
The best way to fix this (in my honest opinion) is to find a simple enough script in the nixpkgs repo and copy it ([this](https://github.com/NixOS/nixpkgs/blob/27343d6e6b710f386aa5df63bdeb16866a782b74/pkgs/tools/misc/pws/default.nix#L2) should do it, it uses the `makeWrapper` function to do the thing that we are going to do next).
But since this is an analysis of a packaging process, we'll try to fix this manually.

### The light at the end of a tunnel
Taking a look at the generated result, we see this:
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace/result/bin]$ cat terraspace
#!/nix/store/96ky1zdkpq871h2dlk198fz0zvklr1dr-bash-5.1-p16/bin/sh -e
exec /nix/store/l9l60bs44jgn59gya59pip2h4rbln66g-terraspace-env/bin/bundle exec /nix/store/yy9sbr2sd4qfn5fdygcqkmibscbcknhq-ruby-2.7.7/bin/ruby /nix/store/vlpslz0zpqgdn9yp019vva1jgr0rlky8-terraspace/share/terraspace/terraspace "$@"
```

There are three different nix-store paths here: `/nix/store/l9l60bs44jgn59gya59pip2h4rbln66g-terraspace-env`, `/nix/store/yy9sbr2sd4qfn5fdygcqkmibscbcknhq-ruby-2.7.7` and `/nix/store/vlpslz0zpqgdn9yp019vva1jgr0rlky8-terraspace`.
The second one is a nix-store path for ruby 2.7.7, and the last one is the one we got running `nix-build`.
Taking a look at `/nix/store/l9l60bs44jgn59gya59pip2h4rbln66g-terraspace-env`, we see:
```
[nix-shell:/nix/store/l9l60bs44jgn59gya59pip2h4rbln66g-terraspace-env]$ ls
bin  lib

[nix-shell:/nix/store/l9l60bs44jgn59gya59pip2h4rbln66g-terraspace-env/bin]$ ls -la | grep terraspace
-r-xr-xr-x 2 root root 1242 Jan  1  1970 terraspace
-r-xr-xr-x 2 root root 1266 Jan  1  1970 terraspace-bundler

[nix-shell:/nix/store/l9l60bs44jgn59gya59pip2h4rbln66g-terraspace-env/bin]$ ./terraspace --version
2.2.7
```

This suggests that our terraspace binary was built in the `bundlerEnv` part of the build script, and we just have to expose it in the `installPhase`:
```sh
mkdir -p $out/bin
bin=$out/bin/terraspace
cat > $bin <<EOF
#!/bin/sh -e
exec ${gems}/bin/terraspace "\$@"
EOF
chmod +x $bin
```

Building it with the new `installPhase` we get:
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'
...
/nix/store/4xrb026k46ql5zxqbfxm753szmclc8qa-terraspace

[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ result/bin/terraspace -v
2.2.7
```
And the thing works this time.

## Final remarks
I hope this post was helpful to someone stuck with building a nix package.
I tried to illustrate some problems I encountered while building terraspace and how I solved them.
The nixpkgs repo is a huge help when writing a package, don't be afraid to look packages up in it.
