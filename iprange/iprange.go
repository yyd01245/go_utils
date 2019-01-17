package iprange

import (
	// "fmt"
	"net"
	"fmt"
	// "errors"
	"strconv"
	"strings"
	"math/big"
	log "github.com/Sirupsen/logrus"
)

//InetAtoN IP address to int64 
func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

//GetNumberIPAddresses 获取网段的 ip 数量， /24 /30
func GetNumberIPAddresses(networkSize int) int {
	return 2 << uint(31-networkSize)
}

//convertQuardsToInt 字符串数组转 int
func convertQuardsToInt(splits []string) []int {
	quardsInt := []int{}

	for _, quard := range splits {
		j, err := strconv.Atoi(quard)
		if err != nil {
			panic(err)
		}
		quardsInt = append(quardsInt, j)
	}

	return quardsInt
}

//GetCidrIpRange 获取 cidr 格式的 IP 范围
func GetCidrIpRange(cidr string) (string, string,int) {
	log.Debugf("check ip valid: %v",cidr)
	ipv4, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Errorf("check IP valid err:%v",err)
		return "","",0
	}
	ipv4Value := ipv4.To4()
	if ipv4Value == nil {
		return "","",0
	}
	ipBegin := ipv4Value.String()

	networkSize,_ := ipv4Net.Mask.Size()

	networkQuads := convertQuardsToInt(strings.Split(ipBegin, "."))
	numberIPAddress := GetNumberIPAddresses(networkSize)
	networkRangeQuads := []string{}
	subnet_mask := 0xFFFFFFFF << uint(32-networkSize)
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[0]&(subnet_mask>>24))+(((numberIPAddress-1)>>24)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[1]&(subnet_mask>>16))+(((numberIPAddress-1)>>16)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[2]&(subnet_mask>>8))+(((numberIPAddress-1)>>8)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[3]&(subnet_mask>>0))+(((numberIPAddress-1)>>0)&0xFF)))

	return ipBegin,strings.Join(networkRangeQuads, "."),numberIPAddress

}