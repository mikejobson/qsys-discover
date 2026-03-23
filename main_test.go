package main

import (
	"testing"
)

func TestParseQDPPacket_Device(t *testing.T) {
	xmlData := []byte(`<QDP>
  <device>
    <n>my-qsys-core</n>
    <type>core</type>
    <platform>core</platform>
    <part_number>Q-SYS Core 110f</part_number>
    <ref>device.core.my-qsys-core</ref>
    <is_virtual>false</is_virtual>
    <lan_a_mac>00:60:74:aa:bb:cc</lan_a_mac>
    <lan_a_ip>192.168.1.100</lan_a_ip>
    <lan_b_ip></lan_b_ip>
    <aux_a_ip></aux_a_ip>
    <lan_a_lldp>14:ab:ec:35:e1:7c+22</lan_a_lldp>
    <web_cfg_url>https://192.168.1.100/</web_cfg_url>
    <hw_rev>1.0</hw_rev>
  </device>
</QDP>`)

	dev, ctrl := parseQDPPacket(xmlData)
	if ctrl != nil {
		t.Errorf("expected nil control, got %+v", ctrl)
	}
	if dev == nil {
		t.Fatal("expected non-nil device, got nil")
	}
	if dev.Name != "my-qsys-core" {
		t.Errorf("Name: want %q, got %q", "my-qsys-core", dev.Name)
	}
	if dev.PartNumber != "Q-SYS Core 110f" {
		t.Errorf("PartNumber: want %q, got %q", "Q-SYS Core 110f", dev.PartNumber)
	}
	if dev.LanAIP != "192.168.1.100" {
		t.Errorf("LanAIP: want %q, got %q", "192.168.1.100", dev.LanAIP)
	}
	if dev.LanAMac != "00:60:74:aa:bb:cc" {
		t.Errorf("LanAMac: want %q, got %q", "00:60:74:aa:bb:cc", dev.LanAMac)
	}
	if dev.Ref != "device.core.my-qsys-core" {
		t.Errorf("Ref: want %q, got %q", "device.core.my-qsys-core", dev.Ref)
	}
	if dev.WebCfgURL != "https://192.168.1.100/" {
		t.Errorf("WebCfgURL: want %q, got %q", "https://192.168.1.100/", dev.WebCfgURL)
	}
}

func TestParseQDPPacket_Control(t *testing.T) {
	xmlData := []byte(`<QDP>
  <control>
    <ref>control.core.my-qsys-core</ref>
    <role>Primary</role>
    <device_ref>device.core.my-qsys-core</device_ref>
    <design_pretty>My Show Design</design_pretty>
    <design_code>abc123</design_code>
    <primary>1</primary>
    <redundant>0</redundant>
  </control>
</QDP>`)

	dev, ctrl := parseQDPPacket(xmlData)
	if dev != nil {
		t.Errorf("expected nil device, got %+v", dev)
	}
	if ctrl == nil {
		t.Fatal("expected non-nil control, got nil")
	}
	if ctrl.DesignName != "My Show Design" {
		t.Errorf("DesignName: want %q, got %q", "My Show Design", ctrl.DesignName)
	}
	if ctrl.DeviceRef != "device.core.my-qsys-core" {
		t.Errorf("DeviceRef: want %q, got %q", "device.core.my-qsys-core", ctrl.DeviceRef)
	}
	if ctrl.Role != "Primary" {
		t.Errorf("Role: want %q, got %q", "Primary", ctrl.Role)
	}
}

func TestParseQDPPacket_Invalid(t *testing.T) {
	dev, ctrl := parseQDPPacket([]byte("not xml"))
	if dev != nil || ctrl != nil {
		t.Error("expected nil device and ctrl for invalid XML")
	}
}

func TestParseQDPPacket_Empty(t *testing.T) {
	dev, ctrl := parseQDPPacket([]byte(`<QDP></QDP>`))
	if dev != nil || ctrl != nil {
		t.Error("expected nil device and ctrl for empty QDP packet")
	}
}

func TestParseQDPPacket_Peripheral(t *testing.T) {
	xmlData := []byte(`<QDP>
  <device>
    <n>my-io-frame</n>
    <type>lcqln</type>
    <platform>lcqln</platform>
    <part_number>QIO-ML2x2</part_number>
    <ref>device.ioframe.my-io-frame</ref>
    <is_virtual>false</is_virtual>
    <lan_a_mac>00:60:74:dd:ee:ff</lan_a_mac>
    <lan_a_ip>192.168.1.101</lan_a_ip>
    <lan_b_ip></lan_b_ip>
    <web_cfg_url>https://192.168.1.101/</web_cfg_url>
    <hw_rev>0.0</hw_rev>
  </device>
</QDP>`)

	dev, ctrl := parseQDPPacket(xmlData)
	if ctrl != nil {
		t.Errorf("expected nil control, got %+v", ctrl)
	}
	if dev == nil {
		t.Fatal("expected non-nil device, got nil")
	}
	if dev.Name != "my-io-frame" {
		t.Errorf("Name: want %q, got %q", "my-io-frame", dev.Name)
	}
	if dev.PartNumber != "QIO-ML2x2" {
		t.Errorf("PartNumber: want %q, got %q", "QIO-ML2x2", dev.PartNumber)
	}
}
