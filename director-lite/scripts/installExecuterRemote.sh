apt -y update
apt -y upgrade
rm installExecuter.sh
wget -nv --no-cache https://raw.githubusercontent.com/rayline/SSClusterManager/master/installExecuter.sh
chmod +x installExecuter.sh
setsid ./installExecuter.sh > ~/install.log 2>&1