package main

import "fmt"
import "io/ioutil"
import "os"
import "luago/binchunk"
import . "luago/vm"
import . "luago/api"
import "luago/state"

//简化版的模拟读取chunk
func main(){
	if len(os.Args) > 1 {
		fmt.Printf("\nLua字节流测试\n")
		data,err := ioutil.ReadFile(os.Args[1])
		if err != nil { panic(err) }
		proto := binchunk.Undump(data)
		list(proto)
	}

	fmt.Printf("\n栈测试\n")
	ls := state.New()
	ls.PushInteger(1)
	ls.PushString("2.0")
	ls.PushString("3.0")
	ls.PushNumber(4.0)
	printStack(ls)

	ls.Arith(LUA_OPADD)
	printStack(ls)
	ls.Arith(LUA_OPBNOT)
	printStack(ls)
	ls.Len(2)
	printStack(ls)
	ls.Concat(3)
	printStack(ls)
	ls.PushBoolean(ls.Compare(1, 2, LUA_OPEQ))
	printStack(ls)
}


//显示chunck的参数
func list(f *binchunk.Prototype){
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _,p :=range f.Protos {
		list(p)
	}
}

//显示chunk的头部
func printHeader(f *binchunk.Prototype){
	//先判断是否是函数
	funcType := "main"
	if f.LineDefined > 0 {funcType = "function"}

	varargFlag := ""
	if f.IsVararg > 0 {varargFlag = "+"}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n",funcType,f.Source,f.LineDefined,f.LastLineDefined,len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues,",f.NumParams,varargFlag,f.MaxStackSize,len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",len(f.LocVars),len(f.Constants),len(f.Protos))
}

//打印指令
func printCode(f *binchunk.Prototype){
	for pc,c := range f.Code{
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d",f.LineInfo[pc])
		}

		i := Instruction(c)
		fmt.Printf("\t%d\t[%s]\t0x%08X\t%s\t",pc+1,line,c,i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

//打印参数细节,如常量表,局部变量表,upvalue表
func printDetail(f *binchunk.Prototype){
	fmt.Printf("constants (%d):\n",len(f.Constants))
	for i,k := range f.Constants{
		fmt.Printf("\t%d\t%s\n",i+1,constantToString(k))
	}

	fmt.Printf("locals (%d):\n",len(f.LocVars))
	for i,locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",i,locVar.VarName,locVar.StartPC+1,locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n",len(f.Upvalues))
	for i,upVar := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",i,upvalName(f,i),upVar.Instack,upVar.Idx)
	}
}

//显示常量的值
func constantToString(k interface{}) string {
	switch k.(type){
	case nil: 				return "nil"
	case bool: 				return fmt.Sprintf("%t",k)
	case float64: 			return fmt.Sprintf("%g",k)
	case int64: 			return fmt.Sprintf("%d",k)
	case string: 			return fmt.Sprintf("%g",k)
	default: 				return "?"
	}
}

func upvalName(f *binchunk.Prototype,idx int) string {
	if len(f.UpvalueNames) > 0{
		return f.UpvalueNames[idx]
	}
	return "-"
}

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()

		fmt.Printf("%d", a)
		if i.BMode() != OpArgN {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case IABx:
		a, bx := i.ABx()

		fmt.Printf("%d", a)
		if i.BMode() == OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == OpArgU {
			fmt.Printf(" %d", bx)
		}
	case IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	}
}


func printStack(ls LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}