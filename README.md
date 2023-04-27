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

2. Open a command prompt or terminal window and navigate to the directory where you extracted the `dcs-hmd-vX.X.X.zip` file.

3. Run the `dcs-hmd.exe` program with the `-i` flag followed by the path to your DCS scripts directory (usually `%USERPROFILE%/Saved Games/DCS.openbeta/Scripts`). For example:

       dcs-hmd.exe -i "%USERPROFILE%/Saved Games/DCS.openbeta/Scripts"

   This will automatically install the required scripts in the specified DCS scripts directory.

4. If you have multiple monitors, run `dcs-hmd.exe` on the monitor where you want the helmet-mounted display (HMD) to appear.

5. Run DCS World in **borderless windowed mode**, and select the Ka-50 helicopter mission.
