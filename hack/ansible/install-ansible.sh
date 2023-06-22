sudo apt-get install software-properties-common -y
sudo add-apt-repository ppa:deadsnakes/ppa -y
sudo apt-get update -y
sudo apt-get install python3.8 -y
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
sudo update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.8 1
python3 --version
python3 get-pip.py --user
python3 -m pip install --user ansible
export PATH=$PATH:/home/ubuntu/.local/bin
ansible --version
