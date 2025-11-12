{
  description = "hcloud-upload-image - Quickly upload any raw disk images into your Hetzner Cloud projects";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ { flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      perSystem = { pkgs, ... }:
        let
          pkg = pkgs.callPackage ./default.nix { };
          app = {
            type = "app";
            program = "${pkg}/bin/hcloud-upload-image";
          };
        in
        {
          packages.default = pkg;
          packages.hcloud-upload-image = pkg;
          apps.default = app;
          apps.hcloud-upload-image = app;
          devShells.default = pkgs.callPackage ./shell.nix { };
          formatter = pkgs.nixpkgs-fmt;
        };
    };
}
