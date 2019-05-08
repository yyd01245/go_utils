package firewall


import (
	"strings"
	// "fmt"
	// "reflect"
	log "github.com/Sirupsen/logrus"	
	GPT "github.com/coreos/go-iptables/iptables"
)



const MARKTABLE = "mangle"
const FIREWALLTABLE = "filter"
const NATTABLE = "nat"

// const FIREWALLCHAIN = "FORWARD"
const FIREWALLCHAIN = "USER-ACL"


func AddGroupChain(chain string,rule []string) error {
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}

	err = ipt.Append(FIREWALLTABLE,FIREWALLCHAIN, rule...)
	if err != nil {
		log.Warnf("Append failed: %v", err)
	}

	return nil
}

func ClearGroupChain(chain string) error {
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}

	// check is exist
	rules, err := ipt.List(FIREWALLTABLE,FIREWALLCHAIN)
	if err != nil {
		log.Warnf("List failed: %v", err)
		return err
	}

	expected := false
	for _,value := range rules {
		if value == ("-N " + FIREWALLCHAIN) {
			log.Debugf("exist chain, match over")
			expected = true
			break;
		}
	}
	if expected  {
		ipt.ClearChain(FIREWALLTABLE,FIREWALLCHAIN)
	}else {
		// create chain
		// chain shouldn't exist
		err = ipt.NewChain(FIREWALLTABLE, FIREWALLCHAIN)
		if err != nil {
			log.Warnf("new chain error: %v",err)
			return err
		}	
		log.Debugf("NewChain success")	
		err = ipt.Append(FIREWALLTABLE, "FORWARD",  "-j",FIREWALLCHAIN )
		if err != nil {
			log.Warnf("append chain failed: %v", err)
			return err
		}	
		// iptables -t filter -A FORWARD -j $USER_CHAIN	
		// txt := fmt.Sprintf("no chain ,do not need clear, rule :%v", rules)
		// log.Warnf(txt)
	}

	return nil
}

func ClearChain(chainTable string,chain string) error{
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}
	ipt.ClearChain(chainTable,chain)
	return nil
}

func AddChain(chainTable string,chain string, rule []string) error{
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}

	err = ipt.Append(chainTable,chain, rule...)
	if err != nil {
		log.Warnf("Append failed: %v", err)
	}

	return nil
}

func InsertChain(chainTable string,chain string, pos int, rule []string) error{
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}

	err = ipt.Insert(chainTable,chain,pos,rule...)
	if err != nil {
		log.Warnf("Append failed: %v", err)
	}

	return nil
}

func DeleteChain(chainTable string,chain string, rule []string) error{
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}

	err = ipt.Delete(chainTable,chain, rule...)
	if err != nil {
		log.Warnf("Delete failed: %v", err)
	}

	return nil
}

func ListChain(chainTable string,chain string, rule string) (bool,error){
	ipt, err := GPT.New()
	ret := false
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return ret,err
	}
	// check list
	rules, err := ipt.List(chainTable,chain)
	if err != nil {
		log.Warnf("List failed: %v", err)
		return ret,err
	}

	for _,value := range rules {
		if strings.Index(value,rule) >= 0 {
			log.Debugf("exist chain, match over")
			ret = true
			break;
		}
	}

	return ret,nil
}
func ListChainAll(chainTable string,chain string) ([]string,error){
	ipt, err := GPT.New()

	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return []string{},err
	}
	// check list
	rules, err := ipt.List(chainTable,chain)
	if err != nil {
		log.Warnf("List failed: %v", err)
		return rules,err
	}

	return rules,nil
}

func DeleteChainbyString(chainTable string,chain string, rule string) error{
	ipt, err := GPT.New()
	if err != nil {
		// panic(fmt.Sprintf("New failed: %v", err))
		log.Warnf("error create ipt instance:%v ",err)
		return err
	}

	err = ipt.Delete(chainTable,chain, rule)
	if err != nil {
		log.Warnf("Delete failed: %v", err)
	}

	return nil
}
	