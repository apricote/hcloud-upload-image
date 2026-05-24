{
  lib,
  buildGoModule,
}:
buildGoModule rec {
  pname = "hcloud-upload-image";
  version =
    builtins.head (builtins.match ".*version = \"([0-9.]+)\".*"
      (builtins.readFile ./internal/version/version.go));

  src = ./.;
  vendorHash = "sha256-YgiSTnCg+7RuSxp0NDNn7QamV7HgdQY+DrbEIv68eJw=";
  env.GOWORK = "off";
  subPackages = ["."];
  goSum = ./go.sum; # make sure to rebuild

  postPatch = ''
    echo 'replace github.com/apricote/hcloud-upload-image/hcloudimages => ./hcloudimages' >> go.mod
  '';


  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ];

  meta = {
    description = "Quickly upload any raw disk images into your Hetzner Cloud projects";
    homepage = "https://github.com/apricote/hcloud-upload-image";
    license = lib.licenses.mit;
    maintainers = [];
    mainProgram = "hcloud-upload-image";
  };
}
