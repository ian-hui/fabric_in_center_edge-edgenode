package NodeUtils

//services used to proccess the web request
import (
	"encoding/json"
	"fabric-edgenode/sdkInit"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 10.0.0.144:8083/register
// 可选择发给哪个kafka
// POST Username:
//
//	    Attribute:
//		   kafkaIp:
//
// 注册身份
func Register(ctx *gin.Context) {
	fmt.Println("<--------service register--------->")
	//提取信息
	kafkaIp := ctx.Query("kafkaIp")
	username := ctx.Query("Username")
	attribute := ctx.QueryArray("Attribute")
	//生成一个uuid作为userid
	// uuid := uuid.New()
	// userid := uuid.String()
	//构建一个user信息结构
	userinformation := sdkInit.UserInfo{
		UserId:    username,
		Username:  username,
		Attribute: attribute,
	}
	res, err := json.Marshal(userinformation)
	if err != nil {
		fmt.Printf("fail to Serialization, err:%v\n", err)
		return
	}
	topic := "register" //操作名
	err = ProducerAsyncSending(res, topic, kafkaIp)
	if err != nil {
		ctx.JSON(403, gin.H{
			"message": err,
		})
	}
	ctx.JSON(200, gin.H{
		"UserID": username,
	})
}

// TODO 返回密文id
// 10.0.0.144:8083/upload
// 可选择发给哪个kafka
// POST	file:
// 把密文上传到节点中，同时节点把密文信息和位置上传到中心节点
func Upload(ctx *gin.Context) {
	fmt.Println("<--------service file upload--------->")
	//提取信息
	fileid := ctx.Query("fileid")
	kafkaIp := ctx.Query("kafkaIp")
	// file := ctx.Query("file")
	//生成一个uuid作为密文id
	// uuid := uuid.New()
	// fileid := uuid.String()
	const kb = 1024
	file := strings.Repeat("a", 512*kb) // 5MB string
	//构建一个user信息结构
	fileinfomation := FileInfo{
		FileId:     fileid,
		Ciphertext: file,
	}

	// 连接kafka
	res, err := json.Marshal(fileinfomation)
	if err != nil {
		fmt.Printf("fail to Serialization, err:%v\n", err)
		return
	}
	fmt.Println(len(res))
	topic := "upload" //操作名
	err = ProducerAsyncSending(res, topic, kafkaIp)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err,
		})
	}
	ctx.JSON(200, gin.H{
		"FileID": fileid,
	})
}

// filerequest ：user request for file
// 10.0.0.144:8083/upload
// POST
func Filerequest(ctx *gin.Context) {
	fmt.Println("<--------service request file--------->")
	//提取信息
	kafkaIp := ctx.Query("kafkaIp")
	FileId := ctx.Query("fileId")
	UserId := ctx.Query("userId")
	FilerequestStruct := FileRequest{
		FileId: FileId,
		UserId: UserId,
	}
	res, err := json.Marshal(FilerequestStruct)
	if err != nil {
		fmt.Printf("fail to Serialization, err:%v\n", err)
		return
	}
	topic := "filereq" //操作名
	err = ProducerAsyncSending(res, topic, kafkaIp)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err,
		})
	}
	ctx.JSON(200, gin.H{
		"message": "filerequest send success",
	})
}
