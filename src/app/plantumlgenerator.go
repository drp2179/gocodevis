package app

import (
	"assembly"
	"fmt"
	"io/ioutil"
	"strings"
	"types"
)

// GeneratePlantUMLFiles - generate the plantuml files for the type collection
func GeneratePlantUMLFiles(master *assembly.MasterTypeCollection, outputFolder string, verbose bool) {
	if verbose {
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------")
	}

	if verbose {
		for _, k := range master.GetDataTypeKeys() {
			v, _ := master.GetDataType(k)
			fmt.Printf("%s\n", *v)
		}

		for _, k := range master.GetScopeKeys() {
			v, _ := master.GetScope(k)
			fmt.Printf("%s\n", v)
		}
	}

	scopeUMLBuilderMap := make(map[string]*strings.Builder)
	scopeNotesBuilderMap := make(map[string]*strings.Builder)

	//
	// generate the scope classes
	//
	for _, scopeName := range master.GetScopeKeys() {
		scope, _ := master.GetScope(scopeName)

		umlSb := &strings.Builder{}
		scopeUMLBuilderMap[scopeName] = umlSb
		notesSb := &strings.Builder{}
		scopeNotesBuilderMap[scopeName] = notesSb

		umlSb.WriteString("@startuml\n")
		umlSb.WriteString("skinparam classAttributeIconSize 0\n")

		umlSb.WriteString(fmt.Sprintf("class %s {\n", scope.Name))

		for _, scopeMethod := range scope.Methods {
			writeMethodToStringBuilder(umlSb, scopeMethod, scopeName, "static")
		}

		umlSb.WriteString("}\n")
	}

	//
	// generate the data types for each scope
	//

	for _, k := range master.GetDataTypeKeys() {
		v, _ := master.GetDataType(k)

		switch vActual := (*v).(type) {
		case *types.InterfaceDescriptor:
			scopeName := vActual.Package
			umlSb, existing := scopeUMLBuilderMap[scopeName]
			if !existing {
				fmt.Printf("Skipping UML generation for interface, no scope tracked for '%s\n", vActual.GetCannonicalName())
				continue
			}

			umlSb.WriteString(fmt.Sprintf("interface \"%s\" {\n", vActual.Name))
			for _, method := range vActual.Methods {
				writeMethodToStringBuilder(umlSb, method, scopeName, "abstract")
			}

			umlSb.WriteString("}\n")

			for _, promoted := range vActual.PromotedTypes {
				promotedName := (*promoted).GetCannonicalName()

				if strings.HasPrefix(promotedName, scopeName) {
					promotedName = (*promoted).GetSimpleName()
				}

				umlSb.WriteString(fmt.Sprintf("\"%s\" <|-- %s\n", promotedName, vActual.Name))
			}
		case *types.StructDescriptor:
			scopeName := vActual.Package

			if vActual.IsBuiltin() {
				if verbose {
					fmt.Printf("Skipping UML generation for builtin type '%s'\n", vActual.Name)
				}
				continue
			}

			umlSb, existing := scopeUMLBuilderMap[scopeName]
			if !existing {
				if verbose {
					fmt.Printf("Skipping UML generation for struct, no scope tracked for '%s'\n", vActual.GetCannonicalName())
				}
				continue
			}

			umlSb.WriteString(fmt.Sprintf("class \"%s\" {\n", vActual.Name))

			for _, attribute := range vActual.Attributes {
				writeAttributeToStringBuilder(umlSb, attribute, scopeName)
			}

			for _, method := range vActual.Methods {
				writeMethodToStringBuilder(umlSb, method, scopeName, "")
			}

			umlSb.WriteString("}\n")
			for _, promoted := range vActual.PromotedTypes {
				promotedName := (*promoted).GetCannonicalName()

				if _, _, isLocal := types.IsNameLocalToScope(promotedName, scopeName); isLocal {
					promotedName = (*promoted).GetSimpleName()
				}

				umlSb.WriteString(fmt.Sprintf("\"%s\" <|-- %s\n", promotedName, vActual.Name))
			}

			for _, attribute := range vActual.Attributes {
				attributeTypeName := attribute.CannonicalTypeName

				if resolvedName, extractedPrefix, isLocal := types.IsNameLocalToScope(attributeTypeName, scopeName); isLocal {
					attributeTypeName = resolvedName
					if extractedPrefix == types.ArrayDescriptorMarker {
						umlSb.WriteString(fmt.Sprintf("\"%s\" --* %s\n", attributeTypeName, vActual.Name))
					} else {
						umlSb.WriteString(fmt.Sprintf("\"%s\" -- %s\n", attributeTypeName, vActual.Name))
					}
				}
			}
		case *types.PointerDescriptor:
			if verbose {
				fmt.Printf("Skipping UML generation for '%s'\n", vActual.GetCannonicalName())
			}
		case *types.ArrayDescriptor:
			if verbose {
				fmt.Printf("Skipping UML generation for '%s'\n", vActual.GetCannonicalName())
			}
		case *types.MapDescriptor:
			if verbose {
				fmt.Printf("Skipping UML generation for '%s'\n", vActual.GetCannonicalName())
			}
		case *types.ChannelDescriptor:
			if verbose {
				fmt.Printf("Skipping UML generation for '%s'\n", vActual.GetCannonicalName())
			}
		case *types.FunctionDescriptor:
			if verbose {
				fmt.Printf("Skipping UML generation for '%s'\n", vActual.GetCannonicalName())
			}
		case *types.EllipsisDescriptor:
			if verbose {
				fmt.Printf("Skipping UML generation for '%s'\n", vActual.GetCannonicalName())
			}
		case *types.AliasedTypeDescriptor:
			scopeName := vActual.Package
			notesSb, existing := scopeNotesBuilderMap[scopeName]
			if !existing {
				if verbose {
					fmt.Printf("Skipping UML generation for alias, no scope tracked for '%s'\n", vActual.GetCannonicalName())
				}
				continue
			}

			notesSb.WriteString(vActual.String())
			notesSb.WriteString("\n")

		default:
			if verbose {
				fmt.Printf("skipping UML generation for '%s'\n", *v)
			}
		}
	}

	//
	// end the scopes and write the files
	//
	for _, k := range master.GetScopeKeys() {
		umlSb := scopeUMLBuilderMap[k]
		noteSb := scopeNotesBuilderMap[k]

		if noteSb.Len() > 0 {
			umlSb.WriteString(fmt.Sprintf("note as Note%s\n", k))
			umlSb.WriteString(noteSb.String())
			umlSb.WriteString("end note\n")
			umlSb.WriteString(fmt.Sprintf("%s .. Note%s\n", k, k))
		}

		umlSb.WriteString("@enduml\n")

		fileName := outputFolder + k + ".uml.txt"
		uml := umlSb.String()
		umlBytes := []byte(uml)
		err := ioutil.WriteFile(fileName, umlBytes, 0644)
		if err != nil {
			fmt.Printf("ERROR: err while writing to %s\n", fileName)
			fmt.Println(err)
		} else {
			fmt.Printf("wrote out %s\n", fileName)
		}
	}

}
func writeAttributeToStringBuilder(sb *strings.Builder, attribute *types.AttributeDescriptor, scopeName string) {
	if attribute.IsPublic() {
		sb.WriteString("+")
	} else {
		sb.WriteString("-")
	}
	sb.WriteString(" {field} ")

	attributeTypeName := attribute.CannonicalTypeName
	if resolvedName, extractedPrefix, isLocal := types.IsNameLocalToScope(attributeTypeName, scopeName); isLocal {
		attributeTypeName = extractedPrefix + resolvedName
	}

	sb.WriteString(attributeTypeName)
	sb.WriteString(" ")
	sb.WriteString(attribute.Name)
	sb.WriteString("\n")
}

