### job-schedule
使用quartz作为调度器，mysql持久化任务信息，提供简单易用的页面操作; 支持管理大量任务

### 1.快速启动
1. 下载项目
    ```shell
    git clone https://github.com/maolinc/job-schedule.git
    ```
2. 创建数据库表，sql语句位于model目录下
3. 更改etc/job.yaml配置文件的数据库连接地址DB.DataSource
4. 启动
   ```shell
   go run job.go
   ```
   访问http://127.0.0.1:8024；8023端口为内部rpc使用, 8024端口为api和页面使用；
   如果打包部署需要将etc/job.yaml中的Mode改为pro，否则访问不到页面

### 2.添加任务
在internal/jobaction/jobs目录下有两个例子LoginLogJob和SyncOrderJob。 编写任务大致流程：
1. 在internal/jobaction/jobs目录下编写任务，任务必须实现Plug接口，如下：
   ```go
   type Plug interface {
       // Action 定时任务需要实现此接口, res的字符数不能大于5000
       Action(ctx context.Context) (res any, err error)
       // Key 定时任务需要实现此接口, 返回的key与数据库的key一致, 每个插件的key必须不同,且不能为""
       Key() string
       // Compensate 补偿接口, res的字符数不能大于5000
       /**
       由于某种原因定任务某次未执行，可主动调用此方法进行补偿
       params执行参数
       */
       Compensate(ctx context.Context, params map[string]any) (res any, err error)
   }
   ```
   任务的具体逻辑不必放在本项目里面，可以使用rpc方式调用业务逻辑，可在internal/svc/serviceContext.go的ServiceContext添加rpc即可
2. 访问http://127.0.0.1:8024，点击创建任务，填写信息，key必须与步骤1中Key()结果一致
3. 任务创建后点击进去，可以管理任务的运行，以及主动执行补偿方法
![1685618027159](https://github.com/maolinc/job-schedule/assets/82015883/a26ae8af-23fd-45f2-a2d9-93784268c512)

![2](https://github.com/maolinc/job-schedule/assets/82015883/bb3e4576-c622-4e3b-8da6-f744566990cb)

### 3.项目使用到的框架和技术
1. 后端使用go-zero框架编写rpc和http (使用proto同时提供rpc和http服务)
2. 前端使用element-plus、axios、vue3 (页面资源位于static目录下)
3. 定时任务基础框架quartz
4. 数据库操作gorm

### 4.todo
1. 增加在页面http调度
2. 上传插件运行
3. 使用redis同步状态
4. ......
