# qsys-discover

A command-line tool that discovers Q-SYS devices on your local network using the Q-SYS Discovery Protocol (QDP).

## Installation

### Homebrew (macOS — recommended)

```sh
brew install mikejobson/tap/qsys-discover
```

### go install

```sh
go install github.com/mikejobson/qsys-discover@latest
```

Requires Go 1.22 or later. The binary is placed in your `$GOPATH/bin` (or `~/go/bin` by default).

### Build from source

```sh
git clone https://github.com/mikejobson/qsys-discover.git
cd qsys-discover
go build -o qsys-discover
```

## Usage

Run the tool with no arguments:

```sh
qsys-discover
```

Example output:

```
Listening for Q-SYS devices on multicast 224.0.23.175:2467...
Waiting for announcements (5s timeout)...

IP Address           Name                 Part Number          MAC Address        Design Name
----------------------------------------------------------------------------------------------------------------
192.168.1.100        my-qsys-core         Q-SYS Core 110f      00:60:74:aa:bb:cc  My Show Design
192.168.1.101        my-io-frame          QIO-ML2x2            00:60:74:dd:ee:ff

Done.
```

The tool listens for 5 seconds. All Q-SYS devices that announce themselves within that window are listed.

> **Note:** The device must be on the same LAN segment — multicast packets do not cross routers or VLAN boundaries without multicast routing configured.

## How It Works

### Discovery mechanism

`qsys-discover` passively listens for QDP (Q-SYS Discovery Protocol) multicast UDP packets. Every Q-SYS device (Core, peripheral, or a PC running Q-SYS Designer) broadcasts a QDP announcement packet roughly once per second to the multicast group `224.0.23.175` on UDP port **2467**.

No discovery packet needs to be sent — the tool simply joins the multicast group and collects announcements as devices broadcast them.

### Packet format

QDP packets are UTF-8 encoded XML. A typical device announcement looks like:

```xml
<QDP>
  <device>
    <n>my-qsys-core</n>
    <type>core</type>
    <platform>core</platform>
    <part_number>Q-SYS Core 110f</part_number>
    <ref>device.core.my-qsys-core</ref>
    <is_virtual>false</is_virtual>
    <lan_a_mac>00:60:74:aa:bb:cc</lan_a_mac>
    <lan_a_ip>192.168.1.100</lan_a_ip>
    <web_cfg_url>https://192.168.1.100/</web_cfg_url>
    <hw_rev>1.0</hw_rev>
  </device>
</QDP>
```

Q-SYS Cores also send control packets that include the active design name:

```xml
<QDP>
  <control>
    <ref>control.core.my-qsys-core</ref>
    <role>Primary</role>
    <device_ref>device.core.my-qsys-core</device_ref>
    <design_pretty>My Show Design</design_pretty>
    <primary>1</primary>
    <redundant>0</redundant>
  </control>
</QDP>
```

### Fields extracted

| Column      | XML element            |
|-------------|------------------------|
| IP Address  | `<lan_a_ip>`           |
| Name        | `<n>`                  |
| Part Number | `<part_number>`        |
| MAC Address | `<lan_a_mac>`          |
| Design Name | `<design_pretty>` (from paired `<control>` packet) |

### Requirements

- macOS or Linux
- Network interface with multicast (IGMP) support
- Same LAN segment as the Q-SYS devices, or multicast routing configured between segments
- Go 1.22+ if building from source
