---
title: XMonad and NixOS
date: 2022-11-19
tags: 
  - linux
  - nixos
  - xmonad
  - haskell
---
I've been using the XMonad window manager for about two years and NixOS for about six months. 
Despite both having a steep learning curve (I still can't comfortably say that I know exactly how my NixOS config works), XMonad and NixOS create an excellent workflow for me, and I can't imagine using anything else.


## Motivation behind this post
I'll explain the way my XMonad config is built by NixOS.

I wanted to write this because there aren't many "guides" on the compilation of a modular XMonad config using the packages from the master branch of xmonad and xmonad-contrib repo.

I don't know how this post will end up looking.
I'll try to explain (to the best of my ability) how I managed to get my setup working in the hope of helping someone stuck on the same thing.

You can only find two Reddit discussions that discuss this, and those get the compilation working and don't explain much.


## Some context
I was an Arch Linux user before I started using NixOS.
Arch Linux was great for me. 
It was minimal, easy to work with and had all the packages I needed.

It lacked only one thing.
Every time I reinstalled Arch on my desktop or laptop, I needed to install and configure my whole setup (which took some time).

I started writing a script that read a `.yaml` file to find which system packages to install and which config to get.
I found out about NixOS while writing that script and saw it solved all problems I had with Arch (lack of reproducibility and config as code).


### XMonad
For those of you who don't know what XMonad is, it's a _dynamic tyling_ window manager.
Basically, you give it windows, it puts them in a layout.
It's a great window manager.
That is, once you get used to it, it is a pain to configure at first.

Here is an example:

![Image of Windows](/post/images/windows.png "Image of Windows")

Speaking of configuration, there are a couple of ways to configure xmonad.
The simplest way is to install the xmonad package and use `xmonad.hs` file.
[Here](https://nixos.wiki/wiki/XMonad) are the instructions, and [here](https://github.com/mislavzanic/nixos-dotfiles/tree/master/config/xmonad) is my config.

## Setup
For this part, I'll assume that you know what NixOS overlays, modules and flakes are. 
You can read about them [here (overlays)](https://nixos.wiki/wiki/Overlays), [here (modules)](https://nixos.wiki/wiki/NixOS_modules) and [here (flakes)](https://nixos.wiki/wiki/Flakes).
Overlays help you override `nixpkgs` packages.
We'll use them to override `haskellPackages`.

The basic "config building flow" will look like this:

  - we'll get the latest xmonad and xmonad-contrib packages with flakes, 
  - using overlays, we override `pkgs.haskellPackages.xmonad` and `pkgs.haskellPackages.xmonad-contrib`
  - we compile our modular xmonad config using the latest xmonad and xmonad-contrib packages and `cabal2nix` 
  - we add our xmonad package to `pkgs.haskellPackages` using overlays and use it

Also, I'm using a flake skeleton made by hlissner (the creator of [Doom Emacs](https://github.com/doomemacs/doomemacs)), you can find his dotfiles [here](https://github.com/hlissner/dotfiles).


### Adding xmonad and xmonad-contrib as flakes
XMonad and XMonad-contrib repos have their own flakes.
We can add those flakes as inputs to our flake.
```nix
inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";  # for packages on the edge

  xmonad.url = "github:xmonad/xmonad";
  xmonad-contrib.url = "github:xmonad/xmonad-contrib";
};
```
We create our custom xmonad package and override xmonad and xmonad-contrib like this:
```nix
mkPkgs = pkgs: extraOverlays: import pkgs {
  inherit system;
  config.allowUnfree = true;
  overlays = extraOverlays ++ (lib.attrValues self.overlays);
};

pkgs  = mkPkgs nixpkgs [ 
  self.overlay xmonad.overlay xmonad-contrib.overlay (import ./overlays) 
];
```
Finally, our flake:
```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";  # for packages on the edge
    xmonad.url = "github:xmonad/xmonad";
    xmonad-contrib.url = "github:xmonad/xmonad-contrib";
  };
  outputs = inputs @ { self, nixpkgs, nixpkgs-unstable, xmonad, xmonad-contrib, ... }:
  let
    inherit (lib.my) mapModules mapModulesRec mapHosts;
    system = "x86_64-linux";
    mkPkgs = pkgs: extraOverlays: import pkgs {
      inherit system;
      config.allowUnfree = true;
      overlays = extraOverlays ++ (lib.attrValues self.overlays);
    };
    pkgs  = mkPkgs nixpkgs [ 
      self.overlay xmonad.overlay xmonad-contrib.overlay (import ./overlays) 
    ];
    pkgs' = mkPkgs nixpkgs-unstable [];
    lib = nixpkgs.lib.extend
      (self: super: { my = import ./lib { inherit pkgs inputs; lib = self; }; });
  in {
    lib = lib.my;
    overlay = final: prev: {
      unstable = pkgs';
      my = self.packages."${system}";
    };
    overlays = mapModules ./overlays import;
    nixosModules = { dotfiles = import ./.; } // mapModulesRec ./modules import;
    nixosConfigurations = mapHosts ./hosts {};
  };
}
```
We are overriding the default `pkgs.haskellPackages.xmonad` and `pkgs.haskellPackages.xmonad-contrib` with `xmonad.overlay` and `xmonad-contrib.overlay` to get the latest version.


### Cabal2Nix
Our overlay will look like this:
```nix
newPkg: oldPkgs: rec {
  haskellPackages = oldPkgs.haskellPackages.override (old: {
    overrides = oldPkgs.lib.composeExtensions (old.overrides or (_: _: { }))
      (self: super: rec {
        mzanic-xmonad = self.callCabal2nix "mzanic-xmonad" ../config/xmonad { };
      });
  });
}
```
This piece of code added `mzanic-xmonad` package to `haskellPackages`.

Let's see how `cabal2nix ../config/xmonad` output looks like:
```nix
mzanic-xmonad = (
  { mkDerivation
  , base
  , containers
  , hostname
  , lib
  , utf8-string
  , X11 
  , xmonad
  , xmonad-contrib 
  }:
  mkDerivation {
    pname = "mzanic-xmonad";
    version = "0.1.0.0";
    src = ../config/xmonad;
    isLibrary = true;
    isExecutable = true;
    libraryHaskellDepends = [
      base containers hostname utf8-string X11 xmonad xmonad-contrib
    ];
    executableHaskellDepends = [
      base containers hostname utf8-string X11 xmonad xmonad-contrib
    ];
    license = "unknown";
    mainProgram = "mzanic-xmonad";
  }
);
```
Great thing about this is that this, in combination with `xmonad.overlay`, `xmonad-contrib.overlay` and `lib.composeExtensions` compiles our config with the latest and gratest of xmonad and xmonad-contrib.

We are going to use that package later in a module for xmonad:
```nix
{ options, config, pkgs, lib, ... }:
with lib;
with lib.my;
let
  cfg = config.modules.desktop.xmonad;
  configDir = config.dotfiles.configDir;
in {
  options.modules.desktop.xmonad = { enable = mkBoolOpt false; };
  config = mkIf cfg.enable {
    environment.systemPackages = with pkgs; [
      haskellPackages.mzanic-xmonad
    ];

    services = {
      xserver = {
        enable = true;
        displayManager = {
          defaultSession = "none+myxmonad";
          lightdm.enable = true;
          lightdm.greeters.mini = {
            enable = true;
            user = config.user.name;
          };
        };
        windowManager = {
          session = [{
            name = "myxmonad";
            start = ''
              /usr/bin/env mzanic-xmonad &
              waitPID=$!
            '';
          }];
        };
      };
    };
   # ... 
  };
}
```
You can find my XMonad config [here](https://github.com/mislavzanic/nixos-dotfiles/tree/master/config/xmonad) (with cabal files and all files needed for enabling HLS).

## Final remarks
This build process is very complicated.
Even I, who tried to explain the logic behind this, don’t fully understand it, but I hope I helped somebody with this mess of a post.

Nevertheless, I’ll post this and probably rewrite this sometime later.
There aren’t enough posts, guides, and videos on how to build XMonad using source code on any distribution, let alone on NixOS.

This build was the result of many hours, many beers, and a lot of code copy-pasting from GitHub and various posts.
I wanted to get it all written in one place so that others don’t have to suffer my pain.

I also plan on making a post about enabling HLS in an XMonad project on NixOS (this will probably work on any distribution with the Nix package manager).

Until then, you can contact me at [mislavzanic3@gmail.com](mailto:mislavzanic3@gmail.com).
