import abc
import argparse
import collections
import os
import re
import subprocess
import sys
import getpass

from glob import glob

EDITOR = os.environ.get('EDITOR','vi')

Args = collections.namedtuple('Args', 'partition mount_path')
Partitions = collections.namedtuple('Partitions', 'boot root')

############ HELPER METHODS ############
def wrap_cmd(args):
    print('running \'%s\'' % args)
    ret = subprocess.call(args)
    print('cmd returned %s' % ret)
    if ret != 0:
        raise Exception('\'%s\' returned non-zero exit code: %s' % (args, ret))
    return ret

def confirm(msg):
    valid_confirmed = False
    confirmed = ''
    while not valid_confirmed:
        sys.stdout.write('%s? (yY/nN): ' % msg)
        sys.stdout.flush()
        confirmed = read_line()
        valid_confirmed = re.search('^(n|N|y|Y)$', confirmed) is not None
    return confirmed == 'y' or confirmed == 'Y'

def check_file(path):
    path = resolve_path(path)
    if not os.path.exists(path) or not os.path.isfile(path):
        raise Exception('\'%s\' doesn\'t exist or is a file' % path)
    return path

def check_dir(path):
    path = resolve_path(path)
    if not os.path.exists(path) or not os.path.isdir(path):
        raise Exception('\'%s\' doesn\'t exist or is a file' % path)
    return path

def resolve_path(path):
    if not os.path.isabs(path):
        path = os.path.join(os.getcwd(), path)
    path = os.path.normpath(path)
    return path

def read_line():
    return sys.stdin.readline().rstrip('\n')

def is_mounted(dev):
    args = ['cat', '/etc/mtab']
    print('running \'%s\'' % args)
    output = subprocess.check_output(args)
    print('cmd returned \'%s\'' % output)
    return output.decode("utf-8").find(dev) != -1

def append_file(file, text):
    with open(file, 'a') as f:
        f.write(text)

def check_mount(dev: str, mnt: str):
    """Runs checks and then mounts a device. Will create the directory."""
    mnt = resolve_path(mnt)
    if is_mounted(dev):
        raise Exception('\'%s\' is already mounted at %s' % (dev, mnt))
    if os.path.isfile(mnt):
        raise Exception('%s should be a directory. found file' % mnt)
    if not os.path.exists(mnt):
        os.makedirs(mnt, exist_ok=True)
    wrap_cmd(['mount', dev, mnt])

def resolve_dev_partitions(dev_pfx: str) -> Partitions:
    """Resolves the root and boot partitions from the top-level device."""
    args = ['lsblk', dev_pfx, '-rbp', '-o', 'NAME,SIZE']
    print('running \'%s\'' % args)
    output = subprocess.check_output(args)
    print('got output %s', output)
    dev_lines = [line.split(' ') for line in output.decode('utf-8').split('\n') if line.startswith(dev_pfx)]
    partition_lines = [dev_line for dev_line in dev_lines if dev_line[0] != dev_pfx]
    if len(partition_lines) != 2:
        raise Exception('only expected to find 2 devices. found %s' % partition_lines)
    root = partition_lines[0]
    boot = partition_lines[1]
    if int(boot[1]) > int(root[1]):
        tmp = boot
        boot = root
        root = tmp
    return Partitions(boot=boot[0], root=root[0])

def unmount_device(dev):
    """Unmount the devivce."""
    to_unmount = [part for part in glob('%s?*' % dev) if is_mounted(part)]

    if len(to_unmount) != 0:
        if not confirm('Partitions \'%s\' needs to be unmounted first. '
                    'Proceed' % dev):
                raise Exception('User canceled unmount operation')
        wrap_cmd(['umount'] + to_unmount)
        print('partions %s are unmounted' % glob('%s?*' % dev))

########################################

############# Commands #################

def write_image(dev: str, img_path: str) -> None :
    """Format the SD card device with a specific image.
    
    Args:
    dev - The top-level name of the device of the SD card. Must be unmounted.
    img_path - Path to the image file to format the SD card with.
    """
    check_file(img_path)
    wrap_cmd(['dd', 'bs=4M', 'if=%s' % img_path, 'of=%s' % dev, 'conv=fsync',
              'status=progress'])
    wrap_cmd(['sync'])

