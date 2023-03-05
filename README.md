# dcs-hmd
A helmet-mounted display (HMD) for DCS World

## Installation

1. Copy the entire contents of the `scripts` directory to the following folder:  `%USERPROFILE%/Saved Games/DCS.openbeta/Scripts`.

2. In the file `%USERPROFILE%/Saved Games/DCS.openbeta/Scripts/Export.lua`, add the following line at the end of the file:

   `local lfs=require('lfs');dofile(lfs.writedir()..'Scripts/DCSHMD/Export.lua')`

   This line can be found in the `Export.lua.snippet` file.

3. Launch `dcs-hmd.exe` and DCS World, and select a mission that includes the Ka-50 helicopter.