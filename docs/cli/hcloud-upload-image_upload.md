## hcloud-upload-image upload

Upload the specified disk image into your Hetzner Cloud project.

### Synopsis

This command implements a fake "upload", by going through a real server and
snapshots. This does cost a bit of money for the server.

#### Image Size

The image size for raw disk images is only limited by the servers root disk.

The image size for qcow2 images is limited to the rescue systems root disk.
This is currently a memory-backed file system with **960 MB** of space. A qcow2
image not be larger than this size, or the process will error. There is a
warning being logged if hcloud-upload-image can detect that your file is larger
than this size.


```
hcloud-upload-image upload (--image-path=<local-path> | --image-url=<url>) --architecture=<x86|arm> [flags]
```

### Examples

```
  hcloud-upload-image upload --image-path /home/you/images/custom-linux-image-x86.bz2 --architecture x86 --compression bz2 --description "My super duper custom linux"
  hcloud-upload-image upload --image-url https://examples.com/image-arm.raw --architecture arm --labels foo=bar,version=latest
  hcloud-upload-image upload --image-url https://examples.com/image-x86.qcow2 --architecture x86 --format qcow2
```

### Options

```
      --architecture string     CPU architecture of the disk image [choices: x86, arm]
      --compression string      Type of compression that was used on the disk image [choices: bz2, xz]
      --description string      Description for the resulting image
      --format string           Format of the image. [choices: qcow2]
  -h, --help                    help for upload
      --image-path string       Local path to the disk image that should be uploaded
      --image-url string        Remote URL of the disk image that should be uploaded
      --labels stringToString   Labels for the resulting image (default [])
      --server-type string      Explicitly use this server type to generate the image. Mutually exclusive with --architecture.
```

### Options inherited from parent commands

```
  -v, --verbose count   verbose debug output, can be specified up to 2 times
```

### SEE ALSO

* [hcloud-upload-image](hcloud-upload-image.md)	 - Manage custom OS images on Hetzner Cloud.