func writeMethodToStringBuilder(sb *strings.Builder, method *types.MethodDescriptor, scopeName string, extraFlag string) {
	if method.IsPublic() {
		sb.WriteString("+")
	} else {
		sb.WriteString("-")
	}
	sb.WriteString(" {method} ")
	if len(extraFlag) > 0 {
		sb.WriteString(fmt.Sprintf("{%s} ", extraFlag))
	}
	sb.WriteString(method.Name)
	sb.WriteString("(")

	for i, p := range method.Params {
		if i > 0 {
			sb.WriteString(",")
		}

		paramTypeName := p.CannonicalTypeName
		if resolvedName, extractedPrefix, isLocal := types.IsNameLocalToScope(paramTypeName, scopeName); isLocal {
			paramTypeName = extractedPrefix + resolvedName
		}

		sb.WriteString(paramTypeName)
	}

	sb.WriteString(")")

	if len(method.Returns) > 0 {
		sb.WriteString(" ")
	}
	if len(method.Returns) > 1 {
		sb.WriteString("(")
	}
	for i, r := range method.Returns {
		if i > 0 {
			sb.WriteString(",")
		}

		returnTypeName := r.CannonicalTypeName
		if extractedTypeName, extractedPrefix, isLocal := types.IsNameLocalToScope(returnTypeName, scopeName); isLocal {
			returnTypeName = extractedPrefix + extractedTypeName
		}

		sb.WriteString(returnTypeName)
	}
	if len(method.Returns) > 1 {
		sb.WriteString(")")
	}

	sb.WriteString("\n")
}
