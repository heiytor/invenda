{
  description = "Invenda development flake";

  inputs = {
    nixpkgs = {
      url = "https://flakehub.com/f/NixOS/nixpkgs/0.1.*.tar.gz";
    };
  };

  outputs = { self, nixpkgs }:
    let
      goVersion = 22;
      overlays = [
        (final: prev: rec {
          go = prev."go_1_${toString goVersion}";
          nodejs = prev.nodejs_latest;
          pnpm = prev.nodePackages.pnpm;
        })
      ];
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit overlays system; };
      });
    in
    {
      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            zsh
            go_1_22
            gotools
            golangci-lint
            gopls
            node2nix
            nodejs
            pnpm
            python312
            python312Packages.pip
          ];

          shellHook = ''
            echo "Entering Invneda development environment."
            echo "`go version`"
            echo "node version `node --version`"
            echo "python version `python --version`"
            exec zsh
          '';
        };
      });
    };
}
