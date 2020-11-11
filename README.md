# GSP

> **g**it **s**imple **p**ackager

## Config example

The following file ( .gsp-config.json ) should be in the same directory as the gsp binary and the root of your home directory. If you are on Windows that would be C:/Users/$username. The executable file ( gsp.exe ) needs to be in a directory that is on the path variable.  

To use the commits count from the HEAD to the commit you want to stop at, set the commitsCount value in your ( .gsp-config.json ) file and leave the getHashNewest and gitHashOldest values empty in your jason file.  If you want to use the hash keys to fine tune your packages here are some things you must do. The easiest way is to use the Git GUI.  
You access this by running:  
```
gitk&
```
The ampersand will allow the shell and the GUI to be open at the same time. Now how do you fill out the two keys you ask?  Here is the answer and please pay special attention.  
For the key gitHasNewest you are going to select the newest hash. For the key gitHashOldest you are going to select the oldest hash.  Do not get frustrated just run gsp -n and you can always get a preview.

```json
{
  "commitsCount": 3,
  "diffFilter": "AM",
  "directoryNames": [
    "pkg/"
  ],
  "gitHashNewest": "",
  "gitHashOldest": "",
  "outputPath": "output-root-folder",
  "repositoryPath": "git-repository-folder"
}
```

If the options `gitHashNewest` or `gitHashOldest` are not empty, then the option `commitsCount` will be ignored. 
