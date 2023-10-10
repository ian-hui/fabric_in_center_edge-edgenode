package NodeUtils

import (
	"encoding/json"
	"fabric-edgenode/clients"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
)

type groupChooser struct {
	// TODO
	storage_weight            float64
	distance_weight           float64
	standard_deviation_weight float64
	// Nodestru Nodestructure
	curNodeInfo *NodeInfo
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

func (g *groupChooser) chooseGroup(k int) (Cur_group []string, delay float64, sd float64, err error) {
	var (
		nodelist              []NodeInfo
		candidate_nodeinfoMap = make(map[string]*processed_nodeinfo)
		mindist, maxdist      = math.MaxFloat64, 0.0
		minstor, maxstor      = math.MaxInt64, 0
		//创建一个group
	)
	Cur_group = make([]string, 0)

	//获取所有节点的距离以及剩余空间信息
	nodelist_byte, err := clients.GetPeerFabric(g.curNodeInfo.NodePeerName, "node").GetNodeInfoAllRange(g.curNodeInfo.NodePeerName)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetNodeInfoAllRange error: %v", err)
	}
	fmt.Println(string(nodelist_byte))
	//nodelist_byte 转化为 nodelist
	if err := json.Unmarshal(nodelist_byte, &nodelist); err != nil {
		return nil, 0, 0, fmt.Errorf("json.Unmarshal error: %v", err)
	}

	//计算本节点到其他节点的距离,顺便找出dist和storage的最大最小值
	for _, nodeif := range nodelist {

		dist := euclideanDistance(g.curNodeInfo, &nodeif)
		//判断是否是最大最小值
		if dist < mindist {
			mindist = dist
		}
		if dist > maxdist {
			maxdist = dist
		}
		//判断sorage的最大最小值
		nodeStorage, err := strconv.Atoi(nodeif.LeftStorage)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("strconv.Atoi error: %v", err)
		}
		//如果剩下的空间为0,则不计算
		if nodeStorage == 0 {
			continue
		}
		if nodeStorage < minstor {
			minstor = nodeStorage
		}
		if nodeStorage > maxstor {
			maxstor = nodeStorage
		}
		//候选节点
		candidate_nodeinfoMap[nodeif.NodePeerName] = &processed_nodeinfo{
			distance: dist,
			storage:  nodeStorage,
		}
	}
	//计算标准差
	sd = g.getStandardDeviation(nodelist)
	//归一化并计算得分
	for _, nodeif := range candidate_nodeinfoMap {
		dist, stor := normalized(nodeif.distance, float64(nodeif.storage), mindist, maxdist, minstor, maxstor)
		nodeif.score = g.getScore(dist, stor, sd)
	}
	//排序
	sortedMap := sortMap(candidate_nodeinfoMap)

	//如果k大于候选节点的数量,则返回错误(无法服务)
	if len(sortedMap) < k {
		return nil, 0, 0, fmt.Errorf("can not find enough node to form a group")
	}

	for i := 0; i < k; i++ {
		Cur_group = append(Cur_group, sortedMap[i].Key)
		//选择这些节点
		delay += candidate_nodeinfoMap[sortedMap[i].Key].distance * 100
		//减少这些节点的剩余空间
		b, err := clients.GetPeerFabric(g.curNodeInfo.NodePeerName, "node").GetNodeInfo(sortedMap[i].Key, g.curNodeInfo.NodePeerName)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("GetNodeInfo error: %v", err)
		}
		var node NodeInfo
		if err := json.Unmarshal(b, &node); err != nil {
			return nil, 0, 0, fmt.Errorf("json.Unmarshal error: %v", err)
		}
		node.LeftStorage = strconv.Itoa(candidate_nodeinfoMap[sortedMap[i].Key].storage - 1)
		res, err := json.Marshal(node)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("json.Marshal error: %v", err)
		}
		if s, err := clients.GetPeerFabric(g.curNodeInfo.NodePeerName, "node").SetNodeInfo(node.NodePeerName, res); err != nil {
			return nil, 0, 0, fmt.Errorf("SetNodeInfo error: %v", err)
		} else {
			fmt.Println(s)
		}
	}
	//平均延迟
	delay /= float64(k)
	return
}

func euclideanDistance(n1, n2 *NodeInfo) float64 {
	n1x, err := strconv.ParseFloat(n1.LocationX, 64)
	if err != nil {
		fmt.Println(err)
	}
	n1y, err := strconv.ParseFloat(n1.LocationY, 64)
	if err != nil {
		fmt.Println(err)
	}
	n2x, err := strconv.ParseFloat(n2.LocationX, 64)
	if err != nil {
		fmt.Println(err)
	}
	n2y, err := strconv.ParseFloat(n2.LocationY, 64)
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

func (g *groupChooser) getStandardDeviation(nodelist []NodeInfo) float64 {

	sum := 0
	for _, nodeif := range nodelist {
		storage, err := strconv.Atoi(nodeif.LeftStorage)
		if err != nil {
			panic(err)
		}
		sum += storage
	}
	average := float64(sum) / float64(len(nodelist))

	//计算标准差
	var sd float64
	for _, nodeif := range nodelist {
		storage, err := strconv.Atoi(nodeif.LeftStorage)
		if err != nil {
			panic(err)
		}
		sd += math.Pow(float64(storage)-average, 2)
	}

	return sd
}

func (g *groupChooser) getScore(dist float64, storage float64, sd float64) float64 {
	return g.distance_weight*dist + g.storage_weight*storage + g.standard_deviation_weight*sd
}

func (g *groupChooser) randomChooser(k int) (Cur_group []string, delay float64, sd float64, err error) {
	var nodelist []NodeInfo
	nodelist_byte, err := clients.GetPeerFabric(g.curNodeInfo.NodePeerName, "node").GetNodeInfoAllRange(g.curNodeInfo.NodePeerName)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("GetNodeInfoAllRange error: %v", err)
	}
	if err := json.Unmarshal(nodelist_byte, &nodelist); err != nil {
		return nil, 0, 0, fmt.Errorf("json.Unmarshal error: %v", err)
	}
	//获取所有节点的距离以及剩余空间信息
	candidate_nodeinfoMap := make(map[string]*processed_nodeinfo)
	//计算本节点到其他节点的距离,顺便找出dist和storage的最大最小值
	for _, nodeif := range nodelist {
		//判断sorage的最大最小值
		nodeStorage, err := strconv.Atoi(nodeif.LeftStorage)
		if err != nil {
			panic(err)
		}
		//如果剩下的空间为0,则不计算
		if nodeStorage == 0 {
			continue
		}
		//如果本节点的剩余空间不为0
		// if g.curNodeInfo.LeftStorage != "0" {
		//候选节点距离以本节点为中心
		candidate_nodeinfoMap[nodeif.NodePeerName] = &processed_nodeinfo{
			storage:  nodeStorage,
			distance: euclideanDistance(g.curNodeInfo, &nodeif),
		}
		// }
		//  else {
		// 	//候选节点距离以

	}
	if len(candidate_nodeinfoMap) < k {
		return nil, 0, 0, fmt.Errorf("can not find enough node to form a group")
	}
	//计算标准差
	sd = g.getStandardDeviation(nodelist)
	//在candidate中随机挑选k个
	all_candidate := make([]string, 0)
	for k := range candidate_nodeinfoMap {
		all_candidate = append(all_candidate, k)
	}
	for i := 0; i < k; i++ {
		//随机选择一个节点
		rand_nodeid := rand.Intn(len(all_candidate))
		//选择这些节点
		delay += candidate_nodeinfoMap[all_candidate[rand_nodeid]].distance * 100
		//减少这些节点的剩余空间
		num, err := strconv.Atoi(all_candidate[rand_nodeid])
		if err != nil {
			panic(err)
		}
		nodelist[num-1].LeftStorage = strconv.Itoa(candidate_nodeinfoMap[all_candidate[rand_nodeid]].storage - 1)
		//取出这个节点
		all_candidate = append(all_candidate[:rand_nodeid], all_candidate[rand_nodeid+1:]...)
	}
	return
}

func NewGroupChooser(curNodeInfo *NodeInfo, storage_weight, distance_weight, standard_deviation_weight float64) *groupChooser {
	return &groupChooser{
		curNodeInfo:               curNodeInfo,
		storage_weight:            storage_weight,
		distance_weight:           distance_weight,
		standard_deviation_weight: standard_deviation_weight,
	}
}

func (gc *groupChooser) ChooseGroupWithLock(nodestru Nodestructure) (group []string, delay float64, sd float64, err error) {
	//获取分布式锁
	zkl, err := clients.NewZooKeeperLock([]string{nodestru.ZookeeperAddr}, "/mylock", nodestru.KafkaIp)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("NewZooKeeperLock error:%v", err)
	}
	if err := zkl.Lock(); err != nil {
		return nil, 0, 0, fmt.Errorf("Lock error:%v", err)
	}
	//选择group
	group, delay, sd, err = gc.chooseGroup(5)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("chooseGroup error:%v", err)
	}
	//释放分布式锁
	defer zkl.Unlock()
	return
}
