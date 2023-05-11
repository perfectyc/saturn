## /bin/bash
echo "开始安装speedtest"
sudo apt-get install curl
curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | sudo bash
sudo apt-get install speedtest
chho "开始安装docker"
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get update
sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
sudo mkdir -m 0755 -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

### 安装saturn
cd  /data && curl -s https://raw.githubusercontent.com/filecoin-saturn/L1-node/main/.env -o .env
cd  /data && curl -s https://raw.githubusercontent.com/filecoin-saturn/L1-node/main/docker-compose.yml -o docker-compose.yml
cd  /data && curl -s https://raw.githubusercontent.com/filecoin-saturn/L1-node/main/docker_compose_update.sh -o docker_compose_update.sh
cd  /data && chmod +x docker_compose_update.sh
(crontab -l 2>/dev/null; echo "*/5 * * * * cd /data && sh docker_compose_update.sh >> docker_compose_update.log 2>&1") | crontab -
