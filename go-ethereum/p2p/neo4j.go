package p2p

import (
	"fmt"
	"strings"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
)

func getConnected(addr string) (string, string, string) {
	length := len(strings.Split(addr, ":"))
	port := strings.Split(addr, ":")[length-1]
	port = strings.Split(port, "?")[0]

	ip := ""
	for i := 0; i < length-1; i++ {
		ip += strings.Split(addr, ":")[i]
	}
	ID := strings.Replace(ip, ".", "", -1)
	ID = strings.Replace(ID, ":", "", -1)
	ID = strings.Replace(ID, "[", "", -1)
	ID = strings.Replace(ID, "]", "", -1)
	ID = ID + port

	return ip, port, ID
}

func prepareCQL(ID string) string {
	//here
	CQL := "create (node" + ID + ":node" + ID + " { label: {label}, ifLive: {ifLive}, ip: {ip}, port: {port}, client: {client}, connected: {connected}})"
	return CQL
}

func createConnectedNode(ID string, ifLive bool, client string, lable string, ip string, port string) bool {
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo("bolt://neo4j:admin@localhost:7687") //连接到服务器
	if err != nil {
		log.Info("connect error")
		//fmt.Printf("connect error\n")
		//panic(err)
		return false
	}
	defer conn.Close()

	// here
	//lable,ip,port,ID:=get(n.String())
	//CQL := "match (node" + ID + ":node" + ID + ") return node" + ID
	CQL := "match (node" + ID + ":node" + ID + ") where node" + ID + ".label='" + lable + "' return node" + ID
	stmt, err1 := conn.PrepareNeo(CQL)
	if err1 != nil {
		fmt.Printf("match error")
		return false
	}
	queryresult, err2 := stmt.QueryNeo(map[string]interface{}{})
	if err2 != nil || queryresult == nil {
		fmt.Printf(CQL + "\n")
		fmt.Printf("match error2\n")
		return false
	}
	inter, _, err3 := queryresult.All()
	if err3 != nil {
		fmt.Printf("match error3")
		return false
	}
	if len(inter) == 0 { //说明没有match的节点
		stmt.Close()
		//CQL = "create (node" + ID + ":node" + ID + ")"
		CQL = prepareCQL(ID)
		stmt, err1 = conn.PrepareNeo(CQL)
		if err1 != nil {
			fmt.Printf("create error1")
			return false
		}
		// here
		result, err2 := stmt.ExecNeo(map[string]interface{}{"label": lable, "ifLive": ifLive, "ip": ip, "port": port, "client": client, "connected": true})
		if err2 != nil {
			fmt.Printf("create error2")
			return false
		}
		_, err3 = result.RowsAffected()
		if err3 != nil {
			fmt.Printf("create error3")
			return false
		}
	} else { //如果有match到的节点
		stmt.Close()
		CQL = "match (node" + ID + ":node" + ID + ") set node" + ID + ".client=\"" + client + "\", node" + ID + ".connected=true"
		stmt, err1 = conn.PrepareNeo(CQL)
		if err1 != nil {
			fmt.Printf("change error1")
			return false
		}
		_, err2 := stmt.ExecNeo(map[string]interface{}{})
		if err2 != nil {
			fmt.Printf("change error2")
			return false
		}
	}
	return true
}
