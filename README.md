# Simple File Manager
This is the system that we created during RAITE's hackaton.

Some features are not working or missing. This includes:
* Update is not working
* The User Interface for Copy and Move is missing
* File display and File Explorer is not yet implemented

## Root Directory

Created files are stored inside the root/ folder.

## API

The API have basic functions: Create, Read, Update, Delete, Copy and Move. Below are the description of what the URI and certain HTTP method do.

Note: This API is served in localhost:8000. Concatenate http://localhost:8000 into the URIs below
Note: Variables inside < > are required.

|     URI       |   Method   | Additional Data |   Description   |
| ------------- | ---------- | --------------- | ----------------|
| /file/test.txt | GET | None | Read the test.txt inside root/ |
| /file/test2.txt | POST | {path : "", content: [file contents]} | Creates test2.txt inside root/ |
| /file/test2.txt | PUT | {path : "", content: [file contents]} | Change the contents of test2.txt |
| /file/test2.txt | DELETE | None | Deletes test2.txt |
| /dir/newDir | POST | None | Creates new folder inside root/ |
| /file/newDir/test3.txt | POST | {path : "", content: [file contents]} | Creates test3.txt inside root/newDir/ |

## Contributors

* [Carl Joshua Biag](https://github.com/carljoshua) (Back-end)
* [Francisco Muan](https://github.com/franckiko32)  (Front-end)
