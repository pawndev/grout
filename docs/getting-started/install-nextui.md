# Installation Guide for NextUI

This guide will help you install Grout on devices running [NextUI][nextui].

## Tested Devices

Grout has been tested on the following devices running NextUI:

| Manufacturer | Device    |
|--------------|-----------|
| Anbernic     | RG35XX SP |
| TrimUI       | Brick     |
| TrimUI       | Smart Pro |

_NextUI for Anbernic's H700 devices is a prerelease. Grout works on it, but expect rough edges while NextUI finishes
that port._

## Prerequisites

- Device with NextUI installed on an SD card
- Device connected to a Wi-Fi network

## Installation Steps

### Method 1: Pak Store (Recommended)

1. Launch the NextUI Pak Store from the `Tools` menu.
2. Navigate to the `ROM Management` category.
3. Select Grout and install.
4. Launch Grout from the `Tools` menu and enjoy!

### Method 2: Manual Installation

1. Download the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout.pak.zip) for NextUI.
2. Unzip the downloaded archive.
3. Place the `Grout.pak` directory into `SD_ROOT/Tools/<platform>`, where `<platform>` is your device's NextUI
   platform folder (e.g. `tg5040` on a TrimUI Brick or Smart Pro, `my355` on a Miyoo Flip, `h700` on an Anbernic
   RG35XX SP).
4. Launch Grout from the `Tools` menu and enjoy!

## Notes

**Anbernic H700 devices.** The NextUI build for these ships an older text library that cannot draw Grout's icons, so
the Grout package carries its own in `lib/h700`, used only on those devices. Keep that folder when updating by hand.

**Save format.** Downloaded saves are named to match NextUI's **Save format** setting: `Game.gba.sav` for MinUI (the
default), `Game.srm` for RetroArch, and `Game.sav` for Generic. Grout does not rename saves already on the card
if you change the setting later. See [Supported save formats](../usage/save-sync.md#supported-save-formats).

## Update

### Pak Store update (Recommended)

1. Launch the NextUI Pak Store from the `Tools` menu.
2. If there are updates available, you will see an entry in the menu name `Available Pak Updates`. Navigate to it
3. Select Grout and update.
4. Launch Grout from the `Tools` menu and enjoy the latest version!

### In-App update

Grout has a built-in update mechanism. To update Grout, launch the application and navigate to the `Settings` menu. From there,
select `Check for Updates`. If a new version is available, follow the on-screen prompts to download and install the update.

### Manual update

To update Grout, simply download the latest release and replace the existing `Grout.pak` folder in your `SD_ROOT/Tools/<platform>` directory. If you
have made any custom configurations, ensure to back them up before replacing the folder. Be sure to keep the `config.json`
file if you do not want to authenticate again, and configure platforms folder mappings again.

## Next Steps

After installation is complete, check out the [User Guide](../usage/guide.md) to learn how to use Grout.

--8<-- "docs/_includes/cfw-links.md"
