Snap Information
================

#### Setup saned
```
sudo apt install sane-utils

# Enable Saned (required for scanning ; allow connection from the loopback only)
sudo sh -c "echo 127.0.0.1 >> /etc/sane.d/saned.conf"
sudo systemctl enable saned.socket
sudo systemctl start saned.socket


sudo usermod -a -G lp saned
sudo systemctl restart udev.service
```

#### Debugging

- `docker run -v $PWD:$PWD -w $PWD snapcore/snapcraft snapcraft` : https://docs.snapcraft.io/build-snaps/build-on-lxd-docker
- `snap run --shell document-imaging`
- `snap try`


Should look something like:
```
device `net:localhost:genesys:libusb:001:002' is a Canon LiDE 200 flatbed scanner
```

SANE_DEBUG_DLL=3 scanimage -L
