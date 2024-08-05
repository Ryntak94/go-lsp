package keywords

type (
	KeywordType  string
	TypeDataType string
)

const (
	Const           KeywordType = "const"
	FuncKeywordType KeywordType = "func"
	Method          KeywordType = "method"
	Var             KeywordType = "var"
	KType           KeywordType = "type"
)

const (
	Interface         TypeDataType = "interface"
	FuncDataType      TypeDataType = "func"
	PrimitiveDataType TypeDataType = "primitive"
)

type Keyword interface {
	Name() string
	Type() KeywordType
	Package() string
}

type ConstKeyword struct {
	name     string
	kType    KeywordType
	DataType string
	kPackage string
}

func (ck *ConstKeyword) Name() string {
	return ck.name
}

func (ck *ConstKeyword) Type() KeywordType {
	return ck.kType
}

func (ck *ConstKeyword) Package() string {
	return ck.kPackage
}

type VarKeyword struct {
	name     string
	kType    KeywordType
	DataType string
	kPackage string
}

func (vk *VarKeyword) Name() string {
	return vk.name
}

func (vk *VarKeyword) Type() KeywordType {
	return vk.kType
}

func (vk *VarKeyword) Package() string {
	return vk.kPackage
}

type FuncKeyword struct {
	name          string
	kType         KeywordType
	ParameterList []string
	ReturnList    []string
	kPackage      string
	TypeParameters
}

func (fk *FuncKeyword) Name() string {
	return fk.name
}

func (fk *FuncKeyword) Type() KeywordType {
	return fk.kType
}

func (fk *FuncKeyword) Package() string {
	return fk.kPackage
}

type MethodKeyword struct {
	FuncKeyword
	ReceiverType string
}

func (mk *MethodKeyword) Name() string {
	return mk.name
}

func (mk *MethodKeyword) Type() KeywordType {
	return mk.kType
}

func (mk *MethodKeyword) Package() string {
	return mk.kPackage
}

type TypeKeyword struct {
	name     string
	kType    KeywordType
	DataType TypeDataType
	kPackage string
	*TypeProperty
	*InterfaceType
}

func (tk *TypeKeyword) Name() string {
	return tk.name
}

func (tk *TypeKeyword) Type() KeywordType {
	return tk.kType
}

func (tk *TypeKeyword) Package() string {
	return tk.kPackage
}

type ReservedKeyword struct {
	name     string
	kType    KeywordType
	kPackage string
}

func (rk *ReservedKeyword) Name() string {
	return rk.name
}

func (rk *ReservedKeyword) Type() KeywordType {
	return rk.kType
}

func (rk *ReservedKeyword) Package() string {
	return rk.kPackage
}

type TypeParameters struct {
	TypeParameterList []string
}

type InterfaceType struct {
	Method string
}

type KeyValue struct {
	Key   string
	Value string
}

type TypeProperty struct {
	KeyValue
}
