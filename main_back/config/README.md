The `cameras.json` file in this directory contains the list of all cameras (both real and fake) that the system will manage. 

Format:
```json
{
    "cameras": [
        {
            "id": "string",      // Unique identifier for the camera
            "ip": "string",      // IP address of the camera
            "username": "string", // Username for authentication
            "password": "string", // Password for authentication
            "isFake": boolean    // true if this is a fake camera, false if real
        }
    ]
}
```

To add a new camera:
1. Edit the cameras.json file
2. Add a new entry to the "cameras" array
3. Restart the server for changes to take effect

The first non-fake camera in the list will be used as the default camera for operations that require a single camera connection.
