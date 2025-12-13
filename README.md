<!-- trunk-ignore-all(markdownlint/MD033) -->
<!-- trunk-ignore(markdownlint/MD041) -->
<div align="center">

  <img src=".github/resources/isotipo.png" height="180px" width="auto" alt="romm-grout logo">
  <br />
  <img src=".github/resources/logotipo.png" height="45px" width="auto" alt="romm grout logotype">
    <h3 style="font-size: 25px;">
A RomM Client for <a href="https://nextui.loveretro.games">NextUI</a> and <a href="https://muos.dev">muOS</a>
    </h3>

<br>

[![license-badge-img]][license-badge]
[![release-badge-img]][release-badge]
[![stars-badge-img]][stars-badge]
[![downloads-badge-img]][downloads-badge]
[![discord-badge-img]][discord-badge]

<img src=".github/resources/webp/preview.webp" alt="collection flow" width="800px" height="auto">

</div>

## Features

- Download Games Wirelessly From Your RomM Instance
- Download Box Art
- Multi-file games with automatic M3U file creation
- Select multiple games at once
- Save Sync
- Optional Game Details Screen
- Optional Unzipping

### Save Sync

Save Sync allows you to store game saves on your RomM instance and use them across your devices.

How does Save Sync work?

**Downloading**

Grout will scan your local ROM and save directories on your device.

If a game on your device currently has no save file and the active RomM user has a save file for that game, Grout will
pull down that save file to your device.

If a game on your device currently has a save file and the remote save file is newer, Grout will do the following:

- Backup your current save file into a hidden `.backup` folder that is created within your saves directory.
    - This backup will have a timestamp appended to the filename.
- Download the save file from RomM and assign the proper filename to it such that it is recognized by your CFW.

**Save Sync Notes**

- Grout uploads and downloads saves associated with the active RomM user setup in Grout
    - If you have a shared user on your RomM instance that you use for Grout, be sure to be cautious with this feature.

- Grout **does not** sync save states.
    - If you use save states with auto-loading enabled, please be aware you will need to reset your emulator in order to
      load a downloaded save file.

## Installation

### NextUI Setup

Grout has been tested on the following devices running NextUI.

- TrimUI Devices
    - Brick
    - Smart Pro

1. Own a TrimUI Brick or Smart Pro and have a SD Card with NextUI.
2. Connect your device to a Wi-Fi network.
3. The preferred Grout installation method is to use the NextUI Pak Store. You can find the Pak Store in the
   `Tools` menu. Once in the Pak Store, Grout can be found under the `ROM Management` category.
    - Alternatively, download
      the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout.pak.zip) for
      NextUI
    - For manual downloads, unzip the release zip and place the `Grout.pak` directory into `SD_ROOT/Tools/tg5040`.
4. Launch `Grout` from the `Tools` menu and enjoy!

---

### muOS Setup

Grout has been tested on the following devices running muOS 2508.4 Loose Goose.

- Anbernic Devices
    - RG34XX
    - RG35XXSP
    - RG40XXV

- TrimUI Devices
    - Brick
    - Smart Pro

Please help by verifying if it works on other devices!

1. Own a supported device running muOS.
2. Download the [latest Grout release](https://github.com/rommapp/grout/releases/latest/download/Grout.muxapp) for muOS.
3. Transfer the `Grout.muxapp` file an `ARCHIVE` folder on your device.
    - `/mnt/mmc/ARCHIVE` or `/mnt/sdcard/ARCHIVE`
4. Go to Applications and launch Archive Manager.
5. Select `[SDX-APP] Grout` from the list and let it extract to your applications directory.
6. Exit Archive Manager.
7. Find an [input mapping config](/.github/resources/input_mappings) for your device.
    - If one does not exist, please try one for a different device.
    - If that does not work,
      please [create an issue](https://github.com/rommapp/grout/issues/new?template=button-mapping-request.md).
    - A first launch setup process is in the works but is not ready for primetime.
8. Save the input mapping JSON file as `input_mapping.json` and transfer it to `/MUOS/application/Grout`.
    - `/mnt/mmc/MUOS/application` or `/mnt/sdcard/MUOS/application`
9. Select `Apps` on the Main Menu, launch Grout, and enjoy!

**Note:** Grout does not support downloading art on muOS. This will be added in a future release.

## Need Help? Find a Bug? Have an Idea?

Please [create an issue](https://github.com/rommapp/grout/issues/new/choose). Be sure to fill out the template
completely!

## Spread joy!

A good chunk of time has been spent building Grout.

If you feel inclined to pay it forward, go do something nice for someone! ❤️

✌🏻

<!-- Badges -->

[license-badge-img]: https://img.shields.io/github/license/rommapp/grout?style=for-the-badge&color=007C77

[license-badge]: LICENSE

[release-badge-img]: https://img.shields.io/github/v/release/rommapp/grout?sort=semver&style=for-the-badge&color=007C77

[release-badge]: https://github.com/rommapp/grout/releases

[stars-badge-img]: https://img.shields.io/github/stars/rommapp/grout?style=for-the-badge&color=007C77

[stars-badge]: https://github.com/rommapp/grout/stargazers

[downloads-badge-img]: https://img.shields.io/github/downloads/rommapp/grout/total?style=for-the-badge&color=007C77

[downloads-badge]: https://github.com/rommapp/grout/releases

[discord-badge-img]: https://img.shields.io/badge/discord-7289da?style=for-the-badge&color=007C77

[discord-badge]: https://discord.gg/P5HtHnhUDH

<!-- Links -->

[discord-invite]: https://invidget.switchblade.xyz/P5HtHnhUDH

[discord-invite-url]: https://discord.gg/P5HtHnhUDH

[oc-donate-img]: https://opencollective.com/romm/donate/button.png?color=blue

[oc-donate]: https://opencollective.com/romm
