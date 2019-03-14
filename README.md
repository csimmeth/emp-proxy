# emp-proxy
Go Proxy for forwarding HLS to Elemental MediaPackage

Elemental Media Package (EMP) accepts HLS input in a specific format. Not all upstream systems support this format, so a proxy is needed to modify the PUT requests into the format EMP expects.

This proxy was built to receive a Zixi HLS HTTP Push, but it should support other systems which output an unauthenticated HTTP HLS Push.

### What it does
The proxy modifies the HLS stream in the following ways:

1. Change the name of the master manifest from index.m3u8 to channel.m3u8.
1. Modify the PUT path to send all manifests and video segments to the root directory, rather than separate directories per rendition. This is accomplished by a simple find and replace all '/' with '_'.
	- e.g. PUT /channel/channel_1080p60.m3u8 -> PUT /channel_channel_1080p60.m3u8
	- e.g. PUT /channel/channel_1080p60/11228888.ts -> PUT /channel_channel_1080p60_11228888.ts
1. Modify the contents of the master manifest to reflect the new filenames of the rendition manifests.
    - e.g. channel/channel_1080p60.m3u8 -> channel_channel_1080p60.m3u8
1. Modify the contents of each rendition manifest to reflect the new filenames of the ts segments. The segments are nested in the directory of the rendition manifests, so this directory name must be prepended the filename of each ts segment. Currently this is hardcoded as 'channel'
	- e.g. channel_1080p60/11228888.ts -> channel_channel_1080p60_11228888.ts
1. Add HTTP Digest Authentication to the HTTP PUT requests.

### Using it

##### To run in a shell

1. Build the binary with 'go build'
2. Add the EMP username, password, and path to wrapper.sh.
3. Run wrapper.sh.

##### To install as a systemd service

1. Build the binary with 'go build'
2. Add the EMP username, password, and path to emp-proxy.service.
3. Run install.sh. The installer will register the service, install the necessary syslog config, and start the service.
