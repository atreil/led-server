import abc
import argparse
import collections
import os
import re
import subprocess
import sys

from glob import glob

EDITOR = os.environ.get('EDITOR','vi')

Args = collections.namedtuple('Args', 'partition mount_path')

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

def check_mount(dev, mnt):
    mnt = resolve_path(mnt)
    if os.path.isfile(mnt):
        raise Exception('\'%s\' is a file' % mnt)
    if not os.path.isdir(mnt):
        os.mkdir(mnt)
    if not is_mounted(dev):
        wrap_cmd(['mount', dev, mnt])
    return mnt

########################################

############# Commands #################
def write_image(args):
    dev = args.dev
    img_path = args.image
    check_file(img_path)

    to_unmount = [part for part in glob('%s?*' % dev) if is_mounted(part)]

    if len(to_unmount) != 0:
        if not confirm('Partitions \'%s\' needs to be unmounted first. '
                    'Proceed' % dev):
                raise Exception('User canceled unmount operation')
        wrap_cmd(['umount'] + to_unmount)
    print('partions %s are unmounted' % glob('%s?*' % dev))
    wrap_cmd(['dd', 'bs=4M', 'if=%s' % img_path, 'of=%s' % dev, 'conv=fsync',
              'status=progress'])

def setup_ssh(args):
    boot_dir = check_mount(args.partition, args.mount_path)
    wrap_cmd(['touch', os.path.join(boot_dir, 'ssh')])

def setup_wireless(args, ssid, psk):
    root_dir = check_mount(args.partition, args.mount_path)
    wpa_supplicant = os.path.normpath(os.path.join(root_dir, 'etc/wpa_supplicant/wpa_supplicant.conf'))

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

def fix_block(args):
    root_dir = check_mount(args.partition, args.mount_path)
    kill_dir = os.path.normpath(os.path.join(root_dir, 'var/lib/systemd/rfkill/'))
    files = glob(os.path.join(kill_dir, '*:wlan'))
    print(root_dir)
    print('overwriting files %s' % files)
    for file in files:
        with open(file, 'w') as f:
            f.write('0\n')

def set_default_audio_dev(args):
    root_dir = check_mount(args.partition, args.mount_path)
    conf_path = os.path.normpath(os.path.join(root_dir, 'etc/asound.conf'))
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

def setup_usb_dev(args):
    root_dir = check_mount(args.partition, args.mount_path)
    conf_path = os.path.normpath(os.path.join(root_dir, 'usr/share/alsa/alsa.conf'))
    wrap_cmd(['sed', '-i', '/defaults.ctl.card 0/c\defaults.ctl.card 1', conf_path])
    wrap_cmd(['sed', '-i', '/defaults.pcm.card 0/c\defaults.pcm.card 1', conf_path])

def setup_image(args):
    boot_args = Args(partition=args.boot_partition, mount_path=args.boot_mount_path)
    root_args = Args(partition=args.root_partition, mount_path=args.root_mount_path)
    setup_ssh(boot_args)
    setup_wireless(root_args, args.ssid, args.psk)
    fix_block(root_args)
    set_default_audio_dev(root_args)
    setup_usb_dev(root_args)

########################################

parser = argparse.ArgumentParser(description='Prepare an SD Card with Raspbian')
subparsers = parser.add_subparsers()

writeimage_args = subparsers.add_parser(
    name='writeimage',
    description='Writes Raspbian image to an SD card.')
writeimage_args.add_argument(
    '--dev',
    dest='dev',
    type=str,
    required=True,
    help='Device path of the SD Card (looks like /dev/sd* - may be a common prefix).')
writeimage_args.add_argument(
    '--image',
    dest='image',
    type=str,
    required=True,
    help='Path to the Raspbian image.')
writeimage_args.set_defaults(func=write_image)

setupimage_args = subparsers.add_parser(
    name='setupimage',
    description='Sets up image by running setupssh, setupwireless, fixblock. Note that the boot partition will be the smaller than the root partition.')
setupimage_args.add_argument(
    '--root_partition',
    dest='root_partition',
    type=str,
    required=True,
    help='Device path to the rootfs partition.')
setupimage_args.add_argument(
    '--root_mount_path',
    dest='root_mount_path',
    type=str,
    required=True,
    help='Mount path to the rootfs partition.')
setupimage_args.add_argument(
    '--boot_partition',
    dest='boot_partition',
    type=str,
    required=True,
    help='Device path to the boot partition.')
setupimage_args.add_argument(
    '--boot_mount_path',
    dest='boot_mount_path',
    type=str,
    required=True,
    help='Mount path to the boot partition.')
setupimage_args.add_argument(
    '--ssid',
    dest='ssid',
    type=str,
    required=False,
    help='Name of the wireless network to connect to')
setupimage_args.add_argument(
    '--psk',
    dest='psk',
    type=str,
    required=False,
    help='Password to use to connect to wireless network')
setupimage_args.set_defaults(func=setup_image)

def main():
    args = parser.parse_args()
    args.func(args)

if __name__ == "__main__":
    main()