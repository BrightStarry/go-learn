package test

/**
	同步相关
 */
/**
	sync.WaitGroup
	该对象类似java中的CountDownLatch.可以让多个程序相互等待.
	例如Add(20),那么只有调用20次Done()方法后,之前执行该对象的Wait()方法的线程才会继续执行
 */

 /**
 	time.After(10 *time.Second) 会返回一个chan,
 	在10s后,该chan会收到一个数据.

 	如果将这个代码作为case写在select中,那么当其他所有case,在10s内没有数据,就会执行这个case

 	time.Tick(10 *time.Second) 会返回一个chan,
 	每隔10s,该chan会收到一个数据

  */
