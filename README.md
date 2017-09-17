# Document Imaging

A set of scripts that can be used to acquire and archive scanned documents.

## deb Package Installation

Download the latest deb package from releases and install via commands:

1. `sudo dpkg -i document-imaging.deb`
2. `sudo apt-get install -f` if the dependencies aren't already installed
3. `sudo dpkg -i document-imaging.deb` run install again, if dependencies were missing

## First-Time Setup

Now that you have document-imaging installed on your system. It's time to step the metadata configuration.

1. Create the metadata file in `$HOME/.config/document-imaging/metadata.json`
2. The content should include the author information in json format; `{ author: "user" }`

## Development Requirements

Install `imagemagick`, `libtiff-tools`, `sane-utils`, `devscripts`, and `dh-make` on Ubuntu.
