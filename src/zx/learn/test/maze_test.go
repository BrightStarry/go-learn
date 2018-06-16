package test

import (
	"os"
	"fmt"
	"testing"
)

/**
	简单的广度迷宫算法
 */

/**
	将迷宫文件读取为 二维数组
 */
func readMaze(filename string) [][]int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	// 从reader(file)中读取行数和列数
	var row, col int
	fmt.Fscanf(file, "%d %d\r\n", &row, &col)

	// 分片,每个元素都是 []int(表示一行的所有元素),长度为行数
	maze := make([][]int, row)
	// 遍历每一行
	for i := range maze {
		// 构造每一行的 元素个数(也就是列数)
		maze[i] = make([]int, col)
		// 循环一行的所有元素
		for j := range maze[i] {
			// 读取该行每一列的数据
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
		// 在windows中,每次读完一行,需要再次调用该方法,否则每一行的第一个元素会读为0
		// mac中换行符是\n, win中是\r\n
		fmt.Fscanln(file)
	}
	return maze
}

/**
	表示迷宫的某个点的位置
 */
type point struct {
	i, j int
}

/**
	位置和位置的相加操作
 */
func (p point) add(r point ) point {
	return point{p.i+r.i, p.j+r.j}
}

/**
	返回某个二维数组(迷宫,steps等)在当前位置的值,
	以及是否未越界
 */
func (p point) at(grid [][]int) (int,bool){
	if p.i < 0 || p.i >= len(grid){
		return 0, false
	}
	if p.j < 0  || p.j >=len(grid[p.i]) {
		return 0,false
	}
	return grid[p.i][p.j],true
}


/**
	代表四个方向
 */
var directions = [4]point{
	// 上
	{-1,0},
	// 左
	{0,-1},
	// 下
	{1,0},
	// 右
	{0,1},
}

/**
	走迷宫
 */
func walk(maze [][]int, start, end point) [][]int {
	// 构造 存储迷宫走法的数组, 行数列数和迷宫数组一样,只不过不附上值(默认都是0)
	steps := make([][]int, len(maze))
	for i := range steps{
		steps[i] = make([]int,len(maze[i]))
	}

	// 路径从起点开始
	q := []point{start}
	for len(q) > 0{
		// 获取队列头
		current := q[0]
		// 删除队列头
		q = q[1:]

		// 如果当前要探索的点是终点,也就是到达终点了退出
		if current == end{
			break
		}

		// 遍历四个方向
		for _,dir := range directions{
			// 获取某个方向的下个节点
			next := current.add(dir)

			// 需要下个位置在迷宫中是0(1代表墙),
			// 并且steps也是0(表示没走过)
			// 并且不是起点
			val,ok := next.at(maze)
			// 越界,或撞墙
			if !ok || val == 1{
				continue
			}
			val,ok = next.at(steps)
			// 越界,或走过了
			if !ok || val !=0{
				continue
			}
			// 等于起点
			if next == start{
				continue
			}

			// 如果要进行探索
			// 将当前位置的探索次数加1
			currentSteps,_ := current.at(steps)
			steps[next.i][next.j] = currentSteps + 1
			// 将找到的新的位置,放入待探索队列
			q = append(q,next)
		}

	}
	return steps
}

func TestMaze(t *testing.T) {
	maze := readMaze("./maze.in")
	// 打印迷宫
	for _, row := range maze {
		for _, val := range row {
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()

	// 起始点为左上角,0,0,  终点为右下角, 也就是行数-1和列数-1
	steps := walk(maze, point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1})
	// 打印步骤
	for _, row := range steps {
		for _, val := range row {
			fmt.Printf("%3d ", val)
		}
		fmt.Println()
	}
}
