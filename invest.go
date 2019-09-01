package main
import(
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"log"
	"container/list"
	"strconv"
	"strings"
)

type Item struct {
	date string
	value float64
	level float64
}

type Invest struct {
	date string
	pe_value float64
	pe_level float64
	pb_value float64
	pb_level float64
	value float64
	level float64
}

func index_level(file_name string, date_index int, value_index int) (*list.List,int){
	file,err:=os.Open(file_name)
	defer file.Close()

	reader:=csv.NewReader(file)
	k := 0
	l := list.New()
	l_len := 0
	for{
		record,err:=reader.Read()
		if err==io.EOF {
			break
		} else if err!=nil {
//			fmt.Println("Error:",err)
			break
		}
		if k > 0 && len(record) >=4 {
			value,err:=strconv.ParseFloat(record[value_index],64);
			if err==nil {
				l.PushBack(Item{record[date_index],value,0})
				l_len=l_len+1
				fmt.Println(record[date_index],":",record[value_index])
			}
		}
		k = k+1
	}
	fmt.Println("list size=",l_len)

	item_list, item_list_len := count_level(l,l_len)

	return item_list,l_len
}

func count_level(l *list.List, l_len int)(*list.List,int){
	for x := l.Front(); x!=nil; x=x.Next() {
		x_item := x.Value.(Item)
		level := 0;
		for y := l.Front(); y!=nil; y=y.Next() {
			y_item := y.Value.(Item)
			if x_item.value > y_item.value {
				level=level+1;
			}
		}
		x_item.level=(level*100)/l_len
		fmt.Println("date:",x_item.date,",value:",x_item.value,", level=",x_item.level)
	}
	return l,l_len
}

func main(){
	flag.Parse()
	fmt.Println(flag.NArg())
	if flag.NArg() !=5 {
		log.Println("usage: command filename date_index value_index")
		return
	}
	pe_file_name := flag.Arg(0)

	pe_date_index, err := strconv.Atoi(flag.Arg(1))
	if err!=nil {
		fmt.Println("Error:",err)
	}

	pe_value_index, err := strconv.Atoi(flag.Arg(2))
	if err!=nil {
		fmt.Println("Error:",err)
	}

	pb_file_name := flag.Arg(3)

	pb_date_index, err := strconv.Atoi(flag.Arg(4))
	if err!=nil {
		fmt.Println("Error:",err)
	}

	pb_value_index, err := strconv.Atoi(flag.Arg(5))
	if err!=nil {
		fmt.Println("Error:",err)
	}

	pe_list,pe_list_len := index_level(pe_file_name,pe_date_index,pe_value_index)
	pb_list,pb_list_len := index_level(pb_file_name,pb_date_index,pb_value_index)

	if pe_list_len != pb_list_len {
		fmt.Println("Error: pe_list_len=",pe_list_len, ", pb_list_len=",pb_list_len)
		return
	}

	invest_list := list.New()
	item_list := list.New()
	item_list_len := 0
	pe := pe_list.Front()
	pb := pb_list.Front()
	for {
		if pe!=nil || pb!=nil {
			break
		}
		pe_item := pe.Value.(Item)
		pb_item := pb.Value.(Item)
		if strings.EqualFold(pe_item.date, pb_item.date) {
			value := (pe_item.level+pb_item.level)/2
			invest_list.PushBack(Invest{pe_item.date,pe_item.value,pe_item.level,pb_item.value,pb_item.level,value,0})
			item_list.PushBack(Item{pe_item.date,value,0})
		}
		pe=pe.Next()
		pb=pb.Next()

	}

	item_list, item_list_len = count_level(item_list)
	invest := invest_list.Front()
	item := item_list.Front()
	for {
		if invest==nil || item==nil {
			break
		}
		invest_value := invest.Value.(Invest)
		item_value := item.Value.(Item)
		if strings.EqualFold(invest_value.date, item_value.date) {
			invest_value.level=item_value.level
		}
		invest=invest.Next()
		item=item.Next()
	}

	fmt.Println("date pe_value pe_level pb_value pb_level value level")

	for invest := invest_list.Front(); invest!=nil; invest=invest.Next() {
		invest_value := invest.Value.(Invest)
		fmt.Println(invest_value.date," ", invest_value.pe_value," ", invest_value.pe_level," ", invest_value.pb_value, " ",  invest_value.pb_level, " ",  invest_value.value, " ",  invest_value.level)
	}
}
