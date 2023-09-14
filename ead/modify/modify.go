// Package modify provides functions that modify an EAD
package modify

import (
	"fmt"

	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/types"
	"github.com/lestrrat-go/libxml2/xpath"
)

// FABifyEAD modifies an ArchivesSpace-generated EAD []byte slice
// so that it is more compatible with the NYU Libraries Finding Aids Bridge (FAB)
// discovery system (https://github.com/NYULibraries/specialcollections).
//
// Please note:
// This function returns an EAD with subcontainers that have @id and @parent
// attributes.  To be FAB-compatible, the @id attribute must be deleted from the
// subcontainers. (The current lestrrat-go/libxml2 package does not support
// AttributeNode deletion.)
//
// For example:
// <container @id="please-delete-me" @parent="some-aspace-id" @type="Folder">
// still needs to be converted to:
// <container @parent="some-aspace-id" @type="Folder">
//
// The subcontainer @id attributes can be deleted with a tool like XMLStarlet:
// https://xmlstar.sourceforge.net/
// $ xmlstarlet ed -L -d '//_:container[@parent]/@id' my_ead.xml
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

	// register the default namespace
	prefix := `_`
	nsuri := `urn:isbn:1-931666-22-9`
	if err := ctx.RegisterNS(prefix, nsuri); err != nil {
		errors = append(errors, "Failed to register namespace")
		return "", append(errors, err.Error())
	}

	// find all nodes where <origination @label="Creator"> so that
	// "Creator" can be converted to the FAB-compatible value "creator"
	exprString := `//_:origination/@label[.='Creator']`
	nodes := xpath.NodeList(ctx.Find(exprString))
	for _, n := range nodes {
		n.SetNodeValue("creator")
	}

	// find all container nodes that may be the root of a container hierarchy
	// e.g., Box --> Folder --> Item
	//
	// The root of a container hierarchy does not have a @parent attribute.
	// Subcontainers are containers that *do* have a @parent attribute.
	// The FAB requires all subcontainers to have the @parent attribute value
	// set to the @id of the root container.
	exprString = `//_:container[@id and not(@parent)]/@id`
	nodes = xpath.NodeList(ctx.Find(exprString))
	for _, n := range nodes {
		rootID := n.NodeValue()
		err := updateSubContainer(ctx, rootID, rootID)
		if err != nil {
			errors = append(errors, "problem processing subcontainers")
			return "", append(errors, err.Error())
		}
	}

	return doc.String(), errors
}

func updateSubContainer(ctx *xpath.Context, parentID string, rootID string) error {
	// find the ID nodes of all containers whose @parent attribute value == parentID
	exprString := fmt.Sprintf("//_:container[@parent = \"%s\"]", parentID)
	containerNodes := xpath.NodeList(ctx.Find(exprString))

	// recursively process all containers
	for _, containerNode := range containerNodes {
		// get the value of this container's @id attribute
		idAttr, err := containerNode.(types.Element).GetAttribute("id")
		if err != nil {
			return fmt.Errorf("problem accessing @id attribute for subcontainers of parent \"%s\": %s", parentID, err)
		}

		id := idAttr.NodeValue()

		// update any subcontainers for whom this container is a parent
		updateSubContainer(ctx, id, rootID)

		// recursive calls are complete
		// update this node
		containerNode.(types.Element).SetAttribute("parent", rootID)
		containerNode.(types.Element).RemoveAttribute("id")
	}
	return nil
}
