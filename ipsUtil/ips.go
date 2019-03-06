package ipsUtil

import (
	// "fmt"
	"net"
	// "fmt"
	// "errors"
	// "strconv"
	"strings"
	// "math/big"
	"github.com/yyd01245/go_utils/files"
	log "github.com/Sirupsen/logrus"
)

//GetMacAddrs get interface mac address
func GetMacAddrs(ifname string) string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
			log.Errorf("fail to get net interfaces: %v", err)
			return ""
	}
	macAddr := ""
	for _, netInterface := range netInterfaces {
			if netInterface.Name == ifname {
				macAddr = netInterface.HardwareAddr.String()
				break
			}
	}
	return macAddr
}

//CheckIPValid check ip address is valid 
func CheckIPValid(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		// log.Errorf("wrong ipAddr format")
		return false
	}
	ip = ip.To4()
	if ip == nil {
		// log.Errorf("wrong ipAddr to To4 format")
		return false
	}
	return true
}

// CheckPrivateIPValid 检查 cidr 格式 是否是有效的，私有地址段
// ipaddr 带掩码 192.168.0.0/24
func CheckPrivateIPValid(ipaddr string) bool {
	ret := false
	log.Debugf("check ip valid: %v",ipaddr)
	ip, _, err := net.ParseCIDR(ipaddr)
	if err != nil {
		// log.Errorf("check IP valid err:%v",err)
		return ret
	}
	ipv4Value := ip.To4()
	if ipv4Value == nil {
		return ret
	}
	ip0 := net.ParseIP("0.0.0.0")
	if ip.Equal(ip0) == true {
		return ret
	}
	return true
}

//FindIfnameByAddresses 通过 ip 查找网卡
func FindIfnameByAddresses(ipAddr string) (string,error) {
	result := ""
	ifaces, err := net.Interfaces()
	if err != nil {
			return result,err
	}

	for _, ifi := range ifaces {
		addrs, err := ifi.Addrs()
		if err != nil {
			log.Warnf("localAddresses: %v\n", err.Error())
			continue
		}

		ipInput := net.ParseIP(ipAddr)
		for _, a := range addrs {
			log.Debugf("%v -- %v\n", ifi.Name, a)
			ip, ipNet, err := net.ParseCIDR(a.String())
			if err != nil {
				// log.Errorf("check IP valid err:%v",err)
				continue
			}
			ipv4Value := ip.To4()
			if ipv4Value == nil {
				continue
			}
			if ipNet.Contains(ipInput) {
				result = ifi.Name
				log.Debugf("get ip:%v, ifname:%v",ipAddr,ifi.Name)
			}
		}
		// fmt.Printf("%v\n", ifi)
	}
	return result,err
}

// FindIfnameExclude 网卡查找除了 ifname 和 lo 网卡，并且是获取到 Ipv4 的地址
func FindIfnameExclude(Ifname string) ([]string,error) {
	result := []string{}
	ifaces, err := net.Interfaces()
	if err != nil {
			return result,err
	}

	for _, ifi := range ifaces {
		addrs, err := ifi.Addrs()
		if err != nil {
			log.Warnf("localAddresses: %v\n", err.Error())
			continue
		}

		for _, a := range addrs {
			log.Debugf("%v -- %v\n", ifi.Name, a)
			ip, _, err := net.ParseCIDR(a.String())
			if err != nil {
				// log.Errorf("check IP valid err:%v",err)
				continue
			}
			ipv4Value := ip.To4()
			if ipv4Value == nil {
				continue
			}
			if Ifname == ifi.Name || ifi.Name == "lo" {
				continue
			}
			result = append(result,ifi.Name)
		}
		// fmt.Printf("%v\n", ifi)
	}
	return result,err
}

// GetMacAddressByIP 通过 ip 地址获取 mac 地址，通过 dhcp 文件过滤的
func GetMacAddressByIP(ip string,dhcpFile string) (string,error) {
	macAddr := ""
	// 获取 mac 地址, 先通过 dhcp client 获取到当前 mac 地址
	// 如果获取不到，则通过 uci show cascade 得到 mac 地址
	output,err := files.ReadFileAll(dhcpFile)
	if err != nil {
		return macAddr,err
	}
	log.Debugf("read dhcp file %v",output)
	outputLine := strings.Split(output,"\n")
	for _,value := range outputLine {
		if strings.Index(value,ip) >= 0 {
			// find
			log.Debugf("find mac addr in dhcp file :%v",value)
			data := strings.Split(value," ")
			macAddr = data[1]
			break;
		}
	}
	return macAddr,nil
}

// GetIPByMacAddress 通过 mac 地址获取 ip，通过 dhcp 文件过滤的
func GetIPByMacAddress(macAddr string,dhcpFile string) (string,error) {
	ip := ""
	// 获取 mac 地址, 先通过 dhcp client 获取到当前 mac 地址
	// 如果获取不到，则通过 uci show cascade 得到 mac 地址
	output,err := files.ReadFileAll(dhcpFile)
	if err != nil {
		return ip,err
	}
	log.Debugf("read dhcp file %v",output)
	outputLine := strings.Split(output,"\n")
	for _,value := range outputLine {
		if strings.Index(value,macAddr) >= 0 {
			// find
			log.Debugf("find mac addr in dhcp file :%v",value)
			data := strings.Split(value," ")
			ip = data[2]
			break;
		}
	}
	return ip,nil
}

// LocalIPs 返回所有的 ipv4 地址
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

// GetIPv4ByInterface return IPv4 address from a specific interface IPv4 addresses
func GetIPv4ByInterface(name string) ([]string, error) {
	var ips []string

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

