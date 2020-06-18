# **只夕——珍惜时间** 后端

​    开发者：DimDoremy       2020-6-15          

# 概述

​    该项目的后端使用的Golang语言进行的开发，具体使用情况如下：

1.  网络框架使用Gin框架，所有前后端交互使用POST请求实现。

2.  数据库方面使用了MySQL作为存储层，GORM模型来处理与数据库的操作。

3.  使用NoSQL的Redis作为服务器缓存池建立缓存层，通过官方内置的reflect包实现缓存数据的解析。

4.  服务器部署使用Nginx，通过反向代理的方式运行服务器程序。

5.  运维方面，API文档的设计使用swaggo包通过注解方式实现API文档的部署并可以让所有参与工作的人员查阅和测试API端口请求。

6.  迭代方面使用AIR包，做到热更新。

# 具体细节

## Gin框架部分

​    Gin框架的处理使用了路由组的方式进行管理，由于Golang语言是一个面向过程的语言，所以在编写时通过函数将每个需要的路由进行封装。

​    而Gin的Group操作，可以将零散的路由通过一个公共前缀组合起来，方便集中管理（详见图表2）。

 

![图表 1 main入口函数中的路由组注册](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image002.jpg)

![图表 2管理路由组下的所有路由情况（这里以UserRouter为例）](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image004.jpg)



## JSON请求部分

对于使用JSON请求，我封装了一个函数专门用于绑定JSON请求使用的请求体，具体操作如下：

1. 获取gin框架下Context结构体指针context用户进行绑定。

2. 通过输入的bindStruct接口作为JSON的容器绑定到gin中。

3. 绑定成功执行输入参数的函数部分，将匿名函数作为参数传递。

 

![图表 3用于绑定JSON的函数](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image006.jpg)



## MySQL连接部分

​    通过ORM模型，使用GORM进行数据库连接，用于存放用户账号以及账号个人信息等需要长期保存的数据，具体细节如下：

1. 定义全局gorm.DB结构体指针，用于入口文件初始化调用。

2. 连接MySQL数据库并使用Ping()的error类型返回判断是否连接成功。

3. 定义一个单独用来关闭的函数，方便入口函数可以在运行结束时gc。

 

![图表 4 MySQL数据库连接相关内容](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image008.jpg)

![图表 5入口函数加载数据库连接](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image010.jpg)



## MySQL数据库操作

​    因为数据库操作的模型是ORM模型（对象关系映射），所以在这里我们项目通过定义映射使用的结构体实现各种操作数据库的接口让客户端请求调用：

1. 为了可以让定义的结构体对象可以映射到指定数据库的属性上，我们通过使用gorm标签定义的方式实现属性值与结构体参数的一一对应。

2. 因为需要接收来自客户端的请求用来初始化，所以同样在结构体上使用json标签将结构体与JSON对象对应。

3. 具体内部增删改查的实现参考：https://gorm.io/zh_CN/docs/。



![图表 6用于ORM的结构体（以UserData为例）](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image011.jpg)

## Redis连接池

​    使用Redis连接池搭建用于用户招募操作的缓存池，因为招募操作需要数据的频繁更新、创建和删除，所以这里使用效率较高的Redis连接池，通过空间换取时间让响应速度变快，具体细节如下：

1. 建立全局redis.Pool的Pool指针，指向一个连接池空间。

2. 通过构造函数初始化Pool结构体，划分一个8个长连接对象初始化的连接池，限制同时使用结构体的连接数为100000并在连接数量超过上限时阻塞。

3. 在创建连接池方法执行结束后会执行测试连接池方法，同样通过Ping测试连接是否成功。

 

![图表 7 初始化连接池](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image013.jpg)

![图表 8入口函数启动连接池和gc](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image015.jpg)

## Redis数据库操作

​    对于缓存Redis的操作，由于Redis的存储结构是以Key-Value的方式进行的存储，所以具体的操作方式我划分了三种：

①  通过Key找到Value。

②  修改对应Key的Value，找不到Key则执行添加这组Key-Value。

③  删除对应Key-Value。

故将Redis操作拆分成三个函数，执行上述的三个功能。

 

![图表 9通过Key找到Value](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image017.jpg)

![图表 10修改对应Key的Value，找不到Key则执行添加这组Key-Value](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image019.jpg)

![图表 11删除对应Key-Value](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image020.jpg)

## 微信登陆服务器验证

