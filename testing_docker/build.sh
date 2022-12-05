# Get the AUTHORIZED_KEY from ~/.ham.json
# SSHPublicKey
KEY=$(cat ~/.ham.json | json_pp | grep -a "SSHPublicKey" | cut -c 22- | cut -c -176)
sudo docker build -t antonyjr/ubuntu:ham --build-arg AUTHORIZED_KEY="$KEY"
