package test


/**
	select相关
	如果只有一个case收到响应,就执行该case,
	如果同时多个case收到响应,就随机执行其中一个,意味着其他case收到的数据将被抛弃.
	如果所有case都没有收到响应,但有default,就执行default

	在select中可以在一个case 中写 case v <- c:然后将c这个chan初始值设为nil,也就是只定义,不赋值,
	该case在c为nil的时候就不会生效.

	time.After(10 * time.Second) 会返回一个chan,在10s后,该chan会收到一个消息,可以用作延迟器
	time.Ticker(time.Second) 每秒会收到一个消息,可以用作定时器


	将一个代码块用func(){xxxx}()匿名函数包裹,然后在里面写lock和defer unlock,可以保证该代码块的线程安全
	go run -race 参数可以用来检测变量的线程冲突
 */
