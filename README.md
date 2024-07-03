# Telegram 发卡机器人

一个方便部署的Telegram发卡机器人，用户支付成功后会发送货物内容

### 特性
    ①使用易支付通用接口
    ②使用sqlite作为数据库，闭着眼睛部署
    ③一个用户只能开一个订单，防止滥用

### DEMO


### 使用说明
自行编译或者下载右侧Release编译好的运行文件（如果下载运行文件，需要执行命令给运行权限`chmod +x tg_faka_linux`），直接运行即可，运行文件的同一目录下要附带.env文件夹并在文件夹中写好配置文件。程序成功运行后，配置nginx，反代到运行端口即可。

### 支付说明
使用通用易支付接口，[Mopay](mopay.vip)或自行寻找易支付进行接入，支付信息填写到epay_config.json

因为需要接收支付成功的回调信息，所以需要对外开放http请求，运行时附带参数可修改http运行端口 `--port 8087`

### nginx反代参考配置
    location / {
        proxy_pass http://127.0.0.1:8087;
    }


### 配置说明
配置文件保存在`.env`文件夹中，修改即可
#### config.json
    "tg_bot_token": 机器人的token，前往BotFather获取，并设置为群组的管理员
    "admin_tg_id": 管理员Telegram Chat ID,可以在@userinfobot获取
    "order_duration_minutes": 订单持续时间
    "host": 绑定的域名，用于易支付发起订单时拼接回调地址
    "proxy": 代理，一般不用开

#### epay_config.json
    "pid": 易支付的pid，数字用双引号括起来
    "key": 易支付的key
    "url": 易支付发起订单的url，有些易支付后台显示的url不以submit.php结尾，可能要自己加上
    "pay_type": 易支付的支付类型
    
    "notify_url": 保持默认


## 管理员指令

### 添加商品
```
/add_products
商品名1 介绍1 金额1
商品名2 介绍2 金额2
商品名3 介绍3 金额3
```

示例:
```
/add_products
宝宝金水 集“祛痱、止痒、防蚊虫”三种功效为一体，无油腻、无刺激。 8848
```


### 查看商品id
`/view_products`


### 添加库存
```
/add_product_items /商品id/
库存1发货内容
库存2发货内容
库存3发货内容
```

示例:
```
/add_product_items /7db6bc44-7265-4060-91fb-634fbffe11f5/
宝宝金水亲子装
宝宝金水豪华礼包数字版
```

### 清空商品
`/clear_products`
### 清空库存(如有库存，则会以文件发送给管理员)
`/clear_product_items`




Tg: [@nulllllllll](https://t.me/nulllllllll)