​    根据微信小程序的开发文档描述，对于有关微信授权的操作需要先由客户端请求我们使用的后端服务器，再通过后端服务器发送请求给微信小程序的后台，获取Session_key和openid，故在此编写一个用于小程序登录的接口。

![图表 12微信小程序服务器验证请求接口](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image022.jpg)

## Swaggo注解API文档开发

​    Swaggo工具是Java中Swagger的Golang实现，但是不同于Java可以直接使用注解方式配置Swagger，Golang里没有和Java一样“@“的注解，所以这个工具是通过将注释对应解析，生成docs.go文件解析swagger.json实现的如同Swagger的操作。

![图表 13 swagger.json](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image024.jpg)

![图表 14 docs.go](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image026.jpg)

## AIR

​    通过配置AIR包，实现每次重新部署后的热更新。

​    执行程序时，运行air -c .air.conf来监控当前目录对应文件是否改变。

​    具体操作参考：https://github.com/cosmtrek/air

![图表 15 .air.conf](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image028.jpg)

## 反射（Reflect）

​    因为Redis连接使用的redigo包为了读取不同类型的数据，将获取参数的返回值定义为interface{}型，而Golang语言一个很大的问题在于为了效率牺牲了泛型，interface{}这个类型，仅可以储存任意类型的数据但不能弱类型转换成其他类型数据，对于这个问题，我们采用了官方包中的Reflect包获取对应值转换为结构体进行解析。

![图表 16 Reflect实现错误处理和取值](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image030.jpg)

# 现有的问题

## 性能与Redis阻塞

问题：由于Redis数据库连接池使用的阻塞方式运行，会存在当同时请求过多的情况下发生阻塞导致整体的请求时间过长。

这里是我用Golang的hey工具进行压力测试有关招募的接口时的测试报告。

测试规格：1000用户2000并发。

测试接口：/return_recruit

测试情况：

可以看出总共2000次的请求下，出现了214个400的返回值和92个Error请求。92个请求错误原因是因为TCP连接的套接字在Windows里有上限，这里超过了系统限制发生的错误；而214个400的请求绝大多数源于请求被占用。而根据这个测试报告可以看出，1694次的成功里，单1.768s和3.535s就一共1145次成功，成功率为67.59%，在这个相应时间内，可以保证小量用户使用过程中不会太受影响。但是其他阻塞的请求最多的甚至有请求时间达到17.666s之久，所以对于当前服务器后端的性能还亟待优化。

预期解决方式：负载均衡。

![图表 17 hey工具POST测试Redis连接池测试报告](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image032.jpg)

## 招募获取房间人数

问题：项目功能实现之前会议中提到有关招募操作，建立房间应该能看到当前房间人数，这一点需要由服务器端主动发送请求，但是现在服务器后端无法做到主动给客户端发送请求。

预期解决方式：通过添加对Redis的监听，实现当Redis的数据修改时读取两个Redis存储的字节长度判断当前是添加新用户还是老用户离开，并将结果通过http/2.0或者消息队列广播到所有当前房间内的用户。

## 招募的结束判定问题

问题：关于招募结算，我们的实现方式是通过客户端获取Redis缓存与当前时间比较的方式判断是否招募结束，如果结束则结算奖励并发放，用户离开当前房间，房间所有人全部离开时销毁。但是如果有用户始终没有上线，该缓存池的空间就将被持续占用，并且房间的建立者无法新建房间。

预期解决方式：使用cron计时器和Redis的消息队列，通过cron控制招募是否到时，如果到时则读取该招募的所有用户数据，通过数据计算发送奖励的JSON包内容，并将包放入消息队列，以广播形式发送给本招募的所有参与者。

## 缓存雪崩

问题：使用了Redis缓存池的方式来处理招募数据，优点是的确运行效率很高，基本平均在我们实际使用测试时每个请求的相应时间都在几ms左右，但是由于数据存在与内存，一旦服务器出现了突然断电维护，内存中的所有数据将会失效。

![图表 18缓存雪崩原理图](https://github.com/DimDoremy/zxPomodoro-backend/blob/master/readme/clip_image034.jpg)

预期解决方式：

1. 在当前服务器承受的住的情况下搭建Redis集群，分布式存储缓存数据（但是由于当前服务器E5 2650v2单核 + 1G内存 + 40G使用空间，搭建集群略显困难）

2. 设置停机写入存储层，如果发生意外停机将当前数据写入MySQL的一个数据库临时存储。
