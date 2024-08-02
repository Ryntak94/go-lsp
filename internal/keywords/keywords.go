package keywords

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func GenerateKeywords(logger *log.Logger) Trie {
	keywordFile := generateKeywordsTxt(logger)

	scanner := bufio.NewScanner(keywordFile)

	keywordTree := CreateTrie()

	for scanner.Scan() {
		line := scanner.Text()
		keyword, err := getKeyword(line, logger)
		if err != nil || keyword == nil {
			continue
		}
		keywordTree.AddWord(keyword)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	reservedWordsFile := fmt.Sprintf("%s/internal/keywords/keywords.txt", cwd)
	reservedWords, err := os.Open(reservedWordsFile)

	if err != nil {
		panic(err)
	}

	scanner = bufio.NewScanner(reservedWords)
	for scanner.Scan() {
		keyword := NewReservedKeyword(scanner.Text())
		keywordTree.AddWord(keyword)
	}

	return keywordTree
}

func getKeyword(line string, logger *log.Logger) (Keyword, error) {
	var keywordType string
	osSpecificKeywords, _ := regexp.Compile("\\s\\(\\w+(-\\w+)+\\)")
	if matched := osSpecificKeywords.FindString(line); len(matched) > 0 {
		line = osSpecificKeywords.ReplaceAllString(line, "")
	}
	m := regexp.MustCompile("\\s#\\d*")
	line = m.ReplaceAllString(line, "")
	keyword, err := handleKeywordType(keywordType, line, logger)
	return keyword, err

}

func handleKeywordType(keywordType string, line string, logger *log.Logger) (Keyword, error) {
	switch keywordType {
	case "const":
		keyword := handleConst(line)
		if keyword != nil {
			return keyword, nil
		}
	case "var":
		keyword := handleVar(line)
		if keyword != nil {
			return keyword, nil
		}
	case "type":
		keyword := handleType(line)
		if keyword != nil {
			return keyword, nil
		}
	case "func":
		keyword := handleFunc(line, logger)
		if keyword != nil {
			return keyword, nil
		}
	case "method":
		keyword := handleMethod(line, logger)
		if keyword != nil {
			return keyword, nil
		}
	default:
	}
	return nil, errors.New("No keyword")
}

func handleConst(line string) *ConstKeyword {
	return NewConstKeyword(line)
}

func handleVar(line string) *VarKeyword {
	return NewVarKeyword(line)
}

func handleType(line string) *TypeKeyword {
	return NewTypeKeyword(line)
}

func handleFunc(line string, logger *log.Logger) *FuncKeyword {
	keyword, err := NewFuncKeyword(line, FuncKeywordType)
	if err != nil {
		logger.Print(fmt.Errorf("error: %w", err))
		return nil
	}
	return keyword
}

func handleMethod(line string, logger *log.Logger) *MethodKeyword {
	keyword, err := NewMethodKeyword(line)
	if err != nil {
		logger.Print(fmt.Errorf("error: %w", err))
		return nil
	}
	return keyword
}

func NewConstKeyword(line string) *ConstKeyword {
	pkgRegex, err := regexp.Compile("pkg \\S*")
	if err != nil {
		return nil
	}
	pkg := pkgRegex.FindString(line)[4:]
	splitLine := strings.Split(line, " ")
	var dataType string
	if splitLine[4] == "=" {
		dataType = splitLine[5]
	} else {
		dataType = splitLine[4]
	}
	keyword := ConstKeyword{name: splitLine[3], kType: "const", DataType: dataType, kPackage: pkg}
	return &keyword
}

func NewVarKeyword(line string) *VarKeyword {
	pkgRegex, err := regexp.Compile("pkg \\S*")
	if err != nil {
		return nil
	}
	pkg := pkgRegex.FindString(line)[4:]
	splitLine := strings.Split(line, " ")
	keyword := VarKeyword{name: splitLine[3], kType: "var", DataType: splitLine[4], kPackage: pkg}
	return &keyword
}

func NewReservedKeyword(word string) *ReservedKeyword {
	keyword := ReservedKeyword{name: word, kType: "reserved", kPackage: ""}
	return &keyword
}

func NewFuncKeyword(line string, keywordType KeywordType) (*FuncKeyword, error) {
	pkgRegex, err := regexp.Compile("pkg \\S*")
	if err != nil {
		return nil, err
	}
	pkg := pkgRegex.FindString(line)[4:]
	if pkg[len(pkg)-1:] == "," {
		pkg = pkg[:len(pkg)-1]
	}
	var paramList []string
	var returnList []string
	if deprecated, err := regexp.MatchString("deprecated", line); deprecated || err != nil {
		return nil, err
	}
	functionNameRegex, _ := regexp.Compile("[\\w\\d]+((\\[)|(\\())")
	bracketRegex, _ := regexp.Compile("func [a-zA-Z]*\\[.*\\]\\(")
	functionName := functionNameRegex.FindString(line)
	functionName = functionName[:len(functionName)-1]
	typeParameters := TypeParameters{[]string{}}
	if brackets := bracketRegex.FindString(line); len(brackets) > 0 {
		brackets = brackets[strings.Index(brackets, "[") : len(brackets)-1]
		typeParameters.TypeParameterList = strings.Split(brackets[1:len(brackets)-1], ", ")
	}
	var paramString string = ""
	var returnString string = ""
	if line[len(line)-1:] == ")" && line[len(line)-2:] != "()" {
		paramOrReturn := extractParenChunk(line)
		determineParamOrReturn := strings.Split(line, paramOrReturn)[0]

		if determineParamOrReturn[len(determineParamOrReturn)-1:] == " " && strings.Index(determineParamOrReturn, "(") > 0 {
			returnString = paramOrReturn
		} else {
			paramString = paramOrReturn
		}
		if len(returnString) > 0 {
			paramString = extractParenChunk(determineParamOrReturn)
		}
		paramList = strings.Split(paramString[1:len(paramString)-1], ", ")

		if len(returnString) > 0 {
			returnList = strings.Split(returnString[1:len(returnString)-1], ", ")
		}
	} else {
		lastRightParen := strings.LastIndex(line, ")")
		var returnType string
		if line[len(line)-6:] == "func()" {
			line = line[0 : len(line)-2]
			returnType = "func()"
		} else if lastRightParen == len(line)-1 {
			returnType = ""
		} else {
			returnType = line[lastRightParen+2:]
		}
		returnList = []string{returnType}
		paramString = extractParenChunk(line)
		paramList = strings.Split(paramString[1:len(paramString)-1], ", ")
	}
	keyword := FuncKeyword{functionName, keywordType, paramList, returnList, pkg, typeParameters}

	return &keyword, nil
}

func NewMethodKeyword(line string) (*MethodKeyword, error) {
	if deprecated, err := regexp.MatchString("deprecated", line); deprecated || err != nil {
		return nil, err
	}
	receiverRegex, _ := regexp.Compile("method \\(\\S*\\)")
	receiver := receiverRegex.FindString(line)
	receiver = receiver[strings.Index(receiver, "(")+1 : len(receiver)-1]
	funcKeyword, err := NewFuncKeyword(line, Method)
	if err != nil {
		return nil, err
	}
	methodKeyword := MethodKeyword{*funcKeyword, receiver}
	return &methodKeyword, nil
}

func NewTypeKeyword(line string) *TypeKeyword {
	if deprecated, err := regexp.MatchString("deprecated", line); deprecated || err != nil {
		return nil
	}

	pkgRegex, err := regexp.Compile("pkg \\S*")
	if err != nil {
	}
	pkg := pkgRegex.FindString(line)[4:]
	if pkg[len(pkg)-1:] == "," {
		pkg = pkg[:len(pkg)-1]
	}
	nameRegex, err := regexp.Compile("type \\S*")
	name := nameRegex.FindString(line)[5:]

	dataTypeIdx := strings.Index(line, name)
	dataType := line[dataTypeIdx+len(name)+1:]
	var typeDataType TypeDataType
	_ = typeDataType

	if strings.Index(dataType, "interface") == 0 {
		typeDataType = Interface
	} else if strings.Index(dataType, "func(") == 0 {
		typeDataType = FuncDataType
	} else {
		typeDataType = PrimitiveDataType
	}

	var interfaceMethod InterfaceType
	var typeProperty TypeProperty
	if typeDataType == Interface {
		if strings.Index(dataType, "interface,") == 0 {
			interfaceMethod.Method = dataType[11:]
		}
	} else if typeDataType == FuncDataType {
	} else if strings.Index(dataType, "interface") == -1 && strings.Index(dataType, ",") != -1 && strings.Index(dataType, "func(") == -1 {
		dataType = dataType[:strings.Index(dataType, ",")]
		keyValuePair := line[strings.LastIndex(line, ",")+2:]
		_ = keyValuePair
		typeProperty.Key = keyValuePair[:strings.LastIndex(keyValuePair, " ")]
		typeProperty.Value = keyValuePair[strings.LastIndex(keyValuePair, " ")+1:]
	}
	_ = interfaceMethod

	typeKeyword := TypeKeyword{name, KType, typeDataType, pkg, &typeProperty, &interfaceMethod}
	return &typeKeyword
}

func (fk *FuncKeyword) ParameterListString() string {
	return fmt.Sprintf("(%s)", strings.Join(fk.ParameterList, ", "))
}

func (fk *FuncKeyword) ReturnListString() string {
	return fmt.Sprintf("(%s)", strings.Join(fk.ReturnList, ", "))
}

func extractParenChunk(line string) string {
	lastRightParen := strings.LastIndex(line, ")")
	currLastLeftIdx := strings.LastIndex(line, "(")
	findLeftParen, _ := regexp.Compile("\\(")
	findRightParen, _ := regexp.Compile("\\)")
	numLeft := findLeftParen.FindAllString(line[currLastLeftIdx:], -1)
	numRight := findRightParen.FindAllString(line[currLastLeftIdx:], -1)
	for len(numLeft) != len(numRight) {
		currLastLeftIdx = strings.LastIndex(line[:currLastLeftIdx], "(")
		numLeft = findLeftParen.FindAllString(line[currLastLeftIdx:], -1)
		numRight = findRightParen.FindAllString(line[currLastLeftIdx:], -1)
	}
	return line[currLastLeftIdx : lastRightParen+1]
}
