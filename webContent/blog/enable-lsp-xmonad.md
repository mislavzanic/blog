---
title: LSP for XMonad on NixOS
date: 2022-11-28
tags: 
  - linux
  - nixos
  - xmonad
  - haskell
---
Configuring XMonad can be a bit of a pain in the ass. Since it's configured in haskell, LSP can be enabled for XMonad configurations. I'll try to explain how I set up LSP for my config.


## Prerequisites

This setup uses a `shell.nix` file, [direnv](https://direnv.net/), and an editor with direnv integration (emacs in my case, vscode and (n)vim also have it).


## Setup

Firstly, clone the `xmonad` and `xmonad-contrib` repositories in your XMonad config directory.

### shell.nix

We'll use [direnv](https://direnv.net/) for this. If you haven't heard of it, please read about it. Setting up direnv is pretty easy, just create your `shell.nix` or `flake.nix` file and create the `.envrc` with the content:
```sh
use nix # use flake if you are using `flake.nix` file
```

This is the `shell.nix` file:
```nix
{ pkgs ? import <nixpkgs> {} }:
with pkgs;
mkShell {
  buildInputs = [
    haskellPackages.cabal-install
    haskellPackages.haskell-language-server
    haskellPackages.hlint
    haskellPackages.ghcid
    haskellPackages.ormolu
    haskellPackages.implicit-hie
    haskellPackages.X11

    pkg-config

    xorg.libX11
    xorg.libX11.dev

    xorg.libXft
    xorg.libXext
    xorg.libXrandr
    xorg.libXrender
    xorg.libXinerama
    xorg.libXScrnSaver
  ];
}
```
`implicit-hie` is the important program here. We'll use it to generate `hie.yaml` file which is used by `haskell-language-server`.

By this point, our XMonad config directory should look like this:
```sh
❯ ls
drwxr-xr-x   - mzanic 27 Nov 22:55 xmonad
drwxr-xr-x   - mzanic 27 Nov 22:56 xmonad-contrib
.rw-r--r--  22 mzanic 27 Nov 22:56 .envrc
.rw-r--r-- 463 mzanic 27 Nov 22:56 shell.nix
.rw-r--r--   0 mzanic 27 Nov 22:54 xmonad.hs
```
### Cabal

After setting up `shell.nix` and `direnv`, we need to setup the `xmonad-conf.cabal` and `cabal.project` files.
We are doing this so that we can use `gen-hie` to generate the `hie.yaml` file.

Here is how our `xmonad-conf.cabal`: 

```cabal
cabal-version:      2.4
name:               xmonad-conf

executable xmonad-conf
    main-is: xmonad.hs
    build-depends:     X11
                     , base
                     , hostname
                     , containers
                     , utf8-string
                     , xmonad
                     , xmonad-contrib

  default-language:    Haskell2010
  default-extensions:  LambdaCase
  ghc-options:         -Wall
```

and `cabal.project` files look like:

```cabal
packages: . */*.cabal
```

#### gen-hie

We can generate the hie.yaml file by typing in the command `gen-hie > hie.yaml`.
Here what `gen-hie` output looks like:
```yaml
cradle:
  cabal:
    - path: "xmonad/src"
      component: "lib:xmonad"

    - path: "xmonad/./Main.hs"
      component: "xmonad:exe:xmonad"

    - path: "xmonad/tests"
      component: "xmonad:test:properties"

    - path: "xmonad-contrib/./"
      component: "lib:xmonad-contrib"

    - path: "xmonad-contrib/tests"
      component: "xmonad-contrib:test:tests"

    - path: "xmonad-contrib/."
      component: "xmonad-contrib:test:tests"
```

After generating `hie.yaml` file, we need to append this to it:
```yaml
    - path: "./xmonad.hs"
      component: "xmonad-conf:exe:xmonad-conf"
```

The structure of the `component` part is: 
```
<name in the .cabal file>:exe:<name of the executable>
```
in our case, both `xmonad-conf`.


Now, all that is left is for you to open up the editor and enable lsp.
Haskell-language-server should read the config from `hie.yaml` file.

Run `cabal build .` if you got an error, it's a good first step for debugging.

You can test the lsp with this `xmonad.hs` file:
```haskell
import XMonad
import XMonad.Hooks.DynamicLog
import XMonad.Hooks.ManageDocks
import XMonad.Util.Run(spawnPipe)
import XMonad.Util.EZConfig(additionalKeys)
import System.IO


main = do
    xmonad $ defaultConfig
```

Here's how my editor with HLS looks like: 
![LSP image](/post/images/lsp-enabled.png)
![LSP image](/post/images/lsp-enabled2.png)

## Final remarks

I configured the HLS for my xmonad config about 6 months ago, and I tried to do the same for my xmobar config. While doing that, I got stuck and could hardly remember what I did the last time.

So this post is a guide for my future self and for all XMonad users who want LSP and don't know where to start (I was pretty frustrated while configuring HLS for XMonad because of the lack of resources online). Hope this helps!

If you have any questions, contact me at (mailto:mislavzanic3@gmail.com).

