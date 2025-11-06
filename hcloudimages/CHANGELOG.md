# Changelog

## [1.2.0](https://github.com/apricote/hcloud-upload-image/compare/hcloudimages/v1.1.0...hcloudimages/v1.2.0) (2025-11-06)


### Features

* change minimum required Go version to 1.24 ([#130](https://github.com/apricote/hcloud-upload-image/issues/130)) ([5eba2d5](https://github.com/apricote/hcloud-upload-image/commit/5eba2d52fe3aafb4fd0d93403548f4c32bc2b5ac))
* support zstd compression ([#125](https://github.com/apricote/hcloud-upload-image/issues/125)) ([37ebbce](https://github.com/apricote/hcloud-upload-image/commit/37ebbce5179997ac216af274055fc34c777b01e6)), closes [#122](https://github.com/apricote/hcloud-upload-image/issues/122)
* update default x86 server type to cx23 ([#129](https://github.com/apricote/hcloud-upload-image/issues/129)) ([a205619](https://github.com/apricote/hcloud-upload-image/commit/a20561944d0ba9485a6e10e99df15c56a688541d))

## [1.1.0](https://github.com/apricote/hcloud-upload-image/compare/hcloudimages/v1.0.1...hcloudimages/v1.1.0) (2025-05-10)


### Features

* smaller snapshots by zeroing disk first ([#101](https://github.com/apricote/hcloud-upload-image/issues/101)) ([fdfb284](https://github.com/apricote/hcloud-upload-image/commit/fdfb284533d3154806b0936c08015fd5cc64b0fb)), closes [#96](https://github.com/apricote/hcloud-upload-image/issues/96)


### Bug Fixes

* upload from local image generates broken command ([#98](https://github.com/apricote/hcloud-upload-image/issues/98)) ([420dcf9](https://github.com/apricote/hcloud-upload-image/commit/420dcf94c965ee470602db6c9c23c777fda91222)), closes [#97](https://github.com/apricote/hcloud-upload-image/issues/97)

## [1.0.1](https://github.com/apricote/hcloud-upload-image/compare/hcloudimages/v1.0.0...hcloudimages/v1.0.1) (2025-05-09)


### Bug Fixes

* timeout while waiting for SSH to become available ([#92](https://github.com/apricote/hcloud-upload-image/issues/92)) ([e490b9a](https://github.com/apricote/hcloud-upload-image/commit/e490b9a7f394e268fa1946ca51aa998c78c3d46a))

## [1.0.0](https://github.com/apricote/hcloud-upload-image/compare/hcloudimages/v0.3.1...hcloudimages/v1.0.0) (2025-05-04)


### Features

* upload qcow2 images ([#69](https://github.com/apricote/hcloud-upload-image/issues/69)) ([ac3e9dd](https://github.com/apricote/hcloud-upload-image/commit/ac3e9dd7ecd86d1538b6401c3073c7c078c40847))

## [0.3.1](https://github.com/apricote/hcloud-upload-image/compare/hcloudimages/v0.3.0...hcloudimages/v0.3.1) (2024-12-07)


### Bug Fixes

* **cli:** local install fails because of go.mod replace ([#47](https://github.com/apricote/hcloud-upload-image/issues/47)) ([66dc5f7](https://github.com/apricote/hcloud-upload-image/commit/66dc5f70b604ed3ee964576d74f94bdcea710c95))

## [0.3.0](https://github.com/apricote/hcloud-upload-image/compare/hcloudimages/v0.2.0...hcloudimages/v0.3.0) (2024-06-23)


### Features

* set server type explicitly ([#36](https://github.com/apricote/hcloud-upload-image/issues/36)) ([42eeb00](https://github.com/apricote/hcloud-upload-image/commit/42eeb00a0784e13a00a52cf15a8659b497d78d72)), closes [#30](https://github.com/apricote/hcloud-upload-image/issues/30)
* update default x86 server type to cx22 ([#38](https://github.com/apricote/hcloud-upload-image/issues/38)) ([ebe08b3](https://github.com/apricote/hcloud-upload-image/commit/ebe08b345c8f31df73087b091fa39f5fdc195156))


### Bug Fixes

* error early when the image write fails ([#34](https://github.com/apricote/hcloud-upload-image/issues/34)) ([256989f](https://github.com/apricote/hcloud-upload-image/commit/256989f4a37e7b124c0684aab0f34cf5e09559be)), closes [#33](https://github.com/apricote/hcloud-upload-image/issues/33)
