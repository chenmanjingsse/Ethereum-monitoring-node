package discover

import (
	"fmt"
	"strings"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
)

// here
func get(name string) (string, string, string, string) {
	a := strings.Split(name, "@")[1]
	length := len(strings.Split(a, ":"))
	port := strings.Split(a, ":")[length-1]
	port = strings.Split(port, "?")[0]

	ip := ""
	for i := 0; i < length-1; i++ {
		ip += strings.Split(a, ":")[i]
	}

	ID := strings.Replace(ip, ".", "", -1)
	ID = strings.Replace(ID, ":", "", -1)
	ID = strings.Replace(ID, "[", "", -1)
	ID = strings.Replace(ID, "]", "", -1)
	ID = ID + port
	//fmt.Printf(ID+"\n")

	aa := strings.Split(name, "@")[0]
	lable := strings.Split(aa, "//")[1]
	return lable, ip, port, ID // ID是label，label是真正的节点ID
}

// inbound是接入连接数
func prepareCQL(ID string) string {
	//here
	CQL := "create (node" + ID + ":node" + ID + " { label: {label}, ifLive: {ifLive}, ip: {ip}, port: {port}, client: {client}, connected: {connected}})"
	return CQL
}

func createNode(n *node, ifLive bool, client string) bool {
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
	lable, ip, port, ID := get(n.String())
	// 以ip和id共同确定一个节点
	CQL := "match (node" + ID + ":node" + ID + ") where node" + ID + ".label='" + lable + "' and node" + ID + ".ip='" + ip + "' return node" + ID
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
		result, err2 := stmt.ExecNeo(map[string]interface{}{"label": lable, "ifLive": ifLive, "ip": ip, "port": port, "client": client, "connected": false})
		if err2 != nil {
			fmt.Printf("create error2")
			return false
		}
		_, err3 = result.RowsAffected()
		if err3 != nil {
			fmt.Printf("create error3")
			return false
		}
	} else { //match了的话，可能是match到之前没验证过的节点了，将其状态改为已验证
		stmt.Close()
		CQL = "match (node" + ID + ":node" + ID + ") set node" + ID + ".ifLive=true"
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
	//closeDriver(conn)
	return true
}

func CreateRalation(outn *node, inn *node, ifinLive bool, outclient string, inclient string) bool {
	if !createNode(outn, true, outclient) || !createNode(inn, ifinLive, inclient) {
		return false
	}

	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo("bolt://neo4j:admin@localhost:7687") //连接到服务器
	if err != nil {
		log.Info("connect error")
		//panic(err)
		return false
	}
	defer conn.Close()

	//conn := createDriver()

	outlabel, _, _, outID := get(outn.String())
	inlabel, _, _, inID := get(inn.String())
	CQL := "match (node" + outID + ":node" + outID + "),(node" + inID + ":node" + inID + "),p=(node" + outID + ")-[]-(node" + inID + ")  where node" + outID + ".label='" + outlabel + "' or node" + inID + ".label='" + inlabel + "' return p"
	stmt, err1 := conn.PrepareNeo(CQL)
	if err1 != nil {
		fmt.Printf("relation error1")
		return false
	}
	queryresult, err2 := stmt.QueryNeo(map[string]interface{}{})
	if err2 != nil || queryresult == nil {
		fmt.Printf("match error2")
		return false
	}
	inter, _, err3 := queryresult.All()
	if err3 != nil {
		fmt.Printf("match error3")
		return false
	}
	if len(inter) == 0 { //没有match的
		stmt.Close()
		CQL = "match (node" + outID + ":node" + outID + "),(node" + inID + ":node" + inID + ") where node" + outID + ".label='" + outlabel + "' or node" + inID + ".label='" + inlabel + "' create (node" + outID + ")-[r:have]->(node" + inID + ")"
		stmt, err1 := conn.PrepareNeo(CQL)
		if err1 != nil {
			fmt.Printf("relation error1")
			return false
		}
		_, err2 := stmt.ExecNeo(map[string]interface{}{})
		if err2 != nil {
			fmt.Printf("relation error")
			return false
		}
	}
	return true
}

// func main() {
// 	driver := bolt.NewDriver()
// 	conn, err := driver.OpenNeo("bolt://neo4j:admin@localhost:7687") //连接到服务器
// 	if err != nil {
// 		fmt.Printf("connect error\n")
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	stmt, err := conn.PrepareNeo("create (n:NODE)")
// 	if err != nil {
// 		fmt.Printf("create error\n")
// 		panic(err)
// 	}
// 	result, err := stmt.ExecNeo(map[string]interface{}{})
// 	if err != nil {
// 		fmt.Printf("exec error\n")
// 		panic(err)
// 	}
// 	numResult, err := result.RowsAffected()
// 	if err != nil {
// 		fmt.Printf("result error\n")
// 		panic(err)
// 	}
// 	fmt.Printf("Created rows: %d\n", numResult)
// }
