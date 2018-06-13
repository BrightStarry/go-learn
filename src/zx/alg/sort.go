package alg

/*基础排序*/

/**
	选择排序  O(n^2)
	从数组中依次找出最小的元素,放到前面(第一个,第二个第三个...)
 */
func selectSort(arr []int, size int) []int {
	// 遍历数组, 第一次遍历[0:],第二次[1:]
	for i := 0; i < size; i++ {
		// 寻找[i,n)区间内的最小值
		// 子循环当前的最小值索引
		minIndex := i
		for j := i + 1; j < size; j++ {
			// 如果该数小于目前的最小数,则记录该数索引
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		// 交换:将该最小值放到前面
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
	return arr
}

/**
	插入排序  O(n^2)
	从一个无序数组中,依次拿出元素,插入一个有序数组的正确位置,进行排序
	如果数组本身就比较有序,但不是完全有序,插入排序的效率就很高
	例如将随机数的范围变小,然后就会发现选择排序效率基本不变,但插入排序效率提高了很多
	或者使用generateRandomArray方法生成近乎有序的数组,会发现性能更高
 */
func insertSort(arr []int, size int) []int {
	// 默认第一个数为有序，从第二个元素开始遍历
	for i := 1; i < size; i++ {
		// 当 当前元素i 大于 前个元素（该元素属于有序队列），则当前元素也属于有序，无需操作
		if arr[i] > arr[i-1] {
			continue
		}

		// 遍历有序队列
		/**
			方案一， 比较出比当前元素i小的元素的位置，然后删除当前元素i，然后再该位置前一个位置插入 当前元素i的值
			因为切片插入删除需要大量复制，所以很慢
			10000个元素，执行时间122+ms
		 */
		//temp := arr[i]
		//for j:=0; j < i;j++{
		//	// 如果当前元素i小于 元素j， 则
		//	if arr[i] < arr[j]{
		//		// 删除 元素i
		//		arr = append(arr[:i],arr[i+1:]...)
		//		// 插入它到元素j-1的位置
		//		arr =  append(arr[:j],append([]int{temp},arr[j:]...)...)
		//	}
		//}

		/**
			方案二：从后往前（从i元素开始）遍历有序队列，如果 元素x 小于 元素x-1,就交换两者位置
			10000个元素，执行时间:30+ms
			1,5,6,7,3,...
			1,5,6,3,7
			1,5,3,6,7
			1,3,5,6,7
		 */
		//for j:=i; j >0 && arr[j] < arr[j-1]; j--{
		//	arr[j],arr[j-1] = arr[j-1],arr[j]
		//}

		/**
			方案三：优化了方案二，不再每次进行交换
			而是先将元素i提取出来,然后从有序队列末尾依次开始比较,如果元素i不应放在x位置,就将x位置的元素往后复制一份.
			10000个元素，执行时间 15+ms
		 */
		temp := arr[i]
		var j int
		for j = i; j > 0 && arr[j-1] > temp; j-- {
			arr[j] = arr[j-1]
		}
		arr[j] = temp
	}
	return arr
}

/**
	冒泡排序 O(n^2) TODO 等待优化
	外循环,对于一个size为x的元素，需要进行x-1次外循环
	内循环,每次都从0开始，一次比较相邻两个元素，将小的元素放到前面
	如此，每一次外循环后，数组末尾就会有多一个有序的元素
	10000个元素，执行时间120+ms
 */
func bubbleSort(arr []int, size int) []int {
	for i := 0; i < size-1; i++ {
		for j := 0; j < size-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	return arr
}

/**
	希尔排序 O(nlogn) TODO 等待优化
	和快速排序一样采用分治法，唯一能和快速排序性能一较高低的算法，能突破O(n^2)

	数组 9,8,7,6,5,4,3,2,1 如果初始增量size/2为4，则
	序列号 1,2,3,4,1,2,3,4,1
	数组   9,8,7,6,5,4,3,2,1
	然后分别循环每个子序列，对其进行插入排序
	然后再将增量/2 = 2...直到增量为

	10000个元素，执行时间 900微秒
 */
func shellSort(arr []int, size int) []int {
	in := size / 2 // 起始增量,为一半元素
	// 循环到增量为0 (1/2=0)
	for in >= 1 {
		// 循环 in次，循环所有子数组
		for i1:=0; i1<in; i1++{
			// 对每个子数组进行插入排序
			// 增量为in时， 0，0+in,0+in+in为一个子数组， 1,1+in,1+in+in是一个子数组
			for i2:= i1+in;i2 < size;i2+=in {
				temp := arr[i2]
				i3 := i2 - in
				for ;i3 >= i1 && arr[i3] > temp; i3 -= in {
					arr[i3+in] = arr[i3]
				}
				arr[i3+in] =temp
			}
		}
		in = in / 2
	}
	return arr
}


/**
	归并排序 O(nlogn)
	空间占用稍大
 */
func mergeSort(arr []int, size int) {

}
