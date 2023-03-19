# dcs-hmd

DCS-HMD is a helmet-mounted display (HMD) designed for [DCS World](https://www.digitalcombatsimulator.com/), specifically for the Ka-50 helicopter. It displays the current values of:

- Rotor pitch (from 1 to 15)
- Rotor RPM (from 0 to 110)
- Vertical velocity (from -30 to 30)

I plan to add the ability to show the current values of the following parameters at a later date:

- Airspeed
- (Radio) altitude
- Heading
- Attitude indicator (bank/pitch)

## Demo

[![Watch the video](https://markdown-videos.deta.dev/youtube/zoILcRMmNAw)](https://www.youtube.com/watch?v=zoILcRMmNAw)

## Installation

1. Download the `dcs-hmd-vX.X.X.zip` file from the [latest release](https://github.com/dimchansky/dcs-hmd/releases/latest) and extract it.

2. Copy the entire contents of the `scripts` directory to the following folder:  `%USERPROFILE%/Saved Games/DCS.openbeta/Scripts`.

3. Add the following line at the top of the `%USERPROFILE%/Saved Games/DCS.openbeta/Scripts/Export.lua` file:

       local lfs=require('lfs');dofile(lfs.writedir()..'Scripts/DCSHMD/Export.lua')`

   You can find this line in the `scripts/Export.lua.snippet` file.

4. If you have multiple monitors, run `bin/dcs-hmd.exe` on the monitor where you want the helmet-mounted display (HMD) to appear.

5. Run DCS World in **borderless windowed mode**, and select the Ka-50 helicopter mission.
