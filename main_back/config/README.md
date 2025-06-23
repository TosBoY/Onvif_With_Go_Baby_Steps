The `cameras.csv` file in this directory contains the list of all cameras (both real and fake) that the system will manage. 

Format:
```
id,ip,username,password,isFake
camera1,192.168.1.100,admin,password123,false
camera2,192.168.1.101,admin,password456,false
fakecam1,127.0.0.1,admin,password,true
```

Column descriptions:
- `id`: Unique identifier for the camera
- `ip`: IP address of the camera
- `username`: Username for authentication
- `password`: Password for authentication
- `isFake`: Set to "true" if this is a fake camera, "false" if real

To add a new camera:
1. You can add cameras through the API endpoint
2. Or edit the cameras.csv file directly, adding a new row with the camera details
3. Restart the server for changes to take effect if edited manually

The first non-fake camera in the list will be used as the default camera for operations that require a single camera connection.
