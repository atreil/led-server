sudo apt-get update
sudo apt-get install git python3-numpy python3-scipy python3-pyaudio build-essential python-dev scons swig vim python3-pip golang
sudo pip3 install rpi_ws281x
cd ~
mkdir led-server
cd led-server
git clone https://github.com/atreil/audio-reactive-led-strip.git
git clone https://github.com/jgarff/rpi_ws281x.git
git clone https://github.com/atreil/led-server.git

pushd rpi_ws281x
scons
cd python
sudo python3 setup.py build
sudo python3 setup.py install
popd