package lib_writer

import (
	"strings"

	"github.com/rossmacarthur/cases"
)

type writer = func(strings.Builder)

type generator struct {
	frontendName string
	text         strings.Builder
}

func newGenerator(frontendName string) *generator {
	return &generator{frontendName: cases.ToPascal(frontendName)}
}

func (self *generator) writeImportBlock(content func()) {
	self.text.WriteString("import (")
	content()
	self.text.WriteString(")\n\n")
}

func (self *generator) writePackageLabel(name string) {
	self.text.WriteString("package ")
	self.text.WriteString(cases.ToSnake(name))
	self.text.WriteString("\n\n")
}

func (self *generator) writeReturnStatement(statement func()) {
	self.text.WriteString("return ")
	statement()
	self.text.WriteRune('\n')
}

func (self *generator) writeStructConstruction(name string, body func()) {
	self.text.WriteString(name)
	self.text.WriteString("{\n")

	body()

	self.text.WriteString("\n}")
}

func (self *generator) writeStringLiteral(content string) {
	self.text.WriteRune('"')
	self.text.WriteString(strings.ReplaceAll(content, "\"", "\\\""))
	self.text.WriteRune('"')
}

func (self *generator) writeJsonStructItem(name string, ty func()) {
	// does this: name ty `json:"name"`

	self.text.WriteString(cases.ToPascal(name))
	self.text.WriteRune(' ')
	ty()
	self.text.WriteString(" `json:\"")
	self.text.WriteString(name)
	self.text.WriteString("\"`\n")
}

func (self *generator) writeCommaSeperatedList(items []func()) {
	for index, fn := range items {
		if index != 0 {
			self.text.WriteString(", ")
		}

		fn()
	}
}

func (self *generator) writeStruct(name string, body func()) {
	self.text.WriteString("type ")
	self.text.WriteString(name)
	self.text.WriteString(" struct {\n")

	body()

	self.text.WriteString("\n}\n\n")
}

func (self *generator) writeInterface(name string, body func()) {
	self.text.WriteString("type ")
	self.text.WriteString(name)
	self.text.WriteString(" interface {\n")

	body()

	self.text.WriteString("\n}\n\n")
}

func (self *generator) writeMethodStart(name string, ptrTypeAssociation string) {
	self.text.WriteString("func (self *")
	self.text.WriteString(ptrTypeAssociation)
	self.text.WriteString(") ")
	self.text.WriteString(name)
}

func (self *generator) writeFunctionArgs() {

}

func (self *generator) writeFunctionBody(returnType string, body func()) {
	if len(returnType) != 0 {
		self.text.WriteRune(' ')
		self.text.WriteString(returnType)
	}

	self.text.WriteString(" {")
	body()
	self.text.WriteString("}")
}

func (self *generator) writeFunc(name string, ptrTypeAssociation string, args func(), rtr string, body func()) {
	self.text.WriteString("func ")

	if len(ptrTypeAssociation) != 0 {
		self.text.WriteString("(self *")
		self.text.WriteString(ptrTypeAssociation)
		self.text.WriteString(") ")
	}

	self.text.WriteString(name)

	self.text.WriteString("(")
	args()
	self.text.WriteString(")")

	if len(rtr) != 0 {
		self.text.WriteRune(' ')
		self.text.WriteString(rtr)
	}

	self.text.WriteString(" {\n")
	body()
	self.text.WriteString("\n}\n\n")
}

func (self *generator) writeComment(comment string) {
	if len(comment) == 0 {
		return
	}

	self.text.WriteString("// ")

	for _, char := range comment {
		if char == '\n' {
			self.text.WriteString("\n// ")
		} else {
			self.text.WriteRune(char)
		}
	}

	self.text.WriteRune('\n')
}

func stackToPublicName(stack []string) string {
	writer := strings.Builder{}

	for _, item := range stack {
		writer.WriteString(cases.ToPascal(item))
	}

	return writer.String()
}
