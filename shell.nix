{ pkgs ? import <nixpkgs> {} }:
with pkgs;
mkShell {
  buildInputs = [
    go_1_18
    gopls
  ];
  shellHook = ''
    export XDG_DATA_HOME="/home/mzanic/.local/share"
    export GOPATH="$XDG_DATA_HOME/go"
    export GO111MODULE=on
    export PATH=$GOPATH/bin:$PATH
  '';
}