def setup_ssh(mount_path):
    wrap_cmd(['touch', os.path.join(mount_path, 'ssh')])

def setup_wireless(mount_path, ssid, psk):
    wpa_supplicant = os.path.normpath(os.path.join(mount_path, 'etc/wpa_supplicant/wpa_supplicant.conf'))

    if not ssid or not psk:
        instructions = (';; 1. Follow all of the steps.\n'
                    ';; 2. Any lines starting with \';;\' will be removed.\n'
                    ';; 3. Ensure that the following exists somewhere in the file:\n'
                    ';;\n'
                    ';;    ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev\n'
                    ';;    update_config=1\n'
                    ';;    country=US\n'
                    ';;\n'
                    ';; 4. Add the following:\n'
                    ';;\n'
                    ';;    network={\n'
                    ';;        ssid="test"\n'
                    ';;        psk="testpassword"\n'
                    ';;    }\n'
                    ';;\n'
                    ';;    replacing \'test\' with your network name and \'testpassword\' with your password.\n'
                    ';; 5 (Optional). If your network is hidden, add \'scan_ssid=1\' to the \'network\' block')
        append_file(wpa_supplicant, instructions)
        wrap_cmd([EDITOR, wpa_supplicant])
        wrap_cmd(['sed', '-i', '/^;;/d', wpa_supplicant])
    else:
        print('writing file %s' % wpa_supplicant)
        with open(wpa_supplicant, 'w+') as f:
            f.write(
"""ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
country=US

network={
\tssid="%s"
\tpsk="%s"
\tscan_ssid=1
}
""" % (ssid, psk))

def fix_block(mount_path):
    kill_dir = os.path.normpath(os.path.join(mount_path, 'var/lib/systemd/rfkill/'))
    files = glob(os.path.join(kill_dir, '*:wlan'))
    print('overwriting files %s' % files)
    for file in files:
        with open(file, 'w') as f:
            f.write('0\n')

def set_default_audio_dev(mount_path):
    conf_path = os.path.normpath(os.path.join(mount_path, 'etc/asound.conf'))
    print('overwritting file %s' % conf_path)
    with open(conf_path, 'w+') as f:
        f.write(
"""pcm.!default {
    type hw
    card 1
}
ctl.!default {
    type hw
    card 1
}
""")

def setup_usb_dev(mount_path):
    conf_path = os.path.normpath(os.path.join(mount_path, 'usr/share/alsa/alsa.conf'))
    wrap_cmd(['sed', '-i', '/defaults.ctl.card 0/c\defaults.ctl.card 1', conf_path])
    wrap_cmd(['sed', '-i', '/defaults.pcm.card 0/c\defaults.pcm.card 1', conf_path])

def setup_image(args):
    unmount_device(args.dev)
    write_image(args.dev, args.image)
    parts = resolve_dev_partitions(args.dev)
    boot_mount_path = '/media/%s/bootfs' % getpass.getuser()
    root_mount_path = '/media/%s/rootfs' % getpass.getuser()
    check_mount(parts.boot, boot_mount_path)
    check_mount(parts.root, root_mount_path)

    setup_ssh(boot_mount_path)
    setup_wireless(root_mount_path, args.ssid, args.psk)
    fix_block(root_mount_path)
    set_default_audio_dev(root_mount_path)
    setup_usb_dev(root_mount_path)
    unmount_device(args.dev)

    print('SD card is ready')

########################################

parser = argparse.ArgumentParser(description='Prepare an SD Card with Raspbian.')
parser.add_argument(
    '--dev',
    dest='dev',
    type=str,
    required=True,
    help='The top-level name of the device for the SD card.')
parser.add_argument(
    '--image',
    dest='image',
    type=str,
    required=True,
    help='Path to the Raspbian image.')
parser.add_argument(
    '--ssid',
    dest='ssid',
    type=str,
    required=False,
    help='Name of the wireless network to connect to')
parser.add_argument(
    '--psk',
    dest='psk',
    type=str,
    required=False,
    help='Password to use to connect to wireless network')
parser.set_defaults(func=setup_image)

def main():
    args = parser.parse_args()
    args.func(args)

if __name__ == "__main__":
    main()