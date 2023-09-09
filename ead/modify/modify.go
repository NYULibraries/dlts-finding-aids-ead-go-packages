package modify

import (
	"fmt"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/xpath"
)

func ModifyEAD(data []byte) []string {

	var errors = []string{}

	p := parser.New()
	doc, err := p.Parse(data)
	if err != nil {
		errors = append(errors, "Unable to parse XML file")
		return append(errors, err.Error())
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if err != nil {
		errors = append(errors, "Unable to extract root node")
		return append(errors, err.Error())
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		errors = append(errors, "Unable to initialize XPathContext")
		return append(errors, err.Error())
	}
	defer ctx.Free()

	exprString := `//origination[@label]`
	nodes := xpath.NodeList(ctx.Find(exprString))

	for _, n := range nodes {
		fmt.Println(n)
	}

	// all ok, return empty slice
	return []string{}
}
