Snap Information
================

#### Setup saned for network scanning
```
sudo apt install sane-utils
# Enable Saned (required for scanning ; allow connection from the loopback only)
sudo sh -c "echo 127.0.0.1 >> /etc/sane.d/saned.conf"
sudo systemctl enable saned.socket
sudo systemctl start saned.socket
sudo systemctl stop saned.socket
sudo systemctl disable saned.socket
# If problems?
sudo usermod -a -G lp saned
sudo systemctl restart udev.service
```

#### Debugging

- Docker / lxd: https://docs.snapcraft.io/build-snaps/build-on-lxd-docker
- `docker run -v $PWD:$PWD -w $PWD snapcore/snapcraft snapcraft` :
- `snap run --shell document-imaging`
- `SANE_DEBUG_DLL=4 scanimage -L`


#### Random

Network device name should look something like:
```
device `net:localhost:genesys:libusb:001:002' is a Canon LiDE 200 flatbed scanner
```

SANE_CONFIG_DIR
See: https://gitlab.com/sane-project/backends/blob/RELEASE_1_0_25/sanei/sanei_config.c#L92
