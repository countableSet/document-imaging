# Document Imaging

A set of scripts that can be used to acquire and archive scanned documents.

## deb Package Installation

Download the latest deb package from releases and install via commands:

1. `sudo dpkg -i document-imaging.deb`
2. `sudo apt-get install -f` if the dependencies aren't already installed
3. `sudo dpkg -i document-imaging.deb` run install again, if dependencies were missing

## Snap Package Installation

1. `sudo snap install document-imaging --edge`
2. `sudo snap connect document-imaging:raw-usb`

## First-Time Setup

Now that you have document-imaging installed on your system. It's time to step the metadata configuration.

1. Create the metadata file in `$HOME/.config/document-imaging/metadata.json`
2. The content should include the author information in json format; `{ author: "user" }`

## Development Requirements

Install `imagemagick`, `libtiff-tools`, `sane-utils`, `devscripts`, and `dh-make` on Ubuntu.

## Optional Docker Install

0. Compile the code via `make` to build the document-imaging executable file
1. Build the docker image, `docker build -t docimg .` where docimg is the tag name, you can change that to whatever you'd like.
2. To run the images, `docker run --rm -t -i --device=/dev/bus/usb/001/009 -v $HOME/Documents/scan:/scan docimg` where docimg is the tag you picked in setup 1.
  - To find the device path run `lsusb`, where 001 is the bus and 009 is the device
```
$ lsusb
Bus 001 Device 009: ID 04a9:1905 Canon, Inc. CanoScan LiDE 200
```
  - The location on the host machine is included in the volume command in this case: `$HOME/Documents/scan`
3. You can add bash to the end of the run command to not stop the container after every run.
