package controls

import (
	"bufio"
	"fmt"
	"github.com/xujiajun/nutsdb"
	"github.com/xujiajun/nutsdb/ds/zset"
	"go2region/utils"
	"go2region/utils/helper"
	"go2region/utils/logs"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	once         sync.Once
	instance     *IpDataInfo
	bucketName   = "iparea:bucket"
	specialCity  = [...]string{"上海市", "北京市", "重庆市", "天津市"}
	specialPlace = [...]string{"香港", "澳门"}
	specialProv  = [...]string{"内蒙古", "新疆", "西藏", "宁夏", "广西"}
)

type IpDataInfo struct {
	nutsDB     *nutsdb.DB
	ipdatakey  string
	filepath   string
	dbGuard    sync.Mutex
	bucketName string
}

func GetInstance() *IpDataInfo {
	once.Do(func() {
		instance = &IpDataInfo{}
		fmt.Println("IpDataInfo instance...")
	})
	return instance
}

func (ipdb *IpDataInfo) loaderFromFile() (*nutsdb.DB, error) {
	fmt.Println("Start IpDataInfo loaderFromFile...")
	f, err := os.Open(ipdb.filepath)
	if err != nil {
		panic(err)
		return nil, utils.ServerError
	}
	defer f.Close()

	dbpath := "/tmp/nutsdb/go2region"
	files, _ := ioutil.ReadDir(dbpath)
	for _, f := range files {
		name := f.Name()
		if name != "" {
			fmt.Println(dbpath + "/" + name)
			err := os.RemoveAll(dbpath + "/" + name)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	opt := nutsdb.DefaultOptions
	opt.Dir = dbpath //这边数据库会自动创建这个目录文件
	opt.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	db, err := nutsdb.Open(opt)
	if err != nil {
		return nil, utils.ServerError
	}

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		ipdb.parseLine(db, helper.ConvertStr2GBK(line))
	}
	fmt.Println("End IpDataInfo loaderFromFile...")

	ipdb.dbGuard.Lock()
	oldDb := ipdb.nutsDB
	ipdb.nutsDB = db
	ipdb.dbGuard.Unlock()

	return oldDb, nil
}

func (ipdb *IpDataInfo) parseLine(db *nutsdb.DB, line string) {
	line = strings.Replace(line, "\r\n", "", -1)
	line = helper.DeleteExtraSpace(line)
	strArr := strings.Split(line, " ")
	if len(strArr) < 3 {
		return
	}

	remark := ""
	start := strArr[0]
	end := strArr[1]
	value := strArr[2]
	if len(strArr) > 3 {
		remark = strArr[3]
	}
	if helper.CheckIp(start) && helper.CheckIp(end) {
		value = ipdb.handleValue(value, remark)
		if len(value) > 0 {
			ipdb.saveOne(db, start, end, value)
		} else {
			fmt.Println(line)
			fmt.Println(value)
		}
	}
}

func (ipdb *IpDataInfo) handleValue(value string, ext string) string {
	var prov, city, region, remark string

	vArr := strings.Split(value, " ")
	mask := vArr[0]
	remark = ext

	isSpCity := false
	isSpProv := false
	isSpPlace := false
	for _, sp := range specialCity {
		if strings.Contains(mask, sp) {
			isSpCity = true
			// 特殊城市
			index := strings.Index(mask, "市")
			if index > 0 {
				prov = mask[:index+3]
				city = prov
				region = ""
			}
			break
		}
	}

	if !isSpCity {
		for _, sp := range specialPlace {
			if strings.Contains(mask, sp) {
				isSpPlace = true
				// 澳门香港特殊处理
				prov = mask
				city = mask
				region = ""
				break
			}
		}
	}

	if !isSpCity && !isSpPlace {
		for _, sp := range specialProv {
			if strings.Contains(mask, sp) {
				isSpProv = true
				// 特殊省份，不是以省命名
				index := strings.Index(mask, sp)
				if index >= 0 {
					prov = mask[:index+len(sp)]
					city = mask[index+len(sp):]
					region = ""
				}
				break
			}
		}
	}
	// 省xx市  省xx州  省xx地区 省xx县
	if !isSpProv && !isSpCity && !isSpPlace {
		index := strings.Index(mask, "省")
		if index > 0 {
			prov = mask[:index+3]
			cityIndex := strings.Index(mask, "市")
			if cityIndex >= 0 {
				city = mask[index+3 : cityIndex+3]
				region = mask[cityIndex+3:]
			} else {
				cityIndex := strings.LastIndex(mask, "州")
				if cityIndex >= index+3 {
					city = mask[index+3 : cityIndex+3]
					region = mask[cityIndex+3:]
				} else {
					cityIndex := strings.Index(mask, "地区")
					if cityIndex >= index+3 {
						city = mask[index+3 : cityIndex+6]
						region = mask[cityIndex+6:]
					} else {
						// 有些数据直接省后面接了区县
						regionIndex := strings.Index(mask, "县")
						if regionIndex >= index+3 {
							region = mask[index+3 : regionIndex+3]
						}
					}
				}
			}
		}
	}

	// 无省份或者无市的数据丢弃
	if len(prov) == 0 || (len(city) == 0 && len(region) == 0) {
		return ""
	}

	// 类似: 广东省,深圳市,南山区,星巴克咖啡厅
	return fmt.Sprintf("%s,%s,%s,%s", prov, city, region, remark)
}

func (ipdb *IpDataInfo) saveOne(db *nutsdb.DB, start string, end string, value string) {
	ipstart := helper.ZoomInTenThousand(helper.ConverIptoInt(start))
	ipend := helper.ZoomInTenThousand(helper.ConverIptoInt(end)) + 1

	zaddStartFunc := func(tx *nutsdb.Tx) error {
		bucket := bucketName
		key := []byte(strconv.Itoa(int(ipstart)))
		return tx.ZAdd(bucket, key, float64(ipstart), []byte("start-"+value))
	}

	zaddEndFunc := func(tx *nutsdb.Tx) error {
		bucket := bucketName
		key := []byte(strconv.Itoa(int(ipend)))
		return tx.ZAdd(bucket, key, float64(ipend), []byte("end-"+value))
	}

	if err := db.Update(zaddStartFunc); err != nil {
		logs.LogSystem.Error(err)
	}

	if err := db.Update(zaddEndFunc); err != nil {
		logs.LogSystem.Error(err)
	}

}

func (ipdb *IpDataInfo) Close() {
	if ipdb.nutsDB != nil {
		_ = ipdb.nutsDB.Close()
	}
}

func (ipdb *IpDataInfo) Reloader(path string) {

	ipdb.filepath = path

	oldIPDB, err := ipdb.loaderFromFile()

	if oldIPDB != nil && err == nil {
		_ = oldIPDB.Close()
		fmt.Println("oldIPDB close")
	}
}

func (ipdb *IpDataInfo) GetIpInfo(ip int64) (province, city, region, remark string) {
	ip = helper.ZoomInTenThousand(ip)
	db := ipdb.nutsDB
	limitOps := &zset.GetByScoreRangeOptions{Limit: 1}

	zrangeByscoreFunc := func(tx *nutsdb.Tx) error {
		bucket := bucketName
		if nodes, err := tx.ZRangeByScore(bucket, float64(ip), 10000000000000000, limitOps); err != nil {
			return err
		} else {
			if nodes != nil && len(nodes) > 0 {
				value := string(nodes[0].Value)

				arrs := strings.Split(value, "-")
				if len(arrs) > 1 {
					dstValue := arrs[1]
					dstArrs := strings.Split(dstValue, ",")
					province = dstArrs[0]
					city = dstArrs[1]
					region = dstArrs[2]
					remark = dstArrs[3]
				}
			}
		}
		return nil
	}

	if err := db.View(zrangeByscoreFunc); err != nil {
		logs.LogSystem.Error(err)
	}
	return
}
