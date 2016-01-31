#!/bin/bash
scanimage --device-name 'genesys:libusb:002:006' --mode Color --resolution 300 -x 215.9 -y 279.4 --format=tiff -p > `date +"%Y%m%d%H%M%S"`.tiff
