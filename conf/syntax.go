package conf

import (
	"V-switch/tools"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func ConfCheck() {

	ConfigItems := []string{
		"MTU",
		"DEVICENAME",
		"PORT",
		"QUEUE",
		"SWITCHID",
		"TTL",
		"DEBUG",
		"DEVICEADDR",
		"DEVICEMASK",
		"PUBLIC",
		"SEED",
	}

	// everything but 		"PUBLIC" and "SEED",

	for _, c_item := range ConfigItems {

		if !ConfigItemExists(c_item) {
			log.Printf("[CONF][SYNTAX] No %s in configuration. It is mandatory", c_item)
			os.Exit(1)
		}

	}

	// now some checks for syntax

	Mtu, MtuErr := strconv.Atoi(GetConfigItem("MTU"))

	if MtuErr != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of MTU ", GetConfigItem("MTU"))
		os.Exit(1)
	}

	if port, err := strconv.Atoi(GetConfigItem("PORT")); err != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of PORT ", GetConfigItem("PORT"))
		os.Exit(1)
	} else {

		if (port > 65534) || (port < 10) {
			log.Println("[CONF][SYNTAX] Unacceptable value of PORT ", GetConfigItem("PORT"))
			os.Exit(1)
		}
	}

	if queue, err := strconv.Atoi(GetConfigItem("QUEUE")); err != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of QUEUE ", GetConfigItem("QUEUE"))
		os.Exit(1)
	} else {

		if (queue > 2048) || (queue < 16) {
			log.Println("[CONF][SYNTAX] Unacceptable value of QUEUE ", GetConfigItem("QUEUE"))
			os.Exit(1)
		}
	}

	if ttl, err := strconv.Atoi(GetConfigItem("TTL")); err != nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of TTL ", GetConfigItem("TTL"))
		os.Exit(1)
	} else {

		if (ttl > 1000) || (ttl < 30) {
			log.Println("[CONF][SYNTAX] Unacceptable value of TTL ", GetConfigItem("TTL"))
			os.Exit(1)
		}
	}

	if ip := net.ParseIP(GetConfigItem("DEVICEADDR")); ip == nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of DEVICEADDR ", GetConfigItem("DEVICEADDR"))
		os.Exit(1)

	}

	if ip := net.ParseIP(GetConfigItem("DEVICEMASK")); ip == nil {
		log.Println("[CONF][SYNTAX] Unacceptable value of DEVICEMASK ", GetConfigItem("DEVICEMASK"))
		os.Exit(1)

	}

	if len(GetConfigItem("SWITCHID")) < Mtu {
		log.Println("[CONF][SYNTAX] Unacceptable value of SWITCHID", GetConfigItem("SWITCHID"))
		log.Println("[CONF][SYNTAX] It MUST be at least the same lenght of MTU: generating one for you")

		SetConfigItem("SWITCHID", tools.RandSeq(Mtu))
		fmt.Println("SWITCHID = ", GetConfigItem("SWITCHID"))
		os.Exit(1)
	}

	if len(GetConfigItem("DEVICENAME")) > 9 {
		log.Println("[CONF][SYNTAX] Devicename too long", GetConfigItem("DEVICENAME"))
		os.Exit(1)
	}

	if len(GetConfigItem("DEVICENAME")) < 3 {
		log.Println("[CONF][SYNTAX] Devicename too short", GetConfigItem("DEVICENAME"))
		os.Exit(1)
	}

	if GetConfigItem("SEED") == "MASTER" {
		log.Println("[CONF][SYNTAX] NODE IS CONFIGURED AS MASTER, NO SEED ")

	} else {
		if _, aerr := net.ResolveUDPAddr("udp", GetConfigItem("SEED")); aerr != nil {
			log.Println("[CONF][SYNTAX] SEED is not a valid IP:PORT", GetConfigItem("SEED"))
			os.Exit(1)
		}
	}

	if h := GetConfigItem("PUBLIC"); h == "HOSTNAME" {
		log.Println("[CONF][SYNTAX] This node has its own DNS entry/public IP, NO NAT ")

	} else {
		if _, aerr := net.ResolveUDPAddr("udp", h); aerr != nil {
			log.Println("[CONF][SYNTAX] PUBLIC  is not a valid IP:PORT", h)
			os.Exit(1)
		}
	}

	log.Println("[CONF][SYNTAX] Conf syntax OK")

}
