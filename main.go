package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	multicastAddr = "224.0.23.175"
	qdpPort       = 2467
	timeout       = 5 * time.Second
)

// QDPPacket represents the top-level QDP XML packet.
type QDPPacket struct {
	XMLName xml.Name    `xml:"QDP"`
	Device  *QDPDevice  `xml:"device"`
	Control *QDPControl `xml:"control"`
}

// QDPDevice holds device announcement data from a QDP packet.
type QDPDevice struct {
	Name        string `xml:"n"`
	Type        string `xml:"type"`
	Platform    string `xml:"platform"`
	PartNumber  string `xml:"part_number"`
	Ref         string `xml:"ref"`
	IsVirtual   string `xml:"is_virtual"`
	LanAMac     string `xml:"lan_a_mac"`
	LanAIP      string `xml:"lan_a_ip"`
	LanBIP      string `xml:"lan_b_ip"`
	AuxAIP      string `xml:"aux_a_ip"`
	LldpInfo    string `xml:"lan_a_lldp"`
	WebCfgURL   string `xml:"web_cfg_url"`
	HwRev       string `xml:"hw_rev"`
}

// QDPControl holds control info associated with a Q-SYS Core.
type QDPControl struct {
	Ref        string `xml:"ref"`
	Role       string `xml:"role"`
	DeviceRef  string `xml:"device_ref"`
	DesignName string `xml:"design_pretty"`
	DesignCode string `xml:"design_code"`
	Primary    int    `xml:"primary"`
	Redundant  int    `xml:"redundant"`
}

// QSysDevice aggregates device and control info for display.
type QSysDevice struct {
	Name       string
	PartNumber string
	IPAddress  string
	MACAddress string
	DesignName string
	WebURL     string
	IsVirtual  string
}

func main() {
	fmt.Printf("Listening for Q-SYS devices on multicast %s:%d...\n", multicastAddr, qdpPort)
	fmt.Printf("Waiting for announcements (%s timeout)...\n", timeout)

	devices, controls := discoverDevices(timeout)

	// Build display list, deduplicated by IP address
	seen := make(map[string]bool)
	var results []QSysDevice
	for ref, dev := range devices {
		if dev.LanAIP == "" || seen[dev.LanAIP] {
			continue
		}
		seen[dev.LanAIP] = true

		d := QSysDevice{
			Name:       dev.Name,
			PartNumber: dev.PartNumber,
			IPAddress:  dev.LanAIP,
			MACAddress: dev.LanAMac,
			WebURL:     dev.WebCfgURL,
			IsVirtual:  dev.IsVirtual,
		}
		if ctrl, ok := controls[ref]; ok {
			d.DesignName = ctrl.DesignName
		}
		results = append(results, d)
	}

	if len(results) == 0 {
		fmt.Println("No devices found.")
		return
	}

	fmt.Printf("\n%-20s %-20s %-20s %-18s %-30s\n",
		"IP Address", "Name", "Part Number", "MAC Address", "Design Name")
	fmt.Println(strings.Repeat("-", 112))
	for _, d := range results {
		fmt.Printf("%-20s %-20s %-20s %-18s %-30s\n",
			d.IPAddress, d.Name, d.PartNumber, d.MACAddress, d.DesignName)
	}
	fmt.Println("\nDone.")
}

// discoverDevices joins the QDP multicast group and collects device and control
// announcements until the given timeout elapses.
func discoverDevices(dur time.Duration) (map[string]*QDPDevice, map[string]*QDPControl) {
	group := &net.UDPAddr{
		IP:   net.ParseIP(multicastAddr),
		Port: qdpPort,
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, group)
	if err != nil {
		fmt.Printf("Error joining multicast group: %v\n", err)
		fmt.Println("Hint: multicast requires a network interface with multicast support.")
		return nil, nil
	}
	defer conn.Close()

	if err := conn.SetReadBuffer(65536); err != nil {
		fmt.Printf("Warning: could not set read buffer: %v\n", err)
	}
	if err := conn.SetReadDeadline(time.Now().Add(dur)); err != nil {
		fmt.Printf("Warning: could not set read deadline: %v\n", err)
	}

	devices := make(map[string]*QDPDevice)
	controls := make(map[string]*QDPControl)
	buf := make([]byte, 65536)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			// Deadline exceeded or other error — stop listening
			break
		}
		dev, ctrl := parseQDPPacket(buf[:n])
		if dev != nil && dev.Ref != "" {
			if _, exists := devices[dev.Ref]; !exists {
				devices[dev.Ref] = dev
			}
		}
		if ctrl != nil && ctrl.DeviceRef != "" {
			if _, exists := controls[ctrl.DeviceRef]; !exists {
				controls[ctrl.DeviceRef] = ctrl
			}
		}
	}

	return devices, controls
}

// parseQDPPacket parses a raw QDP XML packet and returns device or control info.
func parseQDPPacket(data []byte) (*QDPDevice, *QDPControl) {
	var pkt QDPPacket
	if err := xml.Unmarshal(data, &pkt); err != nil {
		return nil, nil
	}
	return pkt.Device, pkt.Control
}
