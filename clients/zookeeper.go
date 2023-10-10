package clients

import (
	"fmt"
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

type ZooKeeperLock struct {
	conn      *zk.Conn
	lockPath  string
	lockIndex string
}

func NewZooKeeperLock(zkServers []string, lockPath string, nodeaddr string) (*ZooKeeperLock, error) {
	// 连接 ZooKeeper
	conn, _, err := zk.Connect(zkServers, time.Second*5)
	if err != nil {
		return nil, err
	}

	//创建一个/mylock目录
	acl := zk.WorldACL(zk.PermAll)
	_, err = conn.Create(lockPath, []byte{}, 0, acl)
	if err != nil && err != zk.ErrNodeExists {
		return nil, err
	}

	// 创建锁节点
	flags := int32(zk.FlagSequence)
	lockPath, err = conn.Create(lockPath+"/", nil, flags, acl)
	if err != nil {
		fmt.Println("err")
		return nil, err
	}
	_, err = conn.Create(lockPath+"/"+nodeaddr, nil, 0, acl)
	if err != nil {
		fmt.Println("err")
		return nil, err
	}

	// 获取锁节点编号
	lockIndex := lockPath[len(lockPath)-10:]
	fmt.Println(lockIndex, lockPath)
	return &ZooKeeperLock{conn, lockPath, lockIndex}, nil
}

func (l *ZooKeeperLock) Lock() error {
	// 尝试获取锁
	for {
		// 获取所有锁节点
		children, _, err := l.conn.Children(l.lockPath[:len(l.lockPath)-11])
		fmt.Println(children, l.lockPath[:len(l.lockPath)-11])
		if err != nil {
			return err
		}

		// 找到编号最小的锁节点
		minIndex := l.lockIndex
		for _, child := range children {
			if child < minIndex {
				minIndex = child
			}
		}

		// 如果当前节点是编号最小的锁节点，则获取锁成功
		if l.lockIndex == minIndex {
			fmt.Print("get lock success")
			return nil
		}

		// 否则等待一段时间后重试
		time.Sleep(time.Second)
	}
}

func (l *ZooKeeperLock) Unlock() {
	// 删除锁节点
	defer l.conn.Close()
	err := l.conn.Delete(l.lockPath, -1)
	if err != nil {
		log.Println(err)
	}
	return
}
