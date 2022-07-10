{ lib, buildGo118Module, rev ? "dirty", static ? false }:

buildGo118Module rec {
  pname = "sockaddr-cli";
  version = rev;

  src = lib.cleanSource ./.;

  ldflags = [
    "-s"
    "-w"
    "-X github.com/ZentriaMC/sockaddr-cli/internal/core.Version=${version}"
  ] ++ lib.optionals static [
    "-extldflags=-static"
  ];

  CGO_ENABLED = if static then 0 else 1;

  checkPhase = ''
    runHook preCheck

    for pkg in $(find . -type f -name "*_test.go" -print0 | xargs -0 -r dirname -- | sort -z -u | tr '\0' '\n'); do
      buildGoDir test "$pkg"
    done

    runHook postCheck
  '';

  doCheck = true;

  vendorSha256 = "sha256-vKnK4pF+lUSVnVY4hEnbgCPeFLmASW2UQchELr4d/Xc=";
  subPackages = [ "cmd/sockaddr-cli" ];
}
