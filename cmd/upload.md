This command implements a fake "upload", by going through a real server and
snapshots. This does cost a bit of money for the server.

#### Image Size

The image size for raw disk images is only limited by the servers root disk.

The image size for qcow2 images is limited to the rescue systems root disk.
This is currently a memory-backed file system with **960 MB** of space. A qcow2
image not be larger than this size, or the process will error. There is a
warning being logged if hcloud-upload-image can detect that your file is larger
than this size.
