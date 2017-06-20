package plane

import (
	"V-switch/crypt"
	"V-switch/tools"
	"log"
	"strings"
)

func init() {

	go TLVInterpreter()

}

func TLVInterpreter() {

	var my_tlv_enc []byte
	log.Println("[PLANE][TLV][INTERPRETER] Thread starts")

	for {

		my_tlv_enc = <-UdpToPlane

		my_tlv := crypt.FrameDecrypt([]byte(VSwitch.SwID), my_tlv_enc)
		if my_tlv == nil {
			continue
		}

		typ, ln, payload := tools.UnPackTLV(my_tlv)

		if ln == 0 {
			continue
		}

		switch typ {

		// it is a frame
		case "F":
			PlaneToTap <- payload
			// someone is announging itself
		case "A":
			announce := string(payload)
			fields := strings.Split(announce, "|")
			if len(fields) == 3 {
				VSwitch.AddMac(fields[0], fields[1], fields[2])
			}
		case "Q":
			sourcemac := string(payload)
			for alienmac, _ := range VSwitch.SPlane {
				AnnounceAlien(alienmac, string(sourcemac))

			}

		default:
			log.Println("[PLANE][TLV][INTERPRETER] Unknown type, discarded: [ ", typ, " ]")

		}

	}

}

func DispatchTLV(mytlv []byte, mac string) {

	mac = strings.ToUpper(mac)

	if VSwitch.macIsKnown(mac) {

		osocket := VSwitch.SPlane[mac].Socket
		log.Printf("[PLANE][TLV][DISPATCH] Dispatching to %s (%s)", mac, osocket.RemoteAddr().String())
		_, err := osocket.Write([]byte(mytlv))
		if err != nil {
			log.Println("[PLANE][TLV][DISPATCH] cannot dispatch: ", err.Error())
		}

	} else {
		log.Println("[PLANE][TLV][DISPATCH] cannot dispatch, unknown MAC ", mac)

		return
	}

}

func AnnounceLocal(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr + "|" + VSwitch.Fqdn + "|" + VSwitch.IPAdd

	log.Println("[PLANE][ANNOUNCELOCAL] Announcing  ", myannounce)

	tlv := tools.CreateTLV("A", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}

// Announces  port which is not ours
func AnnounceAlien(alien_mac string, mac string) {

	mac = strings.ToUpper(mac)
	alien_mac = strings.ToUpper(alien_mac)

	tmp_endpoint := VSwitch.SPlane[alien_mac].EndPoint
	tmp_ethIP := VSwitch.SPlane[alien_mac].EthIP

	myannounce := alien_mac + "|" + tmp_endpoint.String() + "|" + tmp_ethIP.String()

	log.Println("[PLANE][ANNOUNCEALIEN] Announcing  ", myannounce)

	tlv := tools.CreateTLV("A", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}

func SendQueryToMac(mac string) {

	mac = strings.ToUpper(mac)

	myannounce := VSwitch.HAddr

	tlv := tools.CreateTLV("Q", []byte(myannounce))

	tlv_enc := crypt.FrameEncrypt([]byte(VSwitch.SwID), tlv)

	DispatchTLV(tlv_enc, mac)

}
