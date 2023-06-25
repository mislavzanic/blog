---
title: First NixOS package
tags:
  - nixos
date: 2023-06-25
difficulty: 3
---

I created my first NixOS package a couple of months ago.
The package is called [Terraspace](https://terraspace.cloud/). 
I needed it for some work-related stuff, and I didn't find it in the nixpkgs repo, so I took the opportunity to become a nixpkgs contributor.
I had some problems while creating it that I had to debug alone, so I'll recreate the process of packaging it to help anyone stuck on a similar set of problems.

## NixOS Wiki

Terraspace is a ruby package, so googling "nixos create ruby package" yields [this](https://nixos.wiki/wiki/Packaging/Ruby), so I decided to follow it.

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
so lets create a `Gemfile`:
```ruby
source "https://rubygems.org"
gem "terraspace", '~> 2.2.7'
```
and do what the wiki tells us to.

Here's what generated `Gemfile.lock` looks like:
```lock
GEM
  remote: https://rubygems.org/
  specs:
    activesupport (7.0.5)
      concurrent-ruby (~> 1.0, >= 1.0.2)
      i18n (>= 1.6, < 2)
      minitest (>= 5.1)
      tzinfo (~> 2.0)
    aws-eventstream (1.2.0)
    aws-partitions (1.781.0)
    aws-sdk-core (3.175.0)
      aws-eventstream (~> 1, >= 1.0.2)
      aws-partitions (~> 1, >= 1.651.0)
      aws-sigv4 (~> 1.5)
      jmespath (~> 1, >= 1.6.1)
    aws-sdk-kms (1.67.0)
      aws-sdk-core (~> 3, >= 3.174.0)
      aws-sigv4 (~> 1.1)
    aws-sdk-s3 (1.126.0)
      aws-sdk-core (~> 3, >= 3.174.0)
      aws-sdk-kms (~> 1)
      aws-sigv4 (~> 1.4)
    aws-sigv4 (1.5.2)
      aws-eventstream (~> 1, >= 1.0.2)
    cli-format (0.2.2)
      activesupport
      text-table
      zeitwerk
    concurrent-ruby (1.2.2)
    deep_merge (1.2.2)
    diff-lcs (1.5.0)
    dotenv (2.8.1)
    dsl_evaluator (0.3.1)
      activesupport
      memoist
      rainbow
      zeitwerk
    eventmachine (1.2.7)
    eventmachine-tail (0.6.5)
      eventmachine
    graph (2.11.0)
    hcl_parser (0.2.2)
      rhcl
    i18n (1.14.1)
      concurrent-ruby (~> 1.0)
    jmespath (1.6.2)
    memoist (0.16.2)
    minitest (5.18.1)
    nokogiri (1.15.2-x86_64-linux)
      racc (~> 1.4)
    racc (1.7.1)
    rainbow (3.1.1)
    render_me_pretty (0.9.0)
      activesupport
      rainbow
      tilt
    rexml (3.2.5)
    rhcl (0.1.0)
      deep_merge
    rspec (3.12.0)
      rspec-core (~> 3.12.0)
      rspec-expectations (~> 3.12.0)
      rspec-mocks (~> 3.12.0)
    rspec-core (3.12.2)
      rspec-support (~> 3.12.0)
    rspec-expectations (3.12.3)
      diff-lcs (>= 1.2.0, < 2.0)
      rspec-support (~> 3.12.0)
    rspec-mocks (3.12.5)
      diff-lcs (>= 1.2.0, < 2.0)
      rspec-support (~> 3.12.0)
    rspec-support (3.12.0)
    rspec-terraspace (0.3.3)
      activesupport
      memoist
      rainbow
      rspec
      zeitwerk
    rubyzip (2.3.2)
    terraspace (2.2.7)
      activesupport
      bundler
      cli-format
      deep_merge
      dotenv
      dsl_evaluator
      eventmachine-tail
      graph
      hcl_parser
      memoist
      rainbow
      render_me_pretty
      rexml
      rspec-terraspace (>= 0.3.1)
      terraspace-bundler (>= 0.5.0)
      thor
      tty-tree
      zeitwerk
      zip_folder
    terraspace-bundler (0.5.0)
      activesupport
      aws-sdk-s3
      dsl_evaluator
      memoist
      nokogiri
      rainbow
      rubyzip
      thor
      zeitwerk
    text-table (1.2.4)
    thor (1.2.2)
    tilt (2.2.0)
    tty-tree (0.4.0)
    tzinfo (2.0.6)
      concurrent-ruby (~> 1.0)
    zeitwerk (2.6.8)
    zip_folder (0.1.0)
      rubyzip

PLATFORMS
  x86_64-linux

DEPENDENCIES
  terraspace (~> 2.2.7)

BUNDLED WITH
   2.3.26
```
Runing `bundix` generates a `gemset.nix` file.
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
We need to replace the `sha` of nokogiri in the `gemset.nix` file. 
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
This is what `nokogiri` attrset in `gemset.nix` file looks like:
```nix
nokogiri = {
  dependencies = ["racc"];
  groups = ["default"];
  platforms = [];
  source = {
    remotes = ["https://rubygems.org"];
    sha256 = "sha256-INyAC4++TE9LWxZOaqOrgqNxvLJ+toXBZpYcNN2KItc=";
    type = "gem";
  };
  version = "1.15.2";
};
```
Googling "nokogiri" gets us to the [rubygems site](https://rubygems.org/gems/nokogiri/versions/1.15.2) of it.
Here we find an error. 
It seems like `bundle install` didn't add `mini_portile2` as a nokogiri dependency.
Also, it has `1.15.2-x86_64-linux` as a version, which seems weird.
This is what the `nokogiri` section should look like:
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
/nix/store/lrdsc0g56ng65pijbxq83q2b8cdjybag-terraspace
```
Success! Let's run it!
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ result/bin/terraspace -v
Traceback (most recent call last):
/nix/store/yy9sbr2sd4qfn5fdygcqkmibscbcknhq-ruby-2.7.7/bin/ruby: No such file or directory -- /nix/store/lrdsc0g56ng65pijbxq83q2b8cdjybag-terraspace/share/terraspace/terraspace (LoadError)
```
Ehh... It seems something is wrong with the `installPhase`, probably.
We should take a look at ruby packages in the nixpkgs repo.
[This](https://github.com/NixOS/nixpkgs/blob/27343d6e6b710f386aa5df63bdeb16866a782b74/pkgs/tools/misc/pws/default.nix#L2) seems to have a simple `installPhase` that we can use.
Modified `default.nix` looks like this:
```nix
{ stdenv, bundlerEnv, ruby, makeWrapper }:
let
  gems = bundlerEnv {
    name = "terraspace-env";
    inherit ruby;
    gemdir  = ./.;
  };
in stdenv.mkDerivation {
  name = "terraspace";
  src = ./.;
  nativeBuildInputs = [makeWrapper];
  dontUnpack = true;
  installPhase = ''
    mkdir -p $out/bin
    makeWrapper ${gems}/bin/terraspace $out/bin/terraspace
  '';
}
```
Building it we get:
```
[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ nix-build -E '((import <nixpkgs> {}).callPackage (import ./default.nix) { })'
this derivation will be built:
  /nix/store/wdnm4rf8c35v28wyl2zkmcxpfnf225a9-terraspace.drv
building '/nix/store/wdnm4rf8c35v28wyl2zkmcxpfnf225a9-terraspace.drv'...
patching sources
configuring
no configure script, doing nothing
building
no Makefile, doing nothing
installing
post-installation fixup
shrinking RPATHs of ELF executables and libraries in /nix/store/6ig2k1jrf2xygbxz9dmrwcixr4daa9gz-terraspace
strip is /nix/store/dkw46jgi8i0bq64cag95v4ywz6g9bnga-gcc-wrapper-11.3.0/bin/strip
stripping (with command strip and flags -S) in  /nix/store/6ig2k1jrf2xygbxz9dmrwcixr4daa9gz-terraspace/bin
patching script interpreter paths in /nix/store/6ig2k1jrf2xygbxz9dmrwcixr4daa9gz-terraspace
checking for references to /build/ in /nix/store/6ig2k1jrf2xygbxz9dmrwcixr4daa9gz-terraspace...
/nix/store/6ig2k1jrf2xygbxz9dmrwcixr4daa9gz-terraspace

[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ ls
default.nix  Gemfile  Gemfile.lock  gemset.nix  result  shell.nix

[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ result/bin/terraspace
Usage: terraspace COMMAND [args]

The available commands are listed below.
The primary workflow commands are given first, followed by
less common or more advanced commands.

Main Commands:

  terraspace all SUBCOMMAND  # all subcommands
  terraspace build [STACK]   # Build project.
  terraspace bundle          # Bundle with Terrafile.
  terraspace down STACK      # Destroy infrastructure stack.
  terraspace list            # List stacks and modules.
  terraspace new SUBCOMMAND  # new subcommands
  terraspace plan STACK      # Plan stack.
  terraspace seed STACK      # Build starer seed tfvars file.
  terraspace up STACK        # Deploy infrastructure stack.

Other Commands:

  terraspace clean SUBCOMMAND        # clean subcommands
  terraspace completion *PARAMS      # Prints words for auto-completion.
  terraspace completion_script       # Generates a script that can be eval to setup auto-completion.
  terraspace console STACK           # Run console in built terraform project.
  terraspace fmt                     # Run terraform fmt
  terraspace force_unlock            # Calls terrform force-unlock
  terraspace help [COMMAND]          # Describe available commands or one specific command
  terraspace import STACK ADDR ID    # Import existing infrastructure into your Terraform state
  terraspace info STACK              # Get info about stack.
  terraspace init STACK              # Run init in built terraform project.
  terraspace logs [ACTION] [STACK]   # View and tail logs.
  terraspace output STACK            # Run output.
  terraspace providers STACK         # Show providers.
  terraspace refresh STACK           # Run refresh.
  terraspace setup SUBCOMMAND        # setup subcommands
  terraspace show STACK              # Run show.
  terraspace state SUBCOMMAND STACK  # Run state.
  terraspace summary                 # Summarize resources.
  terraspace test                    # Run test.
  terraspace tfc SUBCOMMAND          # tfc subcommands
  terraspace validate STACK          # Validate stack.
  terraspace version                 # Prints version.

For more help on each command, you can use the -h option. Example:

    terraspace up -h

CLI Reference also available at: https://terraspace.cloud/reference/

[nix-shell:~/.local/dev/nix-tinkering/terraspace]$ result/bin/terraspace -v
2.2.7
```
Success! For real this time!

## Final remarks
I hope this post was helpful to someone stuck with building a nix package.
I tried to illustrate some problems I encountered while building terraspace and how I solved them.
The nixpkgs repo is a huge help when writing a package, don't be afraid to look packages up in it.
