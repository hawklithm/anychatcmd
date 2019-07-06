# anychatcmd  [![star this repo](http://github-svg-buttons.herokuapp.com/star.svg?user=hawklithm&repo=anychatcmd&style=flat&background=1081C1)](http://github.com/hawklithm/anychatcmd) [![fork this repo](http://github-svg-buttons.herokuapp.com/fork.svg?user=hawklithm&repo=anychatcmd&style=flat&background=1081C1)](http://github.com/hawklithm/anychatcmd/fork) ![Build](https://camo.githubusercontent.com/46cb8b3469febc6cdb6fbaea2ef1517c396004e7/68747470733a2f2f7472617669732d63692e6f72672f736a77686974776f7274682f676f6c6561726e2e706e673f6272616e63683d6d6173746572)

## INSTALL

```bash
git clone https://github.com/hawklithm/anychatcmd.git
cd anychatcmd
dep ensure -update -v
go build
./anychatcmd  #启动anychatcmd
```

## 最新资讯

1. 代码重构基本完成

2. 基于iterm协议实现了命令行图片展示
    
    目前仅支持iterm，xterm未验证，后续会移植linux版本
    
    使用方法：
    
    在执行```./anychatcmd```之前执行 ```export WECHAT_TERM=iterm```
    
    效果图：
    
    ![test](https://github.com/hawklithm/anychatcmd/blob/master/test/test.png?raw=true)
    


## 背景

最初考虑做pc版微信替代品的出发点是公司安全方面原因，(公司出于安全性考虑不允许安装pc版wechat，网页版在使用上并不令人满意)
，但是后来在做的过程中发现，不止可以尝试微信cmd版，还可以尝试一些其他软件的cmd版(个人无聊爱好)

之前已完成针对微信的cmd版本([hawklithm/wechatcmd](https://github.com/hawklithm/wechatcmd)),现在开始，在对老版本重构的基础上，探索一下其他的应用，希望能有些比较有意思的事情。在一番调研之后决定采用 

*本代码主要在MAC OS上进行开发测试，针对linux系统的兼容主要基于ubuntu进行考虑的，如果在实际使用中存在什么问题欢迎提出，暂不考虑windows*

目前已完善点：

- [x] termui版本升级到3.0.0，接口兼容问题修复
- [x] 群聊天中发言人显示
- [x] 用户多端登陆时，通过其他端发出的消息的同步
- [x] 切换当前聊天窗口时，历史聊天记录的恢复

**注：本程序目的为日常使用替代pc端微信，所以不会开发自动回复或者聊天机器人抑或是群发之类的功能**


操作方式：

| 按键 | 说明 |
| --- | --- |
| Ctrl+n | 下一个聊天 |
| Ctrl+p | 上一个聊天 |
| Ctrl+j | 下一条聊天记录 |
| Ctrl+k | 上一条聊天记录 |
| Ctrl+w | 展示选中的聊天信息的详情；如果是图片则打开图片，如果是外链则打开外链 |
| Ctrl+c | 退出 |
| Ctrl+a | 开启/关闭消息提醒 |

开发计划：

- [x] 实现微信登陆
- [x] 实现微信认证
- [x] 实现拉取用户信息
- [x] 同步消息
- [x] 自动更新消息
- [x] 聊天
- [x] 群聊
- [x] 支持图片显示
- [x] 支持emoji表情
- [x] 解析分享消息
- [x] 解析公众号消息
- [x] 支持表情包
- [x] 消息提醒
- [ ] 界面优化(用户列表和当前会话分拆，支持群成员展示)

由于整体框架的原因，以下特性计划在代码重构之后再完成了:

- [ ] 用户检索
- [ ] 本地表情包发送(发图片)
- [ ] 自动保存消息到本地
- [ ] vim式操作

代码重构后计划增加的特性:

- [ ] 支持即刻网页版账号登陆(因为本人喜欢刷即刻)
- [ ] 支持Boss直聘网页版及一些自动化功能(纯工作需要.....)


