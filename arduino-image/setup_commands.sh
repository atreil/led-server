# Top line should be commented
if [ ${CONTINUE} -le 0 ]; then
    sudo apt-get update
    sudo apt-get install git python3-numpy python3-scipy python3-pyaudio python3-pyqtgraph build-essential python-dev scons swig vim python3-pip golang
    sudo pip3 install rpi_ws281x
    cd ~
    mkdir led-server
    pushd led-server
    git clone https://github.com/atreil/audio-reactive-led-strip.git
    git clone https://github.com/jgarff/rpi_ws281x.git
    git clone https://github.com/atreil/led-server.git
    popd
fi

pushd led-server

if [ ${CONTINUE} -le 1 ]; then
    pushd rpi_ws281x
    scons
    cd python
    sudo python3 setup.py build
    sudo python3 setup.py install
    popd
fi

if [ ${CONTINUE} -le 2 ]; then
    pushd audio-reactive-led-strip
    cd daemon
    sudo ./install.sh
    if [ $? -ne 0 ]; then
        echo "Failed to install audio-reactive daemon"
        exit 1
    fi
    popd
fi

if [ ${CONTINUE} -le 3 ]; then
    pushd led-server
    cd controller/daemon
    sudo ./install.sh

    if [ $? -ne 0 ]; then
        echo "Failed to install led-server"
        exit 1
    fi
    popd
fi

popd

echo "Done setting up Raspberry Pi"