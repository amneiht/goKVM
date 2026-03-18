# Detail 
This is simple kvm softwere write in go
</br> It no protect mechaim ,so do not use in you product
</br> Thi program can sharing clipboard data , control remote keybroad and mouse 

# Prepare
create config to use 
```
./main -print  > config.conf
chmod 600 config.conf
```
change the default psk to you psk
```
psk = P@ssWoRdVeryLong
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
package : libx11-dev , xvfb 

# Know bug and TODO

- compress mouse move evnet
- Smooth switching between multiple devices
- Only work in linux
