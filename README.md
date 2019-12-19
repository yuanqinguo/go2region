# go2region

基于纯真IP库的源数据，go实现的全内存的IP到省市区的转换服务

## API接口

### 通过IP查询省市区信息

**HTTP GET**  http://10.1.6.54:8080/go2region/ipinfo?ip=183.15.178.65

**返回**

```json
{
	"code": 100,
	"err_msg": "Ok",
	"data": {
		"province": "广东省",
		"city": "深圳市",
		"region": "南山区",
		"remark": "高新科技园xx网吧"
	}
}
```



### 热更新IP源数据

**HTTP POST** http://10.1.6.54:8080/go2region/reload

```json
{
	"code": 100,
	"err_msg": "Ok",
	"data": {
	}
}
```



## 配置运行

```json
项目依赖于NutsDB

make  # 得到执行文件go2region
./bin/go2region # -c {配置文件} -consul consul服务地址   -c和-consul二选一  -d=true表deamon方式运行，需要前台运行请不要携带-d参数
使用consul后台启动：
./bin/go2region -consul http://127.0.0.1:8500 -d=true
使用配置文件后台启动：
./bin/go2region -c go2region.yml -d=true

启动过程中需要加载源数据，加载时长根据纯真IP库导出来的数据量决定，虚拟机测试需要20分钟左右加载53W左右的IP端记录。
```



## IP源数据热更新

下载纯真IP地址数据库的客户端，启动后点击在线升级，升级完成后，点击解压，得到文件，将文件名改为yml配置文件中的文件名后，使用curl调用**热更新IP源数据的接口**，即可在后台进行源数据替换，加载过程比较长，故替换过程将在加载完后进行lock替换，未加载完成前均使用的是上一份源数据进行解析

```
 curl  http://10.1.6.54:8080/go2region/reload
```
## 限制
1. 准确度取决于纯真IP库的源数据
2. 运行机器内存需至少保留一半，否则在热更新数据时将内存不足


