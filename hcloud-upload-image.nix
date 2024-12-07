{ config, pkgs, fetchFromGitHub, ... }:

pkgs.buildGo122Module rec {
  pname = "hcloud-upload-image";
  version = "v0.2.1"; # x-release-please-version
  #src = pkgs.fetchFromGitHub {
  #  owner = "apricote";
  #  repo = "hcloud-upload-image";
  #  rev = version;
  #  sha256 = "KYhmMo/GDiLOWuEtdiY/KDwh1MO3crAFN2SGiL8FFIU=";
  #};
  src = ./.;

  vendorHash = "";
  overrideModAttrs = _: {
    #outputHash = "sha256-/FluEFTvrX07L4kGPH7jPrcRhwdrCz1p4OVErNegNLw=";
    #outputHash = null;
    #outputHashAlgo = null;
    #outputHashMode = null;
  };

  nativeBuildInputs = [
    pkgs.installShellFiles
  ];

  preBuild = ''
    export GOWORK=off
  '';

  subPackages = [ "." ]; # We only need to CLI, fails with the "hcloudimages" module otherwise
  CGO_ENABLED = "0";

  ldflags = [
    "-X github.com/apricote/hcloud-upload-image/internal/version.version=${version}"
    "-X github.com/apricote/hcloud-upload-image/internal/version.versionPrerelease="
  ];

  postInstall = ''
    installShellCompletion --cmd hcloud-upload-image \
      --bash <($out/bin/hcloud-upload-image completion bash) \
      --fish <($out/bin/hcloud-upload-image completion fish) \
      --zsh <($out/bin/hcloud-upload-image completion zsh)
  '';
}
