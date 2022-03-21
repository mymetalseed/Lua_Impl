package api

type LuaType = int
type ArithOp = int
type CompareOp = int	
/*-------GO函数渗透部分--------*/
type GoFunction func(LuaState) int

type LuaState interface {
	/* 基础栈操作 */
	GetTop()								int
	AbsIndex(idx int)						int
	CheckStack(n int)   					bool
	Pop(n int)			
	Copy(fromIdx,toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx,n int)
	SetTop(idx int)
	TypeName(tp LuaType) 					string
	Type(idx int)							LuaType
	IsNone(idx int)							bool
	IsNil(idx int)							bool
	IsNoneOrNil(idx int)					bool
	IsBoolean(idx int)						bool
	IsInteger(idx int)						bool
	IsNumber(idx int)						bool
	IsString(idx int)						bool
	ToBoolean(idx int) 						bool
	ToInteger(idx int)						int64
	ToIntegerX(idx int)						(int64,bool)
	ToNumber(idx int)						float64
	ToNumberX(idx int)						(float64,bool)
	ToString(idx int)						string
	ToStringX(idx int)						(string,bool)
	/* push函数 go->stack */
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)

	//执行算术和按位运算
	Arith(op ArithOp)
	//比较运算
	Compare(idx1,idx2 int, op CompareOp) 	bool
	Len(idx int)
	//string的concat方法
	Concat(n int)

	/*get function (Lua -> stack)*/
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) 						LuaType
	GetField(idx int,k string)				LuaType
	GetI(idx int,i int64)					LuaType

	/*set function (stack -> Lua)*/
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, n int64)

	Load(chunk []byte,chunkName,mode string)int
	Call(nArgs,nResults int)
	PrintStack()

	/*Go函数部分*/
	PushGoFunction(f GoFunction)
	IsGoFunction(idx int) 					bool
	ToGoFunction(idx int) 					GoFunction

	/*操作全局变量*/
	PushGlobalTable()
	GetGlobal(name string) 					LuaType
	SetGlobal(name string)
	Register(name string,f GoFunction)
}