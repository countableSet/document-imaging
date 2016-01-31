# Document Imaging

A set of scripts that can be used to acquire and archive scanned documents.

## Requirements

Install `imagemagick` on Ubuntu.

## Setup
1. Run `scanimage -L`
2. Note the device name similar to 'genesys:libusb:001:024'
3. Modify the script scan-image.sh to include the device
4. Run ./scan-image
