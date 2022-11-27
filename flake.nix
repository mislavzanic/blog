{
  description = "A simple Go package";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    let

      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          src = ./.;
        in rec {
          bin = pkgs.buildGoModule {
            pname = "blog-bin";
            inherit version;
            src = ./.;
            vendorSha256 = "sha256-KlFv+XGgKo7rWa95baZ/aSLBDW00UGkDZgQrjUYDdP0=";
          };

          posts = pkgs.stdenv.mkDerivation {
            pname = "blog-posts";
            inherit (bin) version;
            inherit src;

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -vrf $src/webContent $out
            '';
          };

          static = pkgs.stdenv.mkDerivation {
            pname = "blog-static";
            inherit (bin) version;
            inherit src;

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -vrf $src/js $out
              cp -vrf $src/html $out
              cp -vrf $src/css $out
            '';
          };

          default = pkgs.symlinkJoin {
            name = "blog";
            paths = [ posts static bin ];
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "mislavzanic/blog";
            tag = "dev";

            contents = [ default ];

            config = {
              Cmd = [ "${bin}/bin/blog" ];
              WorkingDir = "${default}";
            };
          };

        });


      apps = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          upload-script = pkgs.writeShellScriptBin "upload-image" ''
            set -eu
            nix build .#docker
            docker load < result
            docker push $USERNAME/blog:dev
          '';
        in {
          upload-script = flake-utils.lib.mkApp { drv = upload-script; };
        }
      );

      defaultPackage = forAllSystems (system: self.packages.${system}.gotsm);
    };
}
