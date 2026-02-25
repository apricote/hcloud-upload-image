# Changelog

## [1.4.0](https://github.com/apricote/hcloud-upload-image/compare/v1.3.0...v1.4.0) (2026-02-25)


### Features

* Add Nix support ([#109](https://github.com/apricote/hcloud-upload-image/issues/109)) ([756f051](https://github.com/apricote/hcloud-upload-image/commit/756f0515707c58bd228349928a6c88f58f1fe56a))

## [1.3.0](https://github.com/apricote/hcloud-upload-image/compare/v1.2.0...v1.3.0) (2025-12-22)


### Features

* add --location flag to specify datacenter region ([#141](https://github.com/apricote/hcloud-upload-image/issues/141)) ([fcbc14a](https://github.com/apricote/hcloud-upload-image/commit/fcbc14aab6d495d2c67d653f9ea1ff56a39a8c2f)), closes [#142](https://github.com/apricote/hcloud-upload-image/issues/142)

## [1.2.0](https://github.com/apricote/hcloud-upload-image/compare/v1.1.0...v1.2.0) (2025-11-06)


### Features

* change minimum required Go version to 1.24 ([#130](https://github.com/apricote/hcloud-upload-image/issues/130)) ([5eba2d5](https://github.com/apricote/hcloud-upload-image/commit/5eba2d52fe3aafb4fd0d93403548f4c32bc2b5ac))
* support zstd compression ([#125](https://github.com/apricote/hcloud-upload-image/issues/125)) ([37ebbce](https://github.com/apricote/hcloud-upload-image/commit/37ebbce5179997ac216af274055fc34c777b01e6)), closes [#122](https://github.com/apricote/hcloud-upload-image/issues/122)
* update default x86 server type to cx23 ([#129](https://github.com/apricote/hcloud-upload-image/issues/129)) ([a205619](https://github.com/apricote/hcloud-upload-image/commit/a20561944d0ba9485a6e10e99df15c56a688541d))

## [1.1.0](https://github.com/apricote/hcloud-upload-image/compare/v1.0.1...v1.1.0) (2025-05-10)


### Features

* smaller snapshots by zeroing disk first ([#101](https://github.com/apricote/hcloud-upload-image/issues/101)) ([fdfb284](https://github.com/apricote/hcloud-upload-image/commit/fdfb284533d3154806b0936c08015fd5cc64b0fb)), closes [#96](https://github.com/apricote/hcloud-upload-image/issues/96)


### Bug Fixes

* upload from local image generates broken command ([#98](https://github.com/apricote/hcloud-upload-image/issues/98)) ([420dcf9](https://github.com/apricote/hcloud-upload-image/commit/420dcf94c965ee470602db6c9c23c777fda91222)), closes [#97](https://github.com/apricote/hcloud-upload-image/issues/97)

## [1.0.1](https://github.com/apricote/hcloud-upload-image/compare/v1.0.0...v1.0.1) (2025-05-09)


### Bug Fixes

* timeout while waiting for SSH to become available ([#92](https://github.com/apricote/hcloud-upload-image/issues/92)) ([e490b9a](https://github.com/apricote/hcloud-upload-image/commit/e490b9a7f394e268fa1946ca51aa998c78c3d46a))

## [1.0.0](https://github.com/apricote/hcloud-upload-image/compare/v0.3.1...v1.0.0) (2025-05-04)


### Features

* **deps:** require Go 1.23 ([#70](https://github.com/apricote/hcloud-upload-image/issues/70)) ([f3fcb62](https://github.com/apricote/hcloud-upload-image/commit/f3fcb623fc00095ab3806fa41dbcb7083c13c5df))
* docs website ([#80](https://github.com/apricote/hcloud-upload-image/issues/80)) ([d144b85](https://github.com/apricote/hcloud-upload-image/commit/d144b85e3dfd933e8fbb09a0e6f5acacb4d05bea))
* publish container image ([#82](https://github.com/apricote/hcloud-upload-image/issues/82)) ([91df729](https://github.com/apricote/hcloud-upload-image/commit/91df729f1cfd636355fc8338f47aefa4ab8b3b84))
* upload qcow2 images ([#69](https://github.com/apricote/hcloud-upload-image/issues/69)) ([ac3e9dd](https://github.com/apricote/hcloud-upload-image/commit/ac3e9dd7ecd86d1538b6401c3073c7c078c40847))

## [0.3.1](https://github.com/apricote/hcloud-upload-image/compare/v0.3.0...v0.3.1) (2024-12-07)


### Bug Fixes

* **cli:** local install fails because of go.mod replace ([#47](https://github.com/apricote/hcloud-upload-image/issues/47)) ([66dc5f7](https://github.com/apricote/hcloud-upload-image/commit/66dc5f70b604ed3ee964576d74f94bdcea710c95))

## [0.3.0](https://github.com/apricote/hcloud-upload-image/compare/v0.2.1...v0.3.0) (2024-06-23)


### Features

* set server type explicitly ([#36](https://github.com/apricote/hcloud-upload-image/issues/36)) ([42eeb00](https://github.com/apricote/hcloud-upload-image/commit/42eeb00a0784e13a00a52cf15a8659b497d78d72)), closes [#30](https://github.com/apricote/hcloud-upload-image/issues/30)
* update default x86 server type to cx22 ([#38](https://github.com/apricote/hcloud-upload-image/issues/38)) ([ebe08b3](https://github.com/apricote/hcloud-upload-image/commit/ebe08b345c8f31df73087b091fa39f5fdc195156))


### Bug Fixes

* error early when the image write fails ([#34](https://github.com/apricote/hcloud-upload-image/issues/34)) ([256989f](https://github.com/apricote/hcloud-upload-image/commit/256989f4a37e7b124c0684aab0f34cf5e09559be)), closes [#33](https://github.com/apricote/hcloud-upload-image/issues/33)

## [0.2.1](https://github.com/apricote/hcloud-upload-image/compare/v0.2.0...v0.2.1) (2024-05-10)


### Bug Fixes

* **cli:** completion requires HCLOUD_TOKEN ([#19](https://github.com/apricote/hcloud-upload-image/issues/19)) ([bb2ca48](https://github.com/apricote/hcloud-upload-image/commit/bb2ca482000f5c780545edb9a03aa9f6bf93d906))

## [0.2.0](https://github.com/apricote/hcloud-upload-image/compare/v0.1.1...v0.2.0) (2024-05-09)


### Features

* packaging for deb, rpm, apk, aur ([#17](https://github.com/apricote/hcloud-upload-image/issues/17)) ([139761c](https://github.com/apricote/hcloud-upload-image/commit/139761cc28050b00bca22573d765f2b94af89bac))
* upload local disk images ([#15](https://github.com/apricote/hcloud-upload-image/issues/15)) ([fcea3e3](https://github.com/apricote/hcloud-upload-image/commit/fcea3e3c6e5ba7383aa69838401903e3f54f910c))
* upload xz compressed images ([#16](https://github.com/apricote/hcloud-upload-image/issues/16)) ([1c943e4](https://github.com/apricote/hcloud-upload-image/commit/1c943e4480ba2042fc3feabf363ec88eb2efbaee))


### Bug Fixes

* update user-agent in CLI ([#5](https://github.com/apricote/hcloud-upload-image/issues/5)) ([b17857c](https://github.com/apricote/hcloud-upload-image/commit/b17857c1fefc0b09da2ed2711b20ba76930dd365))

## [0.1.1](https://github.com/apricote/hcloud-upload-image/compare/v0.1.0...v0.1.1) (2024-05-04)


### Bug Fixes

* CLI does not produce release binaries ([#3](https://github.com/apricote/hcloud-upload-image/issues/3)) ([f373d4c](https://github.com/apricote/hcloud-upload-image/commit/f373d4c2baca9ccc892e6b6abff6dd217f2fdbeb))

## [0.1.0](https://github.com/apricote/hcloud-upload-image/compare/v0.0.1...v0.1.0) (2024-05-04)


### Features

* **cli:** docs grouping and version ([847b696](https://github.com/apricote/hcloud-upload-image/commit/847b696c74ce67c2f18aaa69af60f6c0c5b736c4))
* **cli:** hide redundant log attributes ([9e65452](https://github.com/apricote/hcloud-upload-image/commit/9e654521ae12debf40f181dfe291ad4ded0f7524))
* **cli:** upload command ([b6ae95f](https://github.com/apricote/hcloud-upload-image/commit/b6ae95f55ba134f5ef124d377ed3ad0a556b8cf4))
* documentation and cleanup command ([c9ab40b](https://github.com/apricote/hcloud-upload-image/commit/c9ab40b539bc51ea2611bb0b58ab8aef4ec06eea))
* initial library code ([4f57df5](https://github.com/apricote/hcloud-upload-image/commit/4f57df5b66ed1391155792758737b8f54b7ef2ab))
* log output ([904e5e0](https://github.com/apricote/hcloud-upload-image/commit/904e5e0bed6ba87e0f4063c27a0678a9c85b7371))
