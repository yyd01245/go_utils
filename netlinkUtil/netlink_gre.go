package netlinkUtil

import (
	"fmt"
	"net"
	"errors"
	log "github.com/Sirupsen/logrus"	
	NT "github.com/vishvananda/netlink"
	// "github.com/vishvananda/netlink/nl"
	// "reflect"
	"strings"
	"bytes"
	"strconv"

	"golang.org/x/sys/unix"
	"github.com/yyd01245/go_utils/files"

)

// const MODE_VPOP = "vpop"
// const MODE_POP = "pop"
// const PRIORITY = 215
const zeroCIDR = "0.0.0.0/0"

func LinkAddGreTun(ifName string,ipLocal string, ipRemote string) error{
	// list link 
	localIP := net.ParseIP(ipLocal);
	remoteIP := net.ParseIP(ipRemote);
	
	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	// log.Debugf("---links: %v",links)
	for _, l := range links {

		if l.Attrs().Name == ifName {
			txt := fmt.Sprintf("ifname link:%v exsit",ifName)
			log.Infof(txt)
			return errors.New(txt)
		}
	}
	link := &NT.Gretun{
		LinkAttrs: NT.LinkAttrs{Name: ifName},
		Local:     localIP,
		Remote:    remoteIP,
	}

	// log.Debugf("---link:%v",link)

	if err := NT.LinkAdd(link); err != nil {
		log.Errorf("Link Add interface:%v, error: %v",ifName,err)
		return err
	}

	base := link.Attrs()

	result, err := NT.LinkByName(base.Name)
	if err != nil {
		log.Errorf("Link byname error: %v",err)
		return err
	}

	rBase := result.Attrs()

	if base.Index != 0 {
		if base.Index != rBase.Index {
			txt := fmt.Sprintf("index is %d, should be %d", rBase.Index, base.Index)
			log.Errorf(txt)
			return errors.New(txt)
		}
	}

	links, err = NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	flag := false
	for _, l := range links {
		if l.Attrs().Name == link.Attrs().Name {
			log.Infof("Link gre properly:%v",l)
			flag = true
			break;
		}
	}
	if flag == false {
		log.Errorf("link gre add failed!!!")
		return errors.New("link gre add failed!!!")
	}
	// up 
	NT.LinkSetUp(link)
	return nil
}

func LinkListGreTun(ifName string,ipLocal string, ipRemote string) error{
	// list link 
	localIP := net.ParseIP(ipLocal);
	remoteIP := net.ParseIP(ipRemote);

	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	// log.Infof("---links: %v",links)
	for _, link := range links {


		if gretun, ok := link.(*NT.Gretun); ok {
			// get gre tunnel
			if gretun.Attrs().Name == ifName && 
				gretun.Local.Equal(localIP) &&
				gretun.Remote.Equal(remoteIP) {
				txt := fmt.Sprintf("find ifname same link:%v exsit",ifName)
				log.Infof(txt)
			}else if gretun.Attrs().Name == ifName {
				// find same name 
				log.Infof("find gretun same name ")
			}
		}

	}
	// link := &NT.Gretun{
	// 	LinkAttrs: NT.LinkAttrs{Name: ifName},
	// 	Local:     localIP,
	// 	Remote:    remoteIP,
	// }
	return nil
}

func LinkDelGreTun(ifName string,ipLocal string, ipRemote string) error{

	localIP := net.ParseIP(ipLocal);
	remoteIP := net.ParseIP(ipRemote);

	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	// log.Infof("---links: %v",links)
	var link NT.Link
	flag := false
	for _, l := range links {

		if gretun, ok := l.(*NT.Gretun); ok {
			// get gre tunnel
			if gretun.Attrs().Name == ifName && 
				gretun.Local.Equal(localIP) &&
				gretun.Remote.Equal(remoteIP) {
				link = l
				flag = true	
				txt := fmt.Sprintf("find ifname same link:%v exsit",ifName)
				log.Infof(txt)
				break;
			}else if gretun.Attrs().Name == ifName {
				// find same name 
				log.Warnf("find gretun same name, but local and remote unequal!!!")
				link = l
				flag = true	
				break;
			}
		}

	}

	if !flag {
		txt := fmt.Sprintf("gre tunnel:%v not exsit, delete ignore!!!",ifName)
		log.Infof(txt)
		return errors.New(txt)
	}
	if err := NT.LinkDel(link); err != nil {
		log.Errorf("Link Del error: %v",err)
		return err
	}

	links, err = NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	flag = false
	for _, l := range links {
		if l.Attrs().Name == ifName {
			log.Infof("Link gre tunnel properly:%v",l)
			flag = true
			break;
		}
	}
	if flag {
		log.Errorf("link gre tunnel del failed!!!")
		return errors.New("link gre tunnel del failed!!!")
	}
	log.Infof("delete tunnel success!!!")
	return nil
}

