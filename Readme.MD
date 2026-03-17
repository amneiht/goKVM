# Detail 
This is simple kvm softwere write in go
</br> It no protect mechaim ,so do not use in you product
</br> Thi program can sharing clipboard data , control remote keybroad and mouse 

# Using
On remote 
```
sudo ./main -s
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
- Sharing clipboard is not work with ssh
- Add security mechanisms
- Smooth switching between multiple devices
- Only work in linux