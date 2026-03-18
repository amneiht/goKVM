# Detail 
This is simple kvm softwere write in go
</br> It no protect mechaim ,so do not use in you product
</br> Thi program can sharing clipboard data , control remote keybroad and mouse 

# Using
On remote 
```
# run sudo comanc first
sudo ls
export XAUTHORITY=$HOME/.Xauthority
export DISPLAY=":0"
sudo /home/ct129/app/bin/goKVM -s > /dev/null 2>&1  &
```
</br> On local machine which has keybroad and mouse :
```
 sudo ./main -c -ip remote- 
 ```
</br> Switch machine with press right Ctrl + right Shift

# Require
uintput module
package : libx11-dev , xvfb 

# Know bug and TODO
- Add security mechanisms
- Smooth switching between multiple devices
- Only work in linux
- exit program when server not availabe