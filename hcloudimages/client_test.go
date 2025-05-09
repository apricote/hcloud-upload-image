package hcloudimages

import (
	"net/url"
	"testing"
)

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}

func TestAssembleCommand(t *testing.T) {
	tests := []struct {
		name    string
		options UploadOptions
		want    string
		wantErr bool
	}{
		{
			name:    "local raw",
			options: UploadOptions{},
			want:    "bash -c 'set -euo pipefail && dd of=/dev/sda bs=4M && sync'",
		},
		{
			name: "remote raw",
			options: UploadOptions{
				ImageURL: mustParseURL("https://example.com/image.xz"),
			},
			want: "bash -c 'set -euo pipefail && wget --no-verbose -O - \"https://example.com/image.xz\" | dd of=/dev/sda bs=4M && sync'",
		},
		{
			name: "local xz",
			options: UploadOptions{
				ImageCompression: CompressionXZ,
			},
			want: "bash -c 'set -euo pipefail && xz -cd | dd of=/dev/sda bs=4M && sync'",
		},
		{
			name: "remote xz",
			options: UploadOptions{
				ImageURL:         mustParseURL("https://example.com/image.xz"),
				ImageCompression: CompressionXZ,
			},
			want: "bash -c 'set -euo pipefail && wget --no-verbose -O - \"https://example.com/image.xz\" | xz -cd | dd of=/dev/sda bs=4M && sync'",
		},
		{
			name: "local bz2",
			options: UploadOptions{
				ImageCompression: CompressionBZ2,
			},
			want: "bash -c 'set -euo pipefail && bzip2 -cd | dd of=/dev/sda bs=4M && sync'",
		},
		{
			name: "remote bz2",
			options: UploadOptions{
				ImageURL:         mustParseURL("https://example.com/image.bz2"),
				ImageCompression: CompressionXZ,
			},
			want: "bash -c 'set -euo pipefail && wget --no-verbose -O - \"https://example.com/image.bz2\" | xz -cd | dd of=/dev/sda bs=4M && sync'",
		},
		{
			name: "local qcow2",
			options: UploadOptions{
				ImageFormat: FormatQCOW2,
			},
			want: "bash -c 'set -euo pipefail && tee image.qcow2 > /dev/null && qemu-img dd -f qcow2 -O raw if=image.qcow2 of=/dev/sda bs=4M && sync'",
		},
		{
			name: "remote qcow2",
			options: UploadOptions{
				ImageURL:    mustParseURL("https://example.com/image.qcow2"),
				ImageFormat: FormatQCOW2,
			},
			want: "bash -c 'set -euo pipefail && wget --no-verbose -O - \"https://example.com/image.qcow2\" | tee image.qcow2 > /dev/null && qemu-img dd -f qcow2 -O raw if=image.qcow2 of=/dev/sda bs=4M && sync'",
		},

		{
			name: "unknown compression",
			options: UploadOptions{
				ImageCompression: "noodle",
			},
			wantErr: true,
		},

		{
			name: "unknown format",
			options: UploadOptions{
				ImageFormat: "poodle",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assembleCommand(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("assembleCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("assembleCommand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