func DelLinkAddr(ipaddr string,ifName string) error {
	_, address, err := net.ParseCIDR(ipaddr)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}
	var addr = &NT.Addr{IPNet: address}

	link, err := NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("Link byname error: %v",err)
		return err
	}

	err = NT.AddrDel(link, addr)
	if err != nil {
		log.Errorf("AddrDel error: %v",err)
		return err
	}
	return nil
}

//ipLocal /32 ipPeer/32
func GreTunAddrAdd(ifName string,ipLocal string,ipPeer string) error {

	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	// log.Infof("---links: %v",links)
	var link NT.Link
	flag := false
	for _, l := range links {

		if gretun, ok := l.(*NT.Gretun); ok {
			// get gre tunnel
			if gretun.Attrs().Name == ifName {
				// find same name 
				log.Warnf("find gretun same name, but local and remote unequal!!!")
				link = l
				flag = true	
				break;
			}
		}
	}
	if !flag {
		txt := fmt.Sprintf("gre tunnel:%v not exsit, add addr falied!!!",ifName)
		log.Errorf(txt)
		return errors.New(txt)
	}
	local_ip, localNet, err := net.ParseCIDR(ipLocal)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}
	var address = &net.IPNet{IP: local_ip, Mask: localNet.Mask}
	ipAddr, ipNet, err := net.ParseCIDR(ipPeer)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}
	var peer = &net.IPNet{IP: ipAddr, Mask: ipNet.Mask}

	var addr = &NT.Addr{IPNet: address,Peer: peer,}

	err = NT.AddrAdd(link, addr)
	if err != nil {
		log.Errorf("AddrAdd error: %v",err)
		return err
	}
	return nil
}

// func AddRouteForGre(ifname string, routes []string, gateway string) error{
// 	link,err := GetLinkDevice(ifname)
// 	if err != nil {
// 		log.Errorf("error get device link failed!!!")
// 		return err
// 	}
// 	// 检测是否存在
// 	return NetVerifyRoutePatch(link,routes,unix.RT_TABLE_MAIN,unix.RT_SCOPE_UNIVERSE,gateway)
// }

// func AddRouteForPopGre(ifname string, routes []string, gateway string, tableID int) error{
// 	link,err := GetLinkDevice(ifname)
// 	if err != nil {
// 		log.Errorf("error get device link failed!!!")
// 		return err
// 	}

// 	route,err := GetDefaultRouteObject(link,gateway,tableID)
// 	if err = NetFindRouteByLink(link,tableID,route); err!= nil {
// 		// 
// 		log.Infof("begin add route :%v, gateway:%v",route,gateway)
// 		return NetAddOrDelRouteByLink("add",route)	
// 	}
// 	return nil
// 	// 检测是否存在
// 	// netlinkUtil.NetVerifyRoutePatch(link,routes,tableID,unix.RT_SCOPE_UNIVERSE,gateway)
// }

// func DelRouteForGre(ifname string, routes []string, gateway string) error{
// 	link,err := netlinkUtil.GetLinkDevice(ifname)
// 	if err != nil {
// 		log.Errorf("error get device link failed!!!")
// 		return err
// 	}
// 	// 检测是否存在
// 	return netlinkUtil.NetDelRoutePatch(link,routes,unix.RT_TABLE_MAIN,unix.RT_SCOPE_UNIVERSE,gateway)
// }

// func DelRouteForPopGre(ifname string, routes []string, gateway string, tableID int) error{
// 	link,err := netlinkUtil.GetLinkDevice(ifname)
// 	if err != nil {
// 		log.Errorf("error get device link failed!!!")
// 		return err
// 	}
// 	log.Debugf("begin delete default route %v",routes)
// 	route,err := netlinkUtil.GetDefaultRouteObject(link,gateway,tableID)
// 	if err = netlinkUtil.NetFindRouteByLink(link,tableID,route); err == nil {
// 		// 
// 		log.Debugf("begin add route :%v, gateway:%v",route,gateway)
// 		return netlinkUtil.NetAddOrDelRouteByLink("del",route)	
// 	}
// 	return nil
// 	// 检测是否存在
// 	// netlinkUtil.NetVerifyRoutePatch(link,routes,tableID,unix.RT_SCOPE_UNIVERSE,gateway)
// }


func DelTableIDFromName(name string,tableID int) error {
	tables,err := files.ReadFileAll("/etc/iproute2/rt_tables")
	if err != nil {
		log.Errorf("read routes get failed: %v",err)
		return err
	}
	tablesLine := strings.Split(tables,"\n")
	writeList := []string{}
	for _,value := range tablesLine {
		if value == "" {
			continue
		}
		if value == "" || ([]byte(value))[0] == '#' {
			writeList = append(writeList,value)
			continue
		}
		log.Debugf("get tables: %v",value)
		for index,v := range []byte(value) {
			log.Debugf("=====index:%d,%c!",index,v)
		}
		if strings.Index(value,name) > 0 {
			// 尝试 水平定位符号 分割
			data := []string{}
			data = strings.Split(value,"	")
			if len(data) != 2{
				// log.Errorf("tables route first error: %v,len=%v",data,len(data))
				data = strings.Split(value," ")
				if len(data) != 2{
					log.Errorf("tables route error: %v,len=%v",data,len(data))
					writeList = append(writeList,value)

					continue
				}
			}
			if data[1] == name {
				id, _ := strconv.Atoi(data[0])
				if id == tableID {
					// find
					continue
				}
			}
		}
		writeList = append(writeList,value)
	}
	log.Infof("--- get write table list: %v",writeList)
	err = files.WriteStringToFile("/etc/iproute2/rt_tables",writeList)
	if err != nil {
		log.Errorf("error append to route table %v",err)
		return err
	}
	return nil
}


func CompareIPNet(a, b *net.IPNet) bool {
	if a == b {
		return true
	}
	// For unspecified src/dst parseXfrmPolicy would set the zero address cidr
	if (a == nil && b.String() == zeroCIDR) || (b == nil && a.String() == zeroCIDR) {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.IP.Equal(b.IP) && bytes.Equal(a.Mask, b.Mask)
}

func FindAddrByLinkName(ifName string,ipLocal string) bool {
	var link NT.Link

	link, err := NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("find link device:%v error: %v",ifName,err)
		return false
	}

	local_ip, localNet, err := net.ParseCIDR(ipLocal)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return false
	}
	var address = &net.IPNet{IP: local_ip, Mask: localNet.Mask}

	addrList,err := NT.AddrList(link,NT.FAMILY_V4)
	log.Debugf("get addrList :%v",addrList)
	for _, addrData := range addrList {
		log.Debugf("--- addr = %v",addrData)
		if CompareIPNet(addrData.IPNet,address)  {
				// get 
				return true
			}

	}
	return false 
}

func VerfiyLinkInfo(ifName string,ipLocal string,ipPeer string) error{
	var link NT.Link

	link, err := NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("find link device:%v error: %v",ifName,err)
		return err
	}

	local_ip, localNet, err := net.ParseCIDR(ipLocal)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}
	var address = &net.IPNet{IP: local_ip, Mask: localNet.Mask}
	ipAddr, ipNet, err := net.ParseCIDR(ipPeer)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}
	var peer = &net.IPNet{IP: ipAddr, Mask: ipNet.Mask}
	res := errors.New("find link addr err")
	addrList,err := NT.AddrList(link,NT.FAMILY_V4)
	log.Debugf("get addrList :%v",addrList)
	for _, addrData := range addrList {
		log.Debugf("--- addr = %v",addrData)
		if CompareIPNet(addrData.IPNet,address) &&
			CompareIPNet(addrData.Peer,peer) {
				// get 
				return nil
			}

	}
	return res 

}
