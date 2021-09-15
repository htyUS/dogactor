package rank

import (
	"fmt"
	"testing"
)

func TestRank(t *testing.T) {
	rank := New()
	rank.Add("张三", Score(43))
	rank.Add("李四", Score(12))
	rank.Add("王五", Score(678))
	rank.Add("赵六", Score(54, 1000))
	rank.Add("孙七", Score(54, 2000))
	rank.Add("周八", Score(67,1111))
	rank.Add("吴九", Score(67,2222))
	rank.Add("郑十", Score(67,1000))
	rank.Add("郑十", Score(68))
	fmt.Println(rank.Get(1, 100))

	rank.Add("孙七", Score(1, 5))
	fmt.Println(rank.Get(1, 100))

	mem1 := rank.Get(5)
	fmt.Println(mem1[0].Key)


	data := rank.Marshal()
	fmt.Println(data)

	rank2 := New()
	_=rank2.UnMarshal(data)
	fmt.Println(rank2.Get(1, 100))
}