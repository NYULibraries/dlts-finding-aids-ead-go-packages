package modify

import (
	"fmt"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/xpath"
)

func FABifyEAD(data []byte) (string, []string) {

	var errors = []string{}

	p := parser.New()
	doc, err := p.Parse(data)
	if err != nil {
		errors = append(errors, "Unable to parse XML file")
		return "", append(errors, err.Error())
	}
	defer doc.Free()

	root, err := doc.DocumentElement()
	if err != nil {
		errors = append(errors, "Unable to extract root node")
		return "", append(errors, err.Error())
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		errors = append(errors, "Unable to initialize XPathContext")
		return "", append(errors, err.Error())
	}
	defer ctx.Free()

	prefix := `_`
	nsuri := `urn:isbn:1-931666-22-9`
	if err := ctx.RegisterNS(prefix, nsuri); err != nil {
		errors = append(errors, "Failed to register namespace")
		return "", append(errors, err.Error())
	}
	exprString := `//_:origination/@label[.='Creator']`
	nodes := xpath.NodeList(ctx.Find(exprString))

	for _, n := range nodes {
		n.SetNodeValue("creator")
	}

	exprString = `//_:container[@id and not(@parent)]/@id`
	nodes = xpath.NodeList(ctx.Find(exprString))

	for _, n := range nodes {
		rootID := n.NodeValue()
		updateSubContainer(ctx, rootID, rootID)
	}

	// all ok, return empty slice
	return doc.String(), []string{}
}

func updateSubContainer(ctx *xpath.Context, parentID string, rootID string) {
	// find the ID nodes of all containers whose @parent attribute value == parentID
	exprString := fmt.Sprintf("//_:container[@parent = \"%s\"]/@id", parentID)
	nodes := xpath.NodeList(ctx.Find(exprString))

	// recursively process all containers
	for _, n := range nodes {
		// get the value of this @id node
		id := n.NodeValue()
		// update any subcontainers for whom this container is a parent
		updateSubContainer(ctx, id, rootID)

		// recursive calls are complete
		// update this node
		containerNode, _ := n.ParentNode()
		parentAttributeNode, _ := containerNode.Find(`./@parent`)
		// set the parent attribute to the rootID value
		parentAttributeNode.NodeList().First().SetNodeValue(rootID)
	}
}
