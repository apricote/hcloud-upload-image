## hcloud-upload-image cleanup

Remove any temporary resources that were left over

### Synopsis

If the upload fails at any point, there might still exist a server or
ssh key in your Hetzner Cloud project. This command cleans up any resources
that match the label "apricote.de/created-by=hcloud-upload-image".

If you want to see a preview of what would be removed, you can use the official hcloud CLI and run:

    $ hcloud server list -l apricote.de/created-by=hcloud-upload-image
    $ hcloud ssh-key list -l apricote.de/created-by=hcloud-upload-image

This command does not handle any parallel executions of hcloud-upload-image
and will remove in-use resources if called at the same time.

```
hcloud-upload-image cleanup [flags]
```

### Options

```
  -h, --help   help for cleanup
```

### Options inherited from parent commands

```
  -v, --verbose count   verbose debug output, can be specified up to 2 times
```

### SEE ALSO

* [hcloud-upload-image](hcloud-upload-image.md)	 - Manage custom OS images on Hetzner Cloud.

