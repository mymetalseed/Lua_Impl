package vm

import "luago/api"

type Instruction uint32

const MAXARG_Bx = 1 << 18 - 1	
const MAXARG_sBx = MAXARG_Bx >> 1 

//获取指令码的方法集合

func (self Instruction) OpCode() int{
	return int(self & 0x3F)
}

//从IABC模式指令中提取参数  
func (self Instruction) ABC() (a,b,c int){
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}

func (self Instruction) ABx() (a,bx int){
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

func (self Instruction) AsBx() (a,sbx int){
	a,bx := self.ABx()
	return a,bx - MAXARG_sBx
}

func (self Instruction) Ax() int{
	return int(self >> 6)
}

func (self Instruction) OpName() string {
	return opcodes[self.OpCode()].name
}

func (self Instruction) OpMode() byte {
	return opcodes[self.OpCode()].opMode
}

func (self Instruction) BMode() byte {
	return opcodes[self.OpCode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.OpCode()].argCMode
}

func (self Instruction) Execute(vm api.LuaVM) {
	action := opcodes[self.OpCode()].action
	if action != nil {
		action(self, vm)
	} else {
		panic(self.OpName())
	}
}