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
	"strconv"
	"golang.org/x/sys/unix"
	"github.com/yyd01245/go_utils/files"
)

const ADDROUTE = "add"
const DELROUTE = "del"

func LinkAddMacVlan(ifName string,parentIfname string) error{
	// list link 
	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	log.Infof("---links: %v",links)
	var parent NT.Link
	for _, l := range links {
		if l.Attrs().Name == parentIfname {
			// get parent link
			parent = l
		}
		if l.Attrs().Name == ifName {
			txt := fmt.Sprintf("ifname link:%v exsit",ifName)
			log.Infof(txt)
			return errors.New(txt)
		}
	}
	link := &NT.Macvlan{
		LinkAttrs: NT.LinkAttrs{Name: ifName, ParentIndex: parent.Attrs().Index},
		Mode:      NT.MACVLAN_MODE_BRIDGE,
	}
	log.Infof("---link:%v",link)

	if err := NT.LinkAdd(link); err != nil {
		log.Errorf("Link Add interface:%v, parent interface:%v error: %v",ifName,parentIfname,err)
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
			log.Infof("Link macvlan properly:%v",l)
			flag = true
			break;
		}
	}
	if flag == false {
		log.Errorf("link macvlan add failed!!!")
		return errors.New("link macvlan add failed!!!")
	}
	// up 
	NT.LinkSetUp(link)
	return nil
}

func LinkDelMacVlan(ifName string) error{
	// list link 
	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	log.Infof("---links: %v",links)
	var link NT.Link
	flag := false
	for _, l := range links {

		if l.Attrs().Name == ifName {
			txt := fmt.Sprintf("ifname link:%v exsit",ifName)
			log.Infof(txt)
			link = l
			flag = true
			break;
		}
	}
	if !flag {
		txt := fmt.Sprintf("virtural macvlan:%v not exsit!!!",ifName)
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
			log.Infof("Link macvlan properly:%v",l)
			flag = true
			break;
		}
	}
	if flag {
		log.Errorf("link macvlan del failed!!!")
		return errors.New("link macvlan del failed!!!")
	}
	return nil
}

func LinkAddVlan(vlanID int,ifName string,parentIfname string) error {

	// list link 
	links, err := NT.LinkList()
	if err != nil {
		log.Errorf("Link list error: %v",err)
		return err
	}
	log.Infof("---links: %v",links)
	var parent NT.Link
	for _, l := range links {
		if l.Attrs().Name == parentIfname {
			// get parent link
			parent = l
		}
		if l.Attrs().Name == ifName {
			txt := fmt.Sprintf("ifname link:%v exsit",ifName)
			log.Infof(txt)
			return errors.New(txt)
		}
	}
	// &Vlan{LinkAttrs{Name: "bar", ParentIndex: parent.Attrs().Index}, 900})
	link := &NT.Vlan{
		LinkAttrs: NT.LinkAttrs{Name: ifName, ParentIndex: parent.Attrs().Index},
		VlanId: vlanID,
	}
	log.Infof("---link:%v",link)

	if err := NT.LinkAdd(link); err != nil {
		log.Errorf("Link Add interface:%v, parent interface:%v error: %v",ifName,parentIfname,err)
		return err
	}
	_, err = NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("Link byname error: %v",err)
		return err
	}
	// up 
	// NT.LinkSetUp(link)
	if err = NT.LinkSetUp(link); err != nil {
		log.Errorf("setup link device:%v error: %v",ifName,err)
		return err
	}

	return nil
}

// SetLinkAddr 绑定地址为 ip/mask
func SetLinkAddr(ipaddr string,ifName string) error {
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

	err = NT.AddrAdd(link, addr)
	if err != nil {
		log.Errorf("AddrAdd error: %v",err)
		return err
	}
	return nil
}

func DeleteLink(ifName string) error {
	var link NT.Link

	link, err := NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("find link device:%v error: %v",ifName,err)
		return err
	}
	if err := NT.LinkDel(link); err != nil {
		log.Errorf("Link Del link error: %v",err)
		return err
	}
	link, err = NT.LinkByName(ifName)
	if err == nil {
		txt := fmt.Sprintf("delete link device:%v failed!",ifName)
		log.Errorf(txt)
		return errors.New(txt)
	}
	return nil
}

func GetLinkDevice(ifName string) (NT.Link,error) {
	var link NT.Link

	link, err := NT.LinkByName(ifName)
	if err != nil {
		log.Errorf("find link device:%v error: %v",ifName,err)
		return link,err
	}
	// bring the interface up
	if err = NT.LinkSetUp(link); err != nil {
		log.Errorf("setup link device:%v error: %v",ifName,err)
		return link,err
	}
	return link,nil
}

func GetRouteObject(link NT.Link,ipaddr string,outAddr string,tableID int) (*NT.Route,error) {
	_, dst, err := net.ParseCIDR(ipaddr)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return nil,errors.New(txt)
	}
	log.Debugf("dst cidr=%v",dst)
	src := net.ParseIP(outAddr)

	log.Debugf("src ip=%v",src)

	route := &NT.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		// Src:       src,
		Gw:       src,
		// Scope:     unix.RT_SCOPE_LINK,
		// Priority:  13,
		Table:     tableID,
		// Type:      unix.RTN_UNICAST,
		// Tos:       14,
	}
	log.Debugf("--- route: %v",route)
	return route,nil 
}

func GetDefaultRouteObject(link NT.Link,outAddr string,tableID int) (*NT.Route,error) {

	src := net.ParseIP(outAddr)

	log.Infof("src ip=%v",src)

	route := &NT.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       nil,
		// Src:       src,
		Gw:       src,
		// Scope:     unix.RT_SCOPE_LINK,
		Scope: unix.RT_SCOPE_UNIVERSE,
		// Priority:  13,
		Table:     tableID,
		// Type:      unix.RTN_UNICAST,
		// Tos:       14,
	}
	log.Infof("--- route: %v",route)
	return route,nil 
}

func NetAddOrDelRouteByLink(action string,route *NT.Route) error{
	
	// // add a gateway route
	if route == nil {
		txt := fmt.Sprintf("add route error: route is nil")
		log.Errorf(txt)
		return errors.New(txt)
	}
	if action == ADDROUTE {
		log.Debugf("---add route: %v",route)
		if err := NT.RouteReplace(route); err != nil {
			txt := fmt.Sprintf("add route error:%v",err)
			log.Warnf(txt)
			return errors.New(txt)
		}
	}else if action == DELROUTE {
		// list 
		// routes, err := NT.RouteListFiltered(NT.FAMILY_V4, route, NT.RT_FILTER_DST|NT.RT_FILTER_GW|NT.RT_FILTER_TABLE)
		// if err != nil {
		// 	txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		// 	log.Errorf(txt)
		// 	return errors.New(txt)
		// }	
		// log.Infof("route list: %v",routes)
	
		// if len(routes) != 1 {
		// 	txt := fmt.Sprintf("Route not added properly:%v",routes)
		// 	log.Errorf(txt)
		// 	return errors.New(txt)
		// }
		log.Debugf("---delete route: %v",route)
		if err := NT.RouteDel(route); err != nil {
			txt := fmt.Sprintf("add route error:%v",err)
			log.Errorf(txt)
			return errors.New(txt)
		}
	}else {
		txt := "route unkown action type"
		log.Errorf(txt)
		return errors.New(txt)
	}
	return nil
}

