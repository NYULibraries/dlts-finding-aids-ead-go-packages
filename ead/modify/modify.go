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
// Those modifications are:
// 1.) change <origination label="Creator">
//
//	to     <origination label="creator">
//
// 2.) for <container> hierarchies, set all subcontainer @parent
//
//	attribute values = to the @id of the root container
//	and delete the @id attribute from all subcontainers
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

	// rootIDs is a slice of strings containing the @id attribute values for all root containers
	var rootIDs []string

	// subContainers is a map keyed by the @parent attribute values of subcontainers
	subContainers := make(map[string]types.Node)

	// find all containers and divide them between root containers and subcontainers
	exprString = `//_:container`
	containerNodes := xpath.NodeList(ctx.Find(exprString))
	for _, containerNode := range containerNodes {

		// if a containerNode has a @parent attribute, then it is a subcontainer node
		// otherwise it is a root container node
		parentAttributeNode, err := containerNode.(types.Element).GetAttribute("parent")

		if err == nil {
			// subcontainer branch
			parentID := parentAttributeNode.NodeValue()
			if subContainers[parentID] != nil {
				errors = append(errors, fmt.Sprintf("error: detected multiple subcontainers with the same parentID: %s", parentID))
				errors = append(errors, "check if this EAD has already been \"FABified\"")
				return "", append(errors, err.Error())
			}
			subContainers[parentID] = containerNode
		} else {
			// root container branch
			idAttributeNode, err := containerNode.(types.Element).GetAttribute("id")
			if err != nil {
				errors = append(errors, "error: root container @id attribute is missing")
				return "", append(errors, err.Error())
			}
			rootIDs = append(rootIDs, idAttributeNode.NodeValue())
		}
	}

	// process all subcontainer hierarchies
	for _, rootID := range rootIDs {
		err := updateSubContainersUsingMap(subContainers, rootID, rootID)
		if err != nil {
			errors = append(errors, "error updating subcontainers of root container with @id=\"%s\"", rootID)
			return "", append(errors, err.Error())
		}
	}

	return doc.String(), errors
}

func updateSubContainersUsingMap(subContainers map[string]types.Node, parentID string, rootID string) error {
	// find the subcontainer whose parent == parentID
	// if there are no subcontainers we are at the end of the hierarchy, so return nil
	// otherwise, update this node and recursively call this function to process any children

	containerNode := subContainers[parentID]
	if containerNode == nil {
		return nil
	}

	idAttributeNode, err := containerNode.(types.Element).GetAttribute("id")
	if err != nil {
		return fmt.Errorf("problem accessing @id attribute for children of container with @id=\"%s\": %s", parentID, err)
	}

	// save this container's @id attribute value before deleting the idAttributeNode
	id := idAttributeNode.NodeValue()

	// FABify this container
	containerNode.(types.Element).SetAttribute("parent", rootID)
	containerNode.(types.Element).RemoveAttribute("id")

	err = updateSubContainersUsingMap(subContainers, id, rootID)
	return err
}
