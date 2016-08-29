# Document Imaging

A set of scripts that can be used to acquire and archive scanned documents.

## deb Package Installation

Download the latest deb package from releases and install via commands:

1. `sudo dpkg -i document-imaging.deb`
2. `sudo apt-get install -f` if the dependencies aren't already installed
3. `sudo dpkg -i document-imaging.deb` run install again, if dependencies were missing

## Requirements

Install `imagemagick', 'libtiff-tools`, and `sane-utils` on Ubuntu.

## Setup for Legacy Bash Script
1. Run `scanimage -L`
2. Note the device name similar to 'genesys:libusb:001:024'
3. Modify the script scan-image.sh to include the device
4. Run ./scan-image
