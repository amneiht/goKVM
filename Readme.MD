# Detail 
This is simple kvm softwere write in go
</br> It no protect mechaim ,so do not use in you product
</br> Thi program can sharing clipboard data , control remote keybroad and mouse 

# Prepare
create config to use 
</br> Onserver
```
./main -print-server  > config.conf
chmod 600 config.conf
```
</br> On client 
```
./main -print-client  > config.conf
```
change the default psk to you psk
```
psk = P@ssWoRdVeryLong
```
Currently, it only supports automatic switching by hovering the mouse over the left and right sides. And it only supports one main screen.

Auto swich is setting by key "switch" on global session . You should set this value differently on the client and server 
```
[global]
switch = left    # client
------------------
[global]
switch = right  # sever
```
# Using
On remote 
```
# run sudo comanand first
sudo ls
export XAUTHORITY=$HOME/.Xauthority
export DISPLAY=":0"
sudo ./main -s -f <config file>  > /dev/null 2>&1  &
```
</br> On local machine which has keybroad and mouse :
```
 sudo ./main -c <ip remote> -f <config file>
 ```
</br> Switch machine with press right Ctrl + right Shift

# Require
uintput module
package : libx11-dev xvfb libx11-dev libxtst-dev
libxt-dev libxinerama-dev libx11-xcb-dev libxkbcommon-dev libxkbcommon-x11-dev libxkbfile-dev

