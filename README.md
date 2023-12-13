# redisHammer

简单的redis操作, 主要针对db scan、key scan、key del等不方便用命令行操作的动作进行简化

通用参数:

```
-addr: redis地址
-password: redis密码
-db: redis数据库
```

## scan

扫描所有的key

-scan: 扫描所有key
-keyType: 是否输出key类型

```shell
./bin/redisHammer 
  -addr=${host}:${port} 
  -password=${password} 
  -db=${db} 
  -scan 
  -keyType
  -match=
```

## get

获取string类型的key

```shell
./bin/redisHammer 
  -addr=${host}:${port} 
  -password=${password} 
  -db=${db} 
  -get
  -key=${key}
```

## del

删除string类型的key

```shell
./bin/redisHammer 
  -addr=${host}:${port} 
  -password=${password} 
  -db=${db} 
  -del
  -key=${key}
```

## keyScan

针对hash, set, zset类型的key进行scan

```shell
./bin/redisHammer 
  -addr=${host}:${port} 
  -password=${password} 
  -db=${db} 
  -keyScan
  -key=${key}
  -keyScanMatch=""
```

## info

```shell
./bin/redisHammer 
  -addr=${host}:${port} 
  -password=${password} 
  -info
```

## clientList

```shell
./bin/redisHammer 
  -addr=${host}:${port} 
  -password=${password} 
  -clientList
```

