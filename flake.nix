{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; # unstable for go 1.22
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          system = system;
        };
        hcloud-upload-image = pkgs.callPackage ./hcloud-upload-image.nix { };
      in
      {
        packages = {
          default = hcloud-upload-image;
        };

        devShells = {
          default = pkgs.mkShell {
            packages = [ hcloud-upload-image ];
          };
        };
      }
    );
}
