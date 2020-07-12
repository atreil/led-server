# Setting up image
When in doubt, follow https://www.raspberrypi.org/documentation/installation/installing-images/linux.md.

1. Run `lsblk -p`.
1. Connect your SD Card to your computer.
1. Run `lsblk -p` again. The newly connected device should appear. The card will have the name `/dev/sdX` any may have partitions `/dev/sdX1`.
1. Unmount any partitions (check the `MOUNTPOINT` column) by running  `umount /dev/sdX1` replacing `1` with whatever partition number you see. 
1. In a terminal window, write the image to the card with the command below, making sure you replace the input file `if=` argument with the path to your `.img` file, and the `/dev/sdX` in the output file `of=` argument with the correct device name. This is very important, as you will lose all the data on the hard drive if you provide the wrong device name. Make sure the device name is the name of the whole SD card as described above, not just a partition. For example: `sdd`, not `sdds1` or `sddp1`; `mmcblk0`, not `mmcblk0p1`.

    ```
    sudo dd bs=4M if=2020-02-13-raspbian-buster-lite.img of=/dev/sdX conv=fsync
    ```

1. This may take some time. You can check the status in a separate terminal by running `sudo kill -USR1 $(pgrep ^dd)`.
1. Run `lsblk -p` again. There should be two partitions, one much bigger than the other. That partition contains the root file system. Mount it by `sudo mount /dev/sdX1 <mount path>`.

# Setting up wireless config
When in doubt, follow https://www.raspberrypi.org/documentation/configuration/wireless/wireless-cli.md

1. Open `<mount path>/etc/wpa_supplicant/wpa_supplicant.conf` and add

    ```
    network={
        ssid="test"
        psk="testpassword"
    }
    ```

    Replacing `test` with your network name and `testpassword` with the network password. Make sure you use tabs instead of spaces. Save the file.

    If your network is hidden, you will need to add `scan_ssid=1`.

    Also, with newer versions of Raspbian, you will need to make sure the following is at the top:

    ```
    ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
    update_config=1
    country=US
    ```

1. Umount the root file system by running `sudo umount /dev/sdX1`.

# Setting up ssh
When in doubt, follow https://www.raspberrypi.org/documentation/remote-access/ssh/.

1. Mount the boot partition - it should be the smaller of the two. Run `sudo mount /dev/sdX1 <mount path>`.

1. Add a file called `ssh` at the root of the boot partition - `touch <mount path>/ssh`.

1. Umount the noot partition by running `sudo umount /dev/sdX1`.

# Connecting to your pi
1. Check your device is connected - `ping raspberrypi.local`

1. I ran into an issue where `rfkill` was soft-blocking WiFi access. You can fix this by removing the SD card from the PI and mounting it on your computer. Run the following:

    ```
    pi@raspberrypi:~ $ ls /var/lib/systemd/rfkill/
    platform-3f300000.mmcnr:wlan  platform-fe300000.mmcnr:wlan  platform-soc:bluetooth
    ```

    The files ending in `:wlan` are files loaded by `rfkill`. Change the contents of one of them to 0, unmount the SD card, and insert it into your PI.

1. Login to your PI

    ```
    ssh pi@raspberrypi.local
    ```
    
    The default password is `raspberry`.

1. Update your password

    ```
    passwd
    ```

# Pre-setup before setting up LED project
1. While `python2.7` is supported, you'll have the best luck with `python3`. Set that as the default version by running
    
    ```
    sudo update-alternatives --install /usr/bin/python python /usr/bin/python3 1
    ```

    Check the version

    ```
    python --version
    ```

1. Install `git`

    ```
    sudo apt-get update
    sudo apt-get install git
    ```

1. Install python tools. You could use `pip` but in my experience `pip` can bungle the science package installations.

    ```
    sudo apt-get install python3-numpy python3-scipy python-pyqtgraph portaudio19-dev python-pyaudio python3-pyaudio libatlas-base-dev
    ```

1. Download the LED project

    ```
    cd ~
    mkdir led-server
    cd led-server
    git clone https://github.com/scottlawsonbc/audio-reactive-led-strip.git
    ```

1. Follow the steps at https://github.com/scottlawsonbc/audio-reactive-led-strip#installation-for-raspberry-pi.

1. In particular:

    Create/edit `/etc/asound.conf`
    ```
    sudo nano /etc/asound.conf
    ```
    Set the file to the following text
    ```
    pcm.!default {
        type hw
        card 1
    }
    ctl.!default {
        type hw
        card 1
    }
    ```

    Next, set the USB device to as the default device by editing `/usr/share/alsa/alsa.conf`
    ```
    sudo nano /usr/share/alsa/alsa.conf:
    ```
    Change
    ```
    defaults.ctl.card 0
    defaults.pcm.card 0
    ```
    To
    ```
    defaults.ctl.card 1
    defaults.pcm.card 1
    ````

## Test the LED strip
1. cd rpi_ws281x/python/examples
1. sudo nano strandtest.py
1. Configure the options at the top of the file. Enable logic inverting if you are using an inverting logic-level converter. Set the correct GPIO pin and number of pixels for the LED strip. You will likely need a logic-level converter to convert the Raspberry Pi's 3.3V logic to the 5V logic used by the ws2812b LED strip.
1. Run example with 'sudo python strandtest.py'

# Setting up remote server
1. You will need golang. Check the latest versions at https://golang.org/dl/.

    ```
    sudo apt-get install golang
    git clone https://github.com/atreil/led-server.git
    ```

