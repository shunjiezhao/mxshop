
1. web利用token来获取userid
2. 发送给pk 然后利用join来进入 传入find type 选择敌人 随机 复仇
3. 根据类型来判断敌人
   1. 复仇 查对局表,id,id,time,winner; 如果全胜或者没有对局就返回错误
   2. 随机选择 利用 choose 接口
   3. 选择敌人 传入 userid
4. 