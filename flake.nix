{
  description = "A simple Go package";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    iosevka.url = "github:mislavzanic/iosevka";
  };


  outputs = { self, nixpkgs, flake-utils, iosevka }:
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
            vendorSha256 = "sha256-+1Hp7LbYaw1T5+fngatqcHPfFC0drHkiLD+NPIrumCo=";
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

          iosevka = pkgs.stdenvNoCC.mkDerivation {
            name = "blog-iosevka";
            buildInputs = with pkgs; [
              python311Packages.brotli
              python311Packages.fonttools
            ];
            dontUnpack = true;
            buildPhase = ''
              mkdir -p out
              ${pkgs.unzip}/bin/unzip ${
                self.inputs.iosevka.packages.${system}.default
              }/ttf.zip
              for ttf in ttf/*.ttf; do
                cp $ttf out
                name=`basename -s .ttf $ttf`
                pyftsubset \
                    $ttf \
                    --output-file=out/"$name".woff2 \
                    --flavor=woff2 \
                    --layout-features=* \
                    --no-hinting \
                    --desubroutinize \
                    --unicodes="U+0000-0170,U+00D7,U+00F7,U+2000-206F,U+2074,U+20AC,U+2122,U+2190-21BB,U+2212,U+2215,U+F8FF,U+FEFF,U+FFFD,U+00E8"
              done

            '';
            installPhase = ''
              mkdir -p $out/static/css/iosevka
              cp out/* $out/static/css/iosevka
            '';
          };

          static = pkgs.stdenv.mkDerivation {
            pname = "blog-static";
            inherit (bin) version;
            inherit src;

            phases = "installPhase";

            installPhase = ''
              mkdir -p $out
              cp -vrf $src/static $out
            '';
          };

          default = pkgs.symlinkJoin {
            name = "blog";
            paths = [ posts static bin iosevka ];
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
