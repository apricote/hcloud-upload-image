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
  vendorHash = "sha256-UAofTYCfMmFmCoAFeZPfCnPu9g4j3zknpeeSOHp+rWk=";
  env.GOWORK = "off";
  subPackages = ["."];

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
