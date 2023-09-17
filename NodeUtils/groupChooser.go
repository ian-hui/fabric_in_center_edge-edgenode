package NodeUtils

import (
	"encoding/json"
	"fabric-edgenode/clients"
	"fmt"
	"math"
	"sort"
	"strconv"
)

type groupChooser struct {
	// TODO
	weight   float64
	Nodestru Nodestructure
}

type processed_nodeinfo struct {
	distance float64
	storage  int
	score    float64
}

// 将 map 的键和值复制到一个切片中
type kv struct {
	Key   string
	Value float64
}

// func Constructor

func (g *groupChooser) chooseGroup(k int) (group []string, err error) {
	//获取所有节点的距离以及剩余空间信息
	var v []NodeInfo
	nodeinfoMap := make(map[string]*processed_nodeinfo)
	res, err := clients.GetPeerFabric(g.Nodestru.PeerNodeName, "node").GetNodeInfoAllRange(g.Nodestru.PeerNodeName)
	if err != nil {
		return
	}
	if json.Unmarshal(res, &v) != nil {
		return
	}
	mindist, maxdist := math.MaxFloat64, 0.0
	minstor, maxstor := math.MaxInt64, 0
	//计算本节点到其他节点的距离,顺便找出dist和storage的最大最小值
	for _, nodeif := range v {
		//找出本节点,不计算
		if nodeif.NodeIp == g.Nodestru.KafkaIp {
			continue
		}
		dist := euclideanDistance(g.Nodestru.NodeInfo, &nodeif)
		//判断是否是最大最小值
		if dist < mindist {
			mindist = dist
		}
		if dist > maxdist {
			maxdist = dist
		}
		//判断storage的最大最小值
		nodeStorage, err := strconv.Atoi(nodeif.LeftStorage)
		if err != nil {
			return group, err
		}
		if nodeStorage < minstor {
			minstor = nodeStorage
		}
		if nodeStorage > maxstor {
			maxstor = nodeStorage
		}
		nodeinfoMap[nodeif.NodeIp] = &processed_nodeinfo{
			distance: dist,
			storage:  nodeStorage,
		}
	}
	//归一化并计算得分
	for _, nodeif := range v {
		if nodeif.NodeIp == g.Nodestru.KafkaIp {
			continue
		}
		dist, stor := normalized(nodeinfoMap[nodeif.NodeIp].distance, float64(nodeinfoMap[nodeif.NodeIp].storage), mindist, maxdist, minstor, maxstor)
		nodeinfoMap[nodeif.NodeIp].score = g.weight*dist + (1-g.weight)*stor
	}
	//排序
	sortedMap := sortMap(nodeinfoMap)
	for i := 0; i < k; i++ {
		group = append(group, sortedMap[i].Key)
	}
	return
}

func euclideanDistance(n1, n2 *NodeInfo) float64 {
	n1x, err := strconv.ParseFloat(n1.locationX, 64)
	if err != nil {
		fmt.Println(err)
	}
	n1y, err := strconv.ParseFloat(n1.locationY, 64)
	if err != nil {
		fmt.Println(err)
	}
	n2x, err := strconv.ParseFloat(n2.locationX, 64)
	if err != nil {
		fmt.Println(err)
	}
	n2y, err := strconv.ParseFloat(n2.locationY, 64)
	if err != nil {
		fmt.Println(err)
	}
	return math.Sqrt(math.Pow(n1x-n2x, 2) + math.Pow(n1y-n2y, 2))
}

// 距离，剩余空间归一化
func normalized(dist, stor, mindist, maxdist float64, minstor, maxstor int) (float64, float64) {
	dist = (dist - mindist) / (maxdist - mindist)
	stor = float64(stor-float64(minstor)) / float64(maxstor-minstor)
	return dist, stor
}

// 排序
func sortMap(m map[string]*processed_nodeinfo) []kv {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v.score})
	}

	// 使用 sort 包中的函数对切片进行排序
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	return ss
}