// 批量添加路由，不进行匹配添加
func NetAddRoutePatch(link NT.Link,dstRoutes []string,outAddr string,tableID int) error{
	total := 0
	for _,ipAddr := range dstRoutes {
		// 
		route,err := GetRouteObject(link,ipAddr,outAddr,tableID)
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		err = NetAddOrDelRouteByLink("add",route)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("add route total:%v!",total)
	return nil
}

func NetSyncSopeLinkRouteTable(link NT.Link,srcTableID int,dstTableID int) error{
	total := 0
	// 不检查 default route,仅校验 scope link
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		unix.RT_SCOPE_LINK,
		Table:     srcTableID,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	
	dstRoute := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		unix.RT_SCOPE_LINK,
		Table:     dstTableID,
	}

	dstRoutes, err := NT.RouteListFiltered(NT.FAMILY_V4, &dstRoute, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	

	for _,R := range routes {
		flag := false
		if R.Dst == nil {
			// default
			continue
		}
		for index,r := range dstRoutes {
			if r.Dst == nil {
				continue
			}
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst.String() == r.Dst.String() {
				flag = true
				log.Infof("find route:%v, route:%v",r,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		R.Table = dstTableID 
		err = NetAddOrDelRouteByLink("add",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}
	log.Infof("NetSyncSopeLinkRouteTable del route total:%v!",total)
	total = 0
	// 添加
	// log.Infof("last need to delete scope link route:%v,len=%v!",dstRoutes,len(dstRoutes))
	for _,route := range dstRoutes {
		if route.Dst == nil {
			continue
		}
		log.Infof("need to delete scope link  ip:%v",route)
		err = NetAddOrDelRouteByLink("del",&route)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("delete scope link route total:%v!",total)

	return nil
}

func NetVerfiyRouteTable(srcTableID int,dstTableID int) error{
	total := 0
	// 不检查 default route,校验 scope link 和 明细路由
	route := NT.Route{
		Scope: 		unix.RT_SCOPE_LINK|unix.RT_SCOPE_UNIVERSE,
		Table:     srcTableID,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered RT_SCOPE_LINK error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	
	route = NT.Route{
		Scope: 		unix.RT_SCOPE_UNIVERSE,
		Table:     srcTableID,
	}
	routesUN, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
		NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered RT_SCOPE_UNIVERSE error:%v",err)
		log.Errorf(txt)
		// return errors.New(txt)
	}	
	routes = append(routes,routesUN...)

	dstRoute := NT.Route{
		Scope: 		unix.RT_SCOPE_LINK,
		Table:     dstTableID,
	}

	dstRoutes, err := NT.RouteListFiltered(NT.FAMILY_V4, &dstRoute, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered RT_SCOPE_LINK error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	dstRoute = NT.Route{
		Scope: 		unix.RT_SCOPE_UNIVERSE,
		Table:     dstTableID,
	}

	dstRoutesUN, err := NT.RouteListFiltered(NT.FAMILY_V4, &dstRoute, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		// return errors.New(txt)
	}	
	dstRoutes = append(dstRoutes,dstRoutesUN...)
	
	for _,R := range routes {

		// log.Infof("----- get route scope link routes: linkIndex:%v,Ilinkindex:%v,Scope:%v,dst:%v,src:%v,gw:%v,MultiPath:%v,Protocol:%v,Priority:%v,Table:%v,Type:%v,Tos:%v,Flags:%v,MplsDst:%v,NewDst:%v,Encap:%v,MTU:%v,AdvMss:%v!!!",
		// R.LinkIndex,R.ILinkIndex,R.Scope,R.Dst,R.Src,R.Gw,R.MultiPath,R.Protocol,R.Priority,
		// R.Table,R.Type,R.Tos,R.Flags,R.MPLSDst,R.NewDst,R.Encap,R.MTU,R.AdvMSS)
		flag := false
		if R.Dst == nil || (R.Flags == unix.RTNH_F_LINKDOWN ) {
			// default
			continue
		}
		for index,r := range dstRoutes {
			if r.Dst == nil {
				continue
			}
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst.String() == r.Dst.String() {
				flag = true
				log.Infof("find route:%v, route:%v",r,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		R.Table = dstTableID 
		err = NetAddOrDelRouteByLink("add",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}
	log.Infof("NetVerfiyRouteTable del route total:%v!",total)
	total = 0
	// 添加
	// log.Infof("last need to delete scope link route:%v,len=%v!",dstRoutes,len(dstRoutes))

	log.Infof("delete scope link route total:%v!",total)

	return nil
}

func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{}  // 存放不重复主键
	for _, e := range slc{
			l := len(tempMap)
			tempMap[e] = 0
			if len(tempMap) != l{  // 加入map后，map长度变化，则元素不重复
					result = append(result, e)
			}
	}
	return result
}


func ClearRoutePatch(link NT.Link,tableID int) error{
	total := 0

	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Table:     tableID,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_OIF)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	for _,R := range routes {
		// log.Infof("delete route: %v!",R)
		err = NetAddOrDelRouteByLink("del",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}
	log.Infof("ClearRoutePatch del route total:%v!",total)

	return nil
}

func NetSyncRoutePatch(link NT.Link,data []string,tableID int,scope int,gwIP string) error{
	total := 0
	gw := net.ParseIP(gwIP)

	dstRoutes := RemoveRepByMap(data)
	log.Infof("remoce repeat data: len=%d",len(dstRoutes))
	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:			gw,	
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	for _,R := range routes {
		flag := false
		for index,ipAddr := range dstRoutes {
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst == nil {
				continue
			}
			if R.Dst.String() == ipAddr {
				flag = true
				log.Debugf("find route:%v, route:%v",ipAddr,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		err = NetAddOrDelRouteByLink("del",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}
	log.Infof("NetSyncRoutePatch del route total:%v!",total)
	total = 0
	// 添加
	log.Infof("last need to add route len=%v!",len(dstRoutes))
	for _,ipAddr := range dstRoutes {
		// 
		if ipAddr == "" {
			continue
		}
		log.Debugf("begin add ip:%v",ipAddr)
		route,err := GetRouteObject(link,ipAddr,gwIP,tableID)
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		err = NetAddOrDelRouteByLink("add",route)
		if err != nil {
			log.Warnf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("NetSyncRoutePatch add route total:%v!",total)

	return nil
}

func NetVerifyRoutePatch(link NT.Link,dstRoutes []string,tableID int,scope int,gwIP string) error{
	total := 0
	gw := net.ParseIP(gwIP)

	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:			gw,	
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	

	for _,R := range routes {
		flag := false
		for index,ipAddr := range dstRoutes {
			if R.Dst == nil {
				continue
			}
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst.String() == ipAddr {
				flag = true
				log.Infof("find route:%v, route:%v",ipAddr,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if flag {
			continue
		}
	}
	total = 0
	// 添加
	log.Infof("NetVerifyRoutePatch last need to add route:%v,len=%v!",dstRoutes,len(dstRoutes))
	for _,ipAddr := range dstRoutes {
		// 
		if ipAddr == "" {
			continue
		}
		log.Infof("begin add ip:%v",ipAddr)
		route,err := GetRouteObject(link,ipAddr,gwIP,tableID)
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		err = NetAddOrDelRouteByLink("add",route)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",route,err)
			continue
		}
		total++
	}
	log.Infof("NetVerifyRoutePatch add route total:%v!",total)

	return nil
}

// 批量删除路由，不进行匹配保留项，排查 scope
func NetDelRoutePatch(link NT.Link,dstRoutes []string,tableID int,scope int,gwIP string) error{
	total := 0
	gw := net.ParseIP(gwIP)

	// 不能删除 default 和 scope link 以及保留的地址
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:			gw,	
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	for _,R := range routes {
		flag := false
		for index,ipAddr := range dstRoutes {
			// 找到则 60.191.85.65/32 格式为带掩码
			if R.Dst == nil {
				continue
			}
			if R.Dst.String() == ipAddr {
				flag = true
				log.Infof("find route:%v, route:%v",ipAddr,R)
				dstRoutes = append(dstRoutes[:index],dstRoutes[index+1:]...)
				break
			}
		}
		if !flag {
			continue
		}
		// log.Infof("delete route: %v!",R)
		err = NetAddOrDelRouteByLink("del",&R)
		if err != nil {
			log.Errorf("add route:%v failed:%v!!",R,err)
			continue
		}
		total++
	}


	log.Infof("NetDelRoutePatch del route total:%v!",total)
	return nil
}

func NetFindRouteByLink(link NT.Link,tableID int,route *NT.Route) error{


	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_DST|NT.RT_FILTER_SRC|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	log.Infof("route list: %v",routes)
	log.Infof("route list len: %v",len(routes))
	if len(routes) != 1 {
		txt := fmt.Sprintf("RouteListFiltered failed:%v",route)
		log.Errorf(txt)
		return errors.New(txt)
	}
	return nil
}


func NetListRouteByLink(link NT.Link,tableID int,scope int,gwIP string) error{

	// add a gateway route
	gw := net.ParseIP(gwIP)
	route := NT.Route{
		LinkIndex: link.Attrs().Index,
		Scope: 		NT.Scope(scope),
		Table:     tableID,
		Gw:				gw,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE|NT.RT_FILTER_GW)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	log.Infof("route list: %v",routes)
	log.Infof("route list len: %v",len(routes))
	
	return nil
}

func NetListALLRoute(tableID int,scope int) error{


	route := NT.Route{
		Scope: 		NT.Scope(scope),
		Table:     tableID,
	}

	routes, err := NT.RouteListFiltered(NT.FAMILY_V4, &route, 
			NT.RT_FILTER_TABLE|NT.RT_FILTER_SCOPE)
	if err != nil {
		txt := fmt.Sprintf("RouteListFiltered error:%v",err)
		log.Errorf(txt)
		return errors.New(txt)
	}	
	log.Infof("route list: %v",routes)
	log.Infof("route list len: %v",len(routes))
	for i,v := range routes {
		log.Infof("route index:%d, string:%s,struct:%v",i,v.String(),v)
	}
	return nil
}

type OutInfo struct {
	OutAddr string
	DevName string
	Weight  int
}

func GetRouteMultiPathObject(ipaddr string,outPath []OutInfo,tableID int) (*NT.Route,error) {
	_, dst, err := net.ParseCIDR(ipaddr)
	if err != nil {
		txt := fmt.Sprintf("check IP valid err:%v",err)
		log.Errorf(txt)
		return nil,errors.New(txt)
	}
	log.Debugf("dst cidr=%v",dst)

	length := len(outPath)
	if length > 1 {
		NexthopList := []*NT.NexthopInfo{}

		for _,data := range outPath {
			gw := net.ParseIP(data.OutAddr)
			link,err := GetLinkDevice(data.DevName)
			if err != nil {
				log.Errorf("error get device link:%v failed:%v!!!",data.DevName,err)
				continue
			}
			hop := &NT.NexthopInfo{
				LinkIndex: link.Attrs().Index,
				Gw:        gw,
				Hops:      data.Weight,
			}
			NexthopList = append(NexthopList,hop)
		}
		log.Debugf("NexthopList:%v",NexthopList)
		route := &NT.Route{
			// LinkIndex: link.Attrs().Index,
			Dst:       dst,
			// Src:       src,
			// Gw:       src,
			// Scope:     unix.RT_SCOPE_LINK,
			// Priority:  13,
			Table:     tableID,
			MultiPath: NexthopList,
			// Type:      unix.RTN_UNICAST,
			// Tos:       14,
		}
		log.Debugf("--- route: %v",route)
		return route,nil 
	}else if length == 1{
		gw := net.ParseIP(outPath[0].OutAddr)
		link,err := GetLinkDevice(outPath[0].DevName)
		if err != nil {
			txt := fmt.Sprintf("error get device link:%v failed:%v!!!",outPath[0].DevName,err)
			log.Errorf(txt)
			return nil,errors.New(txt)
		}
		log.Debugf("gw ip=%v",gw)
		route := &NT.Route{
			LinkIndex: link.Attrs().Index,
			Dst:       dst,
			// Src:       src,
			Gw:       gw,
			// Scope:     unix.RT_SCOPE_LINK,
			// Priority:  13,
			Table:     tableID,
			// Type:      unix.RTN_UNICAST,
			// Tos:       14,
		}
		log.Debugf("--- route: %v",route)
		return route,nil 
	}
	txt := fmt.Sprintf("check IP valid err:%v",err)
	log.Errorf(txt)
	return nil,errors.New(txt)
}

// 批量添加路由，不进行匹配直接替换
func NetReplaceRoutePatch(dstRoutes []string,outData []OutInfo,tableID int) error{
	total := 0
	var route *NT.Route
	for _,ipAddr := range dstRoutes {
		route,_ = GetRouteMultiPathObject(ipAddr,outData,tableID)
		if route != nil {
			break
		}
	}

	for _,ipAddr := range dstRoutes {
		// 
		_, dst, err := net.ParseCIDR(ipAddr)
		if err != nil {
			txt := fmt.Sprintf("check IP valid err:%v",err)
			log.Errorf(txt)
			continue
		}
		route.Dst = dst
		if err != nil {
			log.Errorf("get route from ipAddr:%s error!!",ipAddr)
			continue
		}
		log.Debugf("---repalce route: %v",route)
		if err := NT.RouteReplace(route); err != nil {
			txt := fmt.Sprintf("add route error:%v",err)
			log.Errorf(txt)
			return errors.New(txt)
		}
		total++
	}
	log.Infof("add route total:%v!",total)
	return nil
}

type NetRule struct {
	Priority          int
	Family            int
	Table             int
	Mark              int
	Mask              int
	TunID             uint
	Goto              int
	Src               *net.IPNet
	Dst               *net.IPNet
	Flow              int
	IifName           string
	OifName           string
	SuppressIfgroup   int
	SuppressPrefixlen int
	Invert            bool
}

func NetGetRuleObject(data NetRule) *NT.Rule {
	rule := NT.NewRule()
	rule.Table = data.Table
	rule.Src = data.Src
	rule.Dst = data.Dst
	rule.Priority = data.Priority
	rule.OifName = data.OifName
	rule.IifName = data.IifName
	rule.Invert = data.Invert
	rule.Mark = data.Mark
	return rule
}


func GetTableIDFromName(name string) int {
	tables,err := files.ReadFileAll("/etc/iproute2/rt_tables")
	if err != nil {
		log.Errorf("read routes get failed: %v",err)
		return -1
	}
	tablesLine := strings.Split(tables,"\n")
	for _,value := range tablesLine {
		if value == "" || ([]byte(value))[0] == '#' {
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
					continue
				}
			}
			if data[1] == name {
				id, _ := strconv.Atoi(data[0])
				return id
			}
		}

	}
	return -1
}

func DeleteTableIDFromName(name string, tableID int) error {
	tables,err := files.ReadFileAll("/etc/iproute2/rt_tables")
	if err != nil {
		log.Errorf("read routes get failed: %v",err)
		return err
	}
	tablesLine := strings.Split(tables,"\n")
	writeData := []string{}
	needClear := false
	for _,value := range tablesLine {
		if value == "" {
			continue
		}
		if ([]byte(value))[0] == '#' {
			writeData = append(writeData,value)
			continue
		}
		log.Debugf("get tables: %v",value)
		for index,v := range []byte(value) {
			log.Debugf("=====index:%d,%c!",index,v)
		}
		if strings.Index(value,name) >= 0 {
			// 尝试 水平定位符号 分割
			data := []string{}
			data = strings.Split(value,"	")

			if len(data) != 2{
				// log.Errorf("tables route first error: %v,len=%v",data,len(data))
				data = strings.Split(value," ")
				if len(data) < 2{
					log.Infof("delete tables route error: %v,len=%v",data,len(data))
					writeData = append(writeData,value)
					continue
				}
			}
			if data[1] == name {
				id, _ := strconv.Atoi(data[0])
				if id == tableID {
					needClear = true
					log.Infof("find tablename=%v, id=%v",name, tableID)
					continue;
				}else {
					log.Infof("find != tablename=%v, id=%v",name, id)
				}
			}
		}
		writeData = append(writeData,value)

	}

	if needClear {
		writeString := strings.Join(writeData, "\n")
		writeString = writeString+"\n"
		log.Infof("need rewrite to file ")
		err := files.WriteStringToFile("/etc/iproute2/rt_tables",writeString)	
		if err != nil {
			// return err
			log.Errorf("write to upwan config error1!");
			return err
		}
	}
	return nil 
	
}

func NetGetRuleObjectList(dstRule []NetRule) []*NT.Rule {
	// 注意 外部的默认值是 0 的话，会修改数据，数值的默认值为 -1
	ruleList := []*NT.Rule{}
	for _,data := range dstRule {
		rule := NT.NewRule()
		rule.Table = data.Table
		rule.Src = data.Src
		rule.Dst = data.Dst
		rule.Priority = data.Priority
		rule.OifName = data.OifName
		rule.IifName = data.IifName
		rule.Invert = data.Invert
		rule.Mark = data.Mark
		ruleList = append(ruleList,rule)
	}
	// for i := range ruleList {
	// 	log.Infof("index:%v,table:%v,Src:%v,Dst:%v,OifName:%v,Prio:%v,IifName:%v,Invert:%v,mark:%v,goto:%v,mask=%v! rule:%v!",
	// 	i,ruleList[i].Table,ruleList[i].Src,ruleList[i].Dst,ruleList[i].OifName,ruleList[i].Priority,
	// 	ruleList[i].IifName,ruleList[i].Invert,ruleList[i].Mark,ruleList[i].Goto,ruleList[i].Mask,ruleList[i])
	// }
	return ruleList
}


func NetAddorDelRule(action string,rule *NT.Rule) error {
	if rule == nil {
		txt := "NetAddorDelRule rule is nil"
		log.Errorf(txt)
		return errors.New(txt)
	}
	if action == "add" {
		if err := NT.RuleAdd(rule); err != nil {
			log.Errorf("NetAddorDelRule add rule error: %v",err)
			return err
		}
	}else {
		if err := NT.RuleDel(rule); err != nil {
			log.Errorf("NetAddorDelRule del rule error: %v",err)
			return err
		}
	}
	return nil 
}

func IpNetEqual(ipn1 *net.IPNet, ipn2 *net.IPNet) bool {
	if ipn1 == ipn2 {
		return true
	}
	if ipn1 == nil || ipn2 == nil {
		return false
	}
	m1, _ := ipn1.Mask.Size()
	m2, _ := ipn2.Mask.Size()
	return m1 == m2 && ipn1.IP.Equal(ipn2.IP)
}

func IsEqualRule(value *NT.Rule,rule *NT.Rule) bool {
	flag := false
	// if value.SuppressIfgroup != rule.SuppressIfgroup {
	// 	log.Infof("SuppressIfgroup not equal: %v,%v",value.SuppressIfgroup,rule.SuppressIfgroup)
	// 	return false
	// }
	// if value.SuppressPrefixlen != rule.SuppressPrefixlen {
	// 	log.Infof("SuppressPrefixlen not equal: %v,%v",value.SuppressPrefixlen,rule.SuppressPrefixlen)
	// 	return false
	// }
	// if value.Family != rule.Family {
	// 	log.Infof("Family not equal: %v,%v",value.Family,rule.Family)
	// 	return false
	// }
	// if value.TunID != rule.TunID {
	// 	log.Infof("TunID not equal: %v,%v",value.TunID,rule.TunID)
	// 	return false
	// }

	// if value.Flow != rule.Flow {
	// 	log.Infof("Flow not equal: %v,%v",value.Flow,rule.Flow)
	// 	return false
	// }
	// if !IpNetEqual(value.Src,rule.Src) || !IpNetEqual(value.Dst,rule.Dst) {
	// 	log.Infof("IpNetEqual not equal:")
	// 	return false
	// }
	// if value.Mark != rule.Mark {
	// 	log.Infof("Mark not equal: %v,%v",value.Mark,rule.Mark)
	// 	return false
	// }
	// amd64 !=
	// if value.Mask != rule.Mask {
	// 	log.Infof("Mask not equal: %v,%v",value.Mask,rule.Mask)
	// 	return false
	// }

	if value.Table == rule.Table &&
		IpNetEqual(value.Src,rule.Src) && 
		IpNetEqual(value.Dst,rule.Dst) &&
		value.Family == rule.Family &&
		value.TunID == rule.TunID &&
		value.Flow == rule.Flow &&
		value.SuppressIfgroup == rule.SuppressIfgroup &&
		value.SuppressPrefixlen == rule.SuppressPrefixlen &&
		value.OifName == rule.OifName &&
		value.Priority == rule.Priority &&
		value.IifName == rule.IifName &&
		value.Mark == rule.Mark &&
		value.Goto == rule.Goto &&
		// amd64 mask error
		// value.Mask == rule.Mask &&
		value.Invert == rule.Invert {
		flag = true
		log.Debugf("find rule: %v",rule)
	}

	return flag
}

// NetVerifyExistRuleList 校验 rulelist 规则存在，没有则增加
func NetVerifyExistRuleList(dstRules []*NT.Rule) error {
	// 
	rules, err := NT.RuleList(unix.AF_INET)
	if err != nil {
		log.Errorf("ListAllRule error:%v",err)
		return err
	}
	log.Infof("dstRules:%v!",dstRules)
	log.Infof("---list len=%d!",len(rules))
	// find this rule
	for i,value := range rules {
		log.Debugf("index:%v,table:%v,Src:%v,Dst:%v,OifName:%v,Prio:%v,IifName:%v,Invert:%v,mark:%v,goto:%v,mask=%v! rule:%v!",
		i,value.Table,value.Src,value.Dst,value.OifName,value.Priority,
		value.IifName,value.Invert,value.Mark,value.Goto,value.Mask,value)
		for index,rule := range dstRules{
			if IsEqualRule(&value,rule) {
				log.Infof("find rule: %v",rule)
				dstRules = append(dstRules[:index],dstRules[index+1:]...)
				break
			}
		}
	}
	total := 0
	for _,rule := range dstRules{
		log.Infof("need add rule :%v",rule)
		NetAddorDelRule("add",rule)
		total++
	}
	log.Infof("add rule total: %v",total)
	return nil 
}

// NetVerifyNotExistRuleList 校验 rulelist 规则不存在，有则删除
func NetVerifyNotExistRuleList(dstRules []*NT.Rule) error {
	// 
	rules, err := NT.RuleList(unix.AF_INET)
	if err != nil {
		log.Errorf("ListAllRule error:%v",err)
		return err
	}
	log.Infof("---list len=%d, dstrule len=%d!",len(rules),len(dstRules))
	// find this rule
	total := 0
	for i,value := range rules {
		log.Debugf("index:%v,table:%v,Src:%v,Dst:%v,OifName:%v,Prio:%v,IifName:%v,Invert:%v,mark:%v,goto:%v! rule:%v!",
		i,value.Table,value.Src,value.Dst,value.OifName,value.Priority,
		value.IifName,value.Invert,value.Mark,value.Goto,value)
		find := false
		for _,rule := range dstRules{
			if IsEqualRule(&value,rule) {
				find = true
				log.Infof("find rule: %v",rule)
				// dstRules = append(dstRules[:index],dstRules[index+1:]...)
				// 匹配多条
				break
			}
		}
		if find {
			// delete 
			total++
			log.Infof("need del rule :%v",value)
			NetAddorDelRule("del",&value)
		}
	}
	log.Infof("delete rule total: %v",total)
	return nil
}

// NetSyncPriorityRuleList 
func NetSyncPriorityRuleList(priority int,dstRules []*NT.Rule) error {
	// 
	rules, err := NT.RuleList(unix.AF_INET)
	if err != nil {
		log.Errorf("ListAllRule error:%v",err)
		return err
	}
	log.Infof("---list len=%d, dstrule len=%d!",len(rules),len(dstRules))
	// find this rule
	total := 0
	for i,value := range rules {
		log.Debugf("index:%v,table:%v,Src:%v,Dst:%v,OifName:%v,Prio:%v,IifName:%v,Invert:%v,mark:%v,goto:%v! rule:%v!",
		i,value.Table,value.Src,value.Dst,value.OifName,value.Priority,
		value.IifName,value.Invert,value.Mark,value.Goto,value)
		find := false
		for index,rule := range dstRules{
			if IsEqualRule(&value,rule) {
				find = true
				log.Infof("find rule: %v",rule)
				dstRules = append(dstRules[:index],dstRules[index+1:]...)
				break
			}
		}
		if priority == value.Priority && !find {
			// delete 
			total++
			log.Infof("need del rule :%v",value)
			NetAddorDelRule("del",&value)
		}
	}
	log.Infof("delete rule total: %v",total)
	total = 0
	for _,rule := range dstRules{
		log.Infof("need add rule :%v",rule)
		NetAddorDelRule("add",rule)
		total++
	}
	log.Infof("add rule total: %v",total)
	return nil
}

func ListAllRule() {
	rules, err := NT.RuleList(unix.AF_INET)
	if err != nil {
		log.Errorf("ListAllRule error:%v",err)
		return
	}
	log.Infof("---list rule len=%d!",len(rules))
	// find this rule
	for i := range rules {
		log.Infof("index:%v,table:%v,Src:%v,Dst:%v,OifName:%v,Prio:%v,IifName:%v,Invert:%v,mark:%v,goto:%v,mask:%v! rule:%v!",
		i,rules[i].Table,rules[i].Src,rules[i].Dst,rules[i].OifName,rules[i].Priority,
		rules[i].IifName,rules[i].Invert,rules[i].Mark,rules[i].Goto,rules[i].Mask,rules[i])
	}
}

// UpdateRouteTable 包含创建，删除，更新路由表，"add","del", add 是 tableid 不一致则更新最新的ID
func UpdateRouteTable(action string,tableName string, tableID int) error {
	// 判断是否存在
	const TABLE_FILE = "/etc/iproute2/rt_tables"
	tables,err := files.ReadFileAll(TABLE_FILE)
	if err != nil {
		log.Errorf("read routes get failed: %v",err)
		return err
	}
	tablesLine := strings.Split(tables,"\n")
	flag := false
	log.Infof("before update tablesLine:%v,len=%d",tablesLine,len(tablesLine))
	length := len(tablesLine)
	if tablesLine[length-1] == "" {
		tablesLine = tablesLine[:length-1]
	}
	for index,value := range tablesLine {
		if value == "" || ([]byte(value))[0] == '#' {
			continue
		}
		log.Debugf("get tables: %v",value)
		for index,v := range []byte(value) {
			log.Debugf("=====index:%d,%c!",index,v)
		}
		if strings.Index(value,tableName) > 0 {
			// 尝试 水平定位符号 分割
			data := []string{}
			data = strings.Split(value,"	")
			if len(data) != 2{
				// log.Errorf("tables route first error: %v,len=%v",data,len(data))
				data = strings.Split(value," ")
				if len(data) != 2{
					log.Errorf("tables route error: %v,len=%v",data,len(data))
					continue
				}
			}
			if data[1] == tableName {
				id, _ := strconv.Atoi(data[0])
				if action == "add" && id == tableID {
					flag = true
				}else if action == "add" && id != tableID && tableID > 0 {
					// 需要更新
					// 删除
					tablesLine = append(tablesLine[:index],tablesLine[index+1:]...)
				}else if action == "del" {
					// 删除
					tablesLine = append(tablesLine[:index],tablesLine[index+1:]...)
					flag = true
					break;
				}
			}
		}

	}
	log.Infof("after update tablesLine:%v,len=%d",tablesLine,len(tablesLine))
	if action == "add" && !flag {
		// 增加
		txt := fmt.Sprintf("%d %s",tableID,tableName)
		tablesLine = append(tablesLine,txt)
	}
	log.Infof("after over tablesLine:%v,len=%d",tablesLine,len(tablesLine))

	// 重新写入
	err = files.WriteListLineToFile(TABLE_FILE,tablesLine)	
	if err != nil {
		return err
	}
	return nil
}