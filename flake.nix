{
  description = "sockaddr-cli";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    docker-tools.url = "github:ZentriaMC/docker-tools";

    docker-tools.inputs.nixpkgs.follows = "nixpkgs";
    docker-tools.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, docker-tools }:
    let
      supportedSystems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];
    in
    flake-utils.lib.eachSystem supportedSystems (system:
      let
        rev = self.rev or "dirty";
        pkgs = nixpkgs.legacyPackages.${system};
      in
      rec {
        packages.sockaddr-cli = pkgs.callPackage ./default.nix {
          inherit rev;
        };

        packages.sockaddr-cli-static = packages.sockaddr-cli.override {
          static = true;
        };

        defaultPackage = packages.sockaddr-cli;

        devShell = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.go_1_18
            pkgs.golangci-lint
            pkgs.gopls
          ];
        };
      });
}
