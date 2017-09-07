Evernote golang API
---

以前使用的dreampuf/evernote-sdk-golang库， 由于evernote api版本的升级，无法正常链接Evernote，
尝试重新生成新的代码

## 准备thrift-go包
```
    go get git.apache.org/thrift.git
```
同时需要将thrift切换倒0.10.0版本    
```
    git checkout -b 0.10.0 origin/0.10.0
```

## Evernote API Generation
用thrift将各个Thrift IDL文件生成代码
```
thrift -strict -nowarn --allow-64bit-consts --allow-neg-keys --gen go ./src/UserStore.thrift
```

## Modified the generated code  
生成代码之后编译时，会有错误信息，主要是生成代码使用append函数的时候类型不匹配，修改生成代码

##
