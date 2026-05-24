## hcloud-upload-image write-to-disk

Write the specified disk image to the root disk of the specified server.

### Synopsis

This command writes the specified image to the target servers root disk. Think of
it as a one-off "upload".

#### Image Size

The image size for raw disk images is only limited by the servers root disk.

The image size for qcow2 images is limited to the rescue systems root disk.
This is currently a memory-backed file system with **960 MB** of space. A qcow2
image not be larger than this size, or the process will error. There is a
warning being logged if hcloud-upload-image can detect that your file is larger
than this size.


```
hcloud-upload-image write-to-disk (--image-path=<local-path> | --image-url=<url>) --server <id-or-name> [flags]
```

### Examples

```
  hcloud-upload-image write-to-disk --image-path /home/you/images/custom-linux-image-x86.bz2 --compression bz2 --server my-server
  hcloud-upload-image write-to-disk --image-url https://examples.com/image-arm.raw --server my-arm-server
  hcloud-upload-image write-to-disk --image-url https://examples.com/image-x86.qcow2 --format qcow2 --server my-x86-server
```

### Options

```
      --compression string   Type of compression that was used on the disk image [choices: bz2, xz, zstd]
      --format string        Format of the disk image. [default: raw, choices: qcow2]
  -h, --help                 help for write-to-disk
      --image-path string    Local path to the disk image
      --image-url string     Remote URL of the disk image
      --server string        ID or name of target server
```

### Options inherited from parent commands

```
  -v, --verbose count   verbose debug output, can be specified up to 2 times
```

### SEE ALSO

* [hcloud-upload-image](hcloud-upload-image.md)	 - Manage custom OS images on Hetzner Cloud.

