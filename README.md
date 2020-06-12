-----

#### Metadata：基于Excel的配置系统



##### 01 为什么要使用Excel做配置，而不使用Apollo, MySQL这样的系统

1. 有些复杂的业务配置，比如游戏，配置的量非常大，并且经常批量刷新修改，这时使用Apollo这种一条数据一条数据修改的方式就很耗时了。
2. 对产品/策划而言，Excel反而是更加熟悉的工具。



##### 02 支持哪些类型？

主要支持两种类型的配置：Template和Config

1. Template是指同一种类型的配置，有多行不同的配置数据，它们之间以id区分。比如商品参数，大量商品的的参数结构是一样的，但是具体的参数又不行。它们在Excel中用不同的行表示
2. Config是指针对这个类型的配置，全局只有一份。比如全局参数配置。它在Excel中只占用一行



```go

type TestTemplate struct {
	Id      int         `xlsx:"column(id)"`          // 按列映射；支持整数；
	Name    string      `xlsx:"name"`                // 支持中文；column()可以省略
	NamePtr *string     `xlsx:"name"`                // 同一列可以映射到多个字段
	Height  float32     `xlsx:"height;default(1.2)"` // 支持浮点数；如果不填，默认值为1.2
	Titles  []string    `xlsx:"titles;split(|)"`     // 支持slice，可以使用使用分隔符，比如空格 " "
	Person  *TestPerson `xlsx:"person"`              // 通过实现UnmarshalBinary接口，可以支持嵌入json字符串；但这里加default({\"Name\":\"Panda\", \"Age\":18}) 之后好像就报错了
}

func main() {
  var url = "xxxx"
  metadata.AddExcel(url)
  
  .....
  var template TestTemplate
	var err = metadata.GetTemplate(1, &template)
  if err != nil{
    .....
  }
}
```



具体参考项目中的测试代码



##### 03 如何在线更新配置？

系统从Excel文件中加载配置，支持本地文件路径和url路径。

如果需要支持在线更新，则建议方式为：将Excel文件上传到某一个网站地址，配置更新时，直接上传覆盖即可。系统会每分钟检查Excel文件是否更新，如果有更新，则会自动下载新的Excel文件。

