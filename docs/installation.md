**Homebrew**

```shell
brew install manala/tap/manala
```

**Debian / Ubuntu**

as root:
```shell
curl -sSL https://raw.githubusercontent.com/manala/packages/main/manala.gpg -o /usr/share/keyrings/manala.gpg
echo "deb [signed-by=/usr/share/keyrings/manala.gpg] https://manala.github.io/packages/apt/ stable main" > /etc/apt/sources.list.d/manala.list
apt update
apt install manala
```

as user with sudo privileges:
```shell
sudo curl -sSL https://raw.githubusercontent.com/manala/packages/main/manala.gpg -o /usr/share/keyrings/manala.gpg
echo "deb [signed-by=/usr/share/keyrings/manala.gpg] https://manala.github.io/packages/apt/ stable main" | sudo tee /etc/apt/sources.list.d/manala.list
sudo apt update
sudo apt install manala
```

**Arch Linux / Manjaro**

as root:
```shell
curl -sSL https://raw.githubusercontent.com/manala/packages/main/manala.gpg | pacman-key --add -
pacman-key --lsign-key 1394DEA3
echo -e "\n[manala]\nServer = https://manala.github.io/packages/aur/stable/\$arch" >> /etc/pacman.conf
pacman -Sy manala
```

as user with sudo privileges:
```shell
curl -sSL https://raw.githubusercontent.com/manala/packages/main/manala.gpg | sudo pacman-key --add -
sudo pacman-key --lsign-key 1394DEA3
echo -e "\n[manala]\nServer = https://manala.github.io/packages/aur/stable/\$arch" | sudo tee -a /etc/apt/sources.list.d/manala.list
sudo pacman -Sy manala
```

**deb / rpm**

Download the `.deb` or `.rpm` from the releases page and install with `dpkg -i` and `rpm -i` respectively.

**Shell script**

```shell
curl -sfL https://raw.githubusercontent.com/manala/manala/main/godownloader.sh | sh
```
