package assembly

import (
	"fmt"
	"go/ast"
	"types"
)

// MasterTypeCollection - the master collection of types
type MasterTypeCollection struct {
	verbose         bool
	scopesMap       map[string]*types.ScopeDescriptor
	dataTypesMap    map[string]*types.TypeDescriptor
	placeholdersMap map[string]*types.PlaceholderDescriptor
}

// SetVerbose - indicate if the be verbose in logging
func (master *MasterTypeCollection) SetVerbose(verbose bool) {
	master.verbose = verbose
}

func newBuiltinType(typeName string) *types.StructDescriptor {
	dt := types.NewBuiltinStructDescriptor(typeName)

	return dt
}

func newBuiltinDataType(typeName string) *types.TypeDescriptor {
	dt := newBuiltinType(typeName)
	var builtinType types.TypeDescriptor = dt
	return &builtinType
}

// NewMasterTypeCollection - create a new master type collection
func NewMasterTypeCollection() *MasterTypeCollection {
	master := MasterTypeCollection{}
	master.scopesMap = make(map[string]*types.ScopeDescriptor)
	master.dataTypesMap = make(map[string]*types.TypeDescriptor)
	master.placeholdersMap = make(map[string]*types.PlaceholderDescriptor)

	master.dataTypesMap["string"] = newBuiltinDataType("string")
	master.dataTypesMap["bool"] = newBuiltinDataType("bool")
	master.dataTypesMap["int"] = newBuiltinDataType("int")
	master.dataTypesMap["int8"] = newBuiltinDataType("int8")
	master.dataTypesMap["int16"] = newBuiltinDataType("int16")
	master.dataTypesMap["int32"] = newBuiltinDataType("int32")
	master.dataTypesMap["int64"] = newBuiltinDataType("int64")
	master.dataTypesMap["uint"] = newBuiltinDataType("uint")
	master.dataTypesMap["uint8"] = newBuiltinDataType("uint8")
	master.dataTypesMap["uint16"] = newBuiltinDataType("uint16")
	master.dataTypesMap["uint32"] = newBuiltinDataType("uint32")
	master.dataTypesMap["uint64"] = newBuiltinDataType("uint64")
	master.dataTypesMap["uintptr"] = newBuiltinDataType("uintptr")
	master.dataTypesMap["byte"] = newBuiltinDataType("byte")
	master.dataTypesMap["rune"] = newBuiltinDataType("rune")
	master.dataTypesMap["float32"] = newBuiltinDataType("float32")
	master.dataTypesMap["float64"] = newBuiltinDataType("float64")
	master.dataTypesMap["complex64"] = newBuiltinDataType("complex64")
	master.dataTypesMap["complex128"] = newBuiltinDataType("complex128")
	master.dataTypesMap["error"] = newBuiltinDataType("error")

	// dt := types.NewStructDescriptor("", "")
	// var typeCast types.TypeDescriptor = dt
	// master.DataTypesMap[typeCast.GetCannonicalName()] = &typeCast

	// intf := types.NewInterfaceDescriptor("", "")
	// typeCast = intf
	// master.DataTypesMap[typeCast.GetCannonicalName()] = &typeCast

	return &master
}

// GetDataType - returns the TypeDescriptor for the cannonical name of the data type
func (master MasterTypeCollection) GetDataType(cannonicalName string) (*types.TypeDescriptor, bool) {
	theType, exists := master.dataTypesMap[cannonicalName]
	return theType, exists
}

// GetScope - returns the TypeDescriptor for the cannonical name of the scope
func (master MasterTypeCollection) GetScope(cannonicalName string) (*types.ScopeDescriptor, bool) {
	theType, exists := master.scopesMap[cannonicalName]
	return theType, exists
}

// GetPlaceholder - returns the TypeDescriptor for the cannonical name of the data type
func (master MasterTypeCollection) GetPlaceholder(cannonicalName string) (*types.PlaceholderDescriptor, bool) {
	theType, exists := master.placeholdersMap[cannonicalName]
	return theType, exists
}

// GetDataTypeKeys returns a list of the cannonical names in the data types map
func (master MasterTypeCollection) GetDataTypeKeys() []string {
	keys := make([]string, 0)

	for k := range master.dataTypesMap {
		keys = append(keys, k)
	}

	return keys
}

// GetScopeKeys returns a list of the cannonical names in the scope map
func (master MasterTypeCollection) GetScopeKeys() []string {
	keys := make([]string, 0)

	for k := range master.scopesMap {
		keys = append(keys, k)
	}

	return keys
}

// UpdateFromFuncDecl - update the master types using the FuncDecl
func (master *MasterTypeCollection) UpdateFromFuncDecl(packageName string, funcDecl *ast.FuncDecl) {
	var typeMethodTarget *types.TypeDescriptor
	methodName := funcDecl.Name.Name

	if funcDecl.Recv != nil {
		_, receiverType := master.extractReceiverTypeFromFieldList(packageName, funcDecl.Recv)
		typeMethodTarget = receiverType
	} else {
		scope := master.scopesMap[packageName]

		if scope == nil {
			scope = types.NewScopeDescriptor(packageName)
			master.addNewScopeType(scope)
			if master.verbose {
				fmt.Printf("added func dec scope '%s'\n", packageName)
			}
		}

		var scopeCast types.TypeDescriptor = scope
		typeMethodTarget = &scopeCast
	}

	params := master.extractAttributesFromFieldList(packageName, funcDecl.Type.Params)
	results := master.extractAttributesFromFieldList(packageName, funcDecl.Type.Results)
	method := types.NewMethodDescriptor(methodName, params, results)

	switch switchedMethodTarget := (*typeMethodTarget).(type) {
	case types.TypeWithMethods:
		switchedMethodTarget.AddMethod(method)

		if master.verbose {
			fmt.Printf("added funcdef method '%s' to '%s' \n", method.Name, switchedMethodTarget.GetSimpleName())
		}
	case *types.PointerDescriptor:
		switch ptrTarget := (*(switchedMethodTarget.RefType)).(type) {
		case types.TypeWithMethods:
			ptrTarget.AddMethod(method)
			if master.verbose {
				fmt.Printf("added funcdef method '%s' to '%s' \n", method.Name, ptrTarget.GetSimpleName())
			}
		default:
			fmt.Printf("ERROR: trying to add method '%s' to a non-method pointer target type target: '%s'\n", methodName, *typeMethodTarget)
		}
	default:
		fmt.Printf("ERROR: trying to add method '%s' to a non-method type target: '%s'\n", methodName, *typeMethodTarget)
	}
}

// UpdateFromTypeSpec - give a type spec, assemble the class definitions
func (master *MasterTypeCollection) UpdateFromTypeSpec(packageName string, typeSpec *ast.TypeSpec) {

	assemblyScope, assemblyScopeExists := master.scopesMap[packageName]
	if !assemblyScopeExists {
		assemblyScope = types.NewScopeDescriptor(packageName)
		master.addNewScopeType(assemblyScope)
	}

	switch theTypeData := typeSpec.Type.(type) {
	case *ast.Ident:
		typeName := typeSpec.Name.Name
		_, aliasedType := master.ExtractIdentType(packageName, theTypeData)
		master.addAliasedDataType(packageName, typeName, aliasedType)
	case *ast.SelectorExpr:
		typeName := typeSpec.Name.Name
		_, aliasedType := master.ExtractSelectorType(theTypeData)
		master.addAliasedDataType(packageName, typeName, aliasedType)
	case *ast.ArrayType:
		typeName := typeSpec.Name.Name
		_, aliasedType := master.extractTypeFromExpr(packageName, theTypeData.Elt)
		master.addAliasedDataType(packageName, typeName, aliasedType)
	case *ast.StructType:
		structName := typeSpec.Name.Name
		var structType types.TypeDescriptor

		dt := types.NewStructDescriptor(packageName, structName)

		cannonicalName := dt.GetCannonicalName()
		if existingDt, exists := master.dataTypesMap[cannonicalName]; exists {
			structType = *existingDt
		} else {
			structType = dt
			master.addNewDataType(&structType)
		}

		var structTypeWithAttributes types.TypeWithAttributes

		switch actual := structType.(type) {
		case types.TypeWithAttributes:
			structTypeWithAttributes = actual
			master.updateStructDescriptorAttributesFromFieldList(packageName, &structTypeWithAttributes, theTypeData.Fields)
		default:
			fmt.Printf("ERROR: UpdateFromTypeSpec is trying to add attributes to '%s' a non attribute type '%T'\n", actual.GetCannonicalName(), actual)
		}

	case *ast.InterfaceType:
		interfaceName := typeSpec.Name.Name

		intfc := types.NewInterfaceDescriptor(packageName, interfaceName)
		master.updateInterfaceFromFieldList(intfc, theTypeData.Methods)

		cannonicalName := intfc.GetCannonicalName()
		if existingDataType, exists := master.dataTypesMap[cannonicalName]; exists {
			switch (*existingDataType).(type) {
			case *types.StructDescriptor:
				var typeCast types.TypeDescriptor = intfc
				master.replaceDataType(&typeCast)
			case *types.InterfaceDescriptor:
				fmt.Printf("ERROR: already have an interface definition for '%s'\n", interfaceName)
			default:
				fmt.Printf("ERROR: already have an %T definition for '%s'\n", *existingDataType, interfaceName)
			}
		} else {
			var newDataType types.TypeDescriptor = intfc
			master.addNewDataType(&newDataType)
		}
	case *ast.FuncType:
		aliasName := typeSpec.Name.Name

		fd := types.NewFunctionDescriptor(packageName, "")
		master.updateFunctionDescriptorParamsFromFieldList(fd, theTypeData.Params)
		master.updateFunctionDescriptorReturnsFromFieldList(fd, theTypeData.Results)

		functionTypeName := fd.GetCannonicalName()
		if _, exists := master.dataTypesMap[functionTypeName]; exists {
			//fmt.Printf("Already have definition for class '%s', how is this possible?\n", cannonicalName)
		} else {
			var newDataType types.TypeDescriptor = fd
			master.addNewDataType(&newDataType)
		}

		aliasedType := master.dataTypesMap[functionTypeName]
		master.addAliasedDataType(packageName, aliasName, aliasedType)
	case *ast.MapType:
		aliasName := typeSpec.Name.Name
		keyTypeName, keyType := master.extractTypeFromExpr(packageName, theTypeData.Key)
		valueTypeName, valueType := master.extractTypeFromExpr(packageName, theTypeData.Value)

		mappedTypeName := types.GetMapTypeNameForTypenames(keyTypeName, valueTypeName)

		var mappedType *types.TypeDescriptor
		var mappedTypeExists bool
		if mappedType, mappedTypeExists = master.dataTypesMap[mappedTypeName]; mappedTypeExists {
			// nothing to do
		} else {
			mapDescriptor := types.NewMapDescriptor(keyType, valueType)
			var mapDescriptorCast types.TypeDescriptor = mapDescriptor
			mappedType = &mapDescriptorCast
			master.addNewDataType(mappedType)
		}

		master.addAliasedDataType(packageName, aliasName, mappedType)
	default:
		fmt.Printf("WARN: ignoring TypeSpec type '%s' of '%s'\n", typeSpec.Name.Name, theTypeData)
	}
}

func (master *MasterTypeCollection) addAliasedDataType(packageName string, aliasName string, aliasedDataType *types.TypeDescriptor) {

	aliasedDataTypeDescriptor := types.NewAliasedTypeDescriptor(packageName, aliasName, aliasedDataType)

	var aliasCast types.TypeDescriptor = aliasedDataTypeDescriptor
	master.addNewDataType(&aliasCast)
	if master.verbose {
		fmt.Printf("added aliased data type '%s' for '%s'\n", aliasedDataTypeDescriptor.GetCannonicalName(), (*aliasedDataType).GetCannonicalName())
	}
}

func (master *MasterTypeCollection) replaceDataType(dt *types.TypeDescriptor) {

	name := (*dt).GetCannonicalName()

	if existingDt, exists := master.dataTypesMap[name]; exists {
		master.dataTypesMap[name] = dt
		if master.verbose {
			fmt.Printf("replaced exist data type '%s' of '%T' with '%T'\n", name, *existingDt, *dt)
		}
	} else {
		fmt.Printf("ERROR: asked to replace '%s' which doesn't exist\n", name)
	}
}

func (master *MasterTypeCollection) addNewDataType(newType *types.TypeDescriptor) {
	name := (*newType).GetCannonicalName()

	master.dataTypesMap[name] = newType

	if master.verbose {
		fmt.Printf("added new data type '%s' of '%T'\n", name, *newType)
	}

	if existingPlaceholder, placeholderExists := master.placeholdersMap[name]; placeholderExists {
		delete(master.placeholdersMap, name)

		switch newTypeWithMethods := (*newType).(type) {
		case types.TypeWithMethods:
			for _, m := range existingPlaceholder.GetMethods() {
				newTypeWithMethods.AddMethod(m)
			}
		}

		if master.verbose {
			fmt.Printf("removed placeholder for data type '%s'\n", name)
		}

		for k, v := range master.dataTypesMap {
			switch vType := (*v).(type) {
			case *types.ArrayDescriptor:
				vType.ReplaceRefTypeIfPlaceholder(newType)
			case *types.PointerDescriptor:
				vType.ReplaceRefTypeIfPlaceholder(newType)
			case *types.EllipsisDescriptor:
				vType.ReplaceRefTypeIfPlaceholder(newType)
			case *types.ChannelDescriptor:
				vType.ReplaceRefTypeIfPlaceholder(newType)
			case *types.MapDescriptor:
				vType.ReplaceKeyTypeIfPlaceholder(newType)
				vType.ReplaceValueTypeIfPlaceholder(newType)
			case *types.FunctionDescriptor:
				// nothing to do here
			case *types.StructDescriptor:
				// nothing to do here
			case *types.InterfaceDescriptor:
				// nothing to do here
			case *types.AliasedTypeDescriptor:
				// nothing to do here
			default:
				fmt.Printf("WARN: ignoring replacing placeholder for data type '%s' '%s'\n", k, vType)
			}
		}
	}
}

func (master *MasterTypeCollection) addNewScopeType(newType *types.ScopeDescriptor) {
	name := (*newType).GetCannonicalName()

	master.scopesMap[name] = newType

	if master.verbose {
		fmt.Printf("added new scope type '%s'\n", name)
	}

	if _, placeholderExists := master.placeholdersMap[name]; placeholderExists {
		delete(master.placeholdersMap, name)

		if master.verbose {
			fmt.Printf("removed placeholder for scope type '%s'\n", name)
		}
	}
}

func (master *MasterTypeCollection) extractTypeFromExpr(defaultPackageName string, typeExpression ast.Expr) (string, *types.TypeDescriptor) {
	var theType *types.TypeDescriptor
	var name string

	switch expression := typeExpression.(type) {
	case *ast.Ident:
		name, theType = master.ExtractIdentType(defaultPackageName, expression)
	case *ast.StarExpr:
		name, theType = master.ExtractPointerType(defaultPackageName, expression)
	case *ast.ArrayType:
		name, theType = master.ExtractArrayType(defaultPackageName, expression)
	case *ast.MapType:
		name, theType = master.ExtractMapType(defaultPackageName, expression)
	case *ast.ChanType:
		name, theType = master.ExtractChannelType(defaultPackageName, expression)
	case *ast.Ellipsis:
		name, theType = master.ExtractEllipsis(defaultPackageName, expression)
	case *ast.SelectorExpr:
		name, theType = master.ExtractSelectorType(expression)
	case *ast.FuncType:
		name, theType = master.ExtractFuncType(defaultPackageName, expression)
	case *ast.InterfaceType:
		name, theType = master.ExtractInterfaceType(expression)
	case *ast.StructType:
		name, theType = master.ExtractStructType(defaultPackageName, expression)
	default:
		fmt.Printf("ERROR: Unknown typeExpression '%s' in extractTypeFromExpr\n", typeExpression)
	}

	return name, theType
}

func (master *MasterTypeCollection) getOrCreatePlaceholderForIdent(defaultPackageName string, ident *ast.Ident) (string, *types.TypeDescriptor) {

	placeholder := types.NewPlaceholderDescriptor(defaultPackageName, ident.Name)
	theName := placeholder.GetCannonicalName()
	var theType types.TypeDescriptor

	if existingPlaceholderType, placeholderExists := master.placeholdersMap[theName]; placeholderExists {
		theType = existingPlaceholderType
	} else {
		master.placeholdersMap[theName] = placeholder
		theType = placeholder

		if master.verbose {
			fmt.Printf("added placeholder for '%s'\n", theName)
		}
	}

	return theName, &theType
}

func (master *MasterTypeCollection) extractAttributesFromFieldList(defaultPackageName string, fieldList *ast.FieldList) []*types.AttributeDescriptor {

	attributes := make([]*types.AttributeDescriptor, 0)

	if fieldList != nil {
		for _, field := range fieldList.List {
			extractedName, extractedType := master.extractTypeFromExpr(defaultPackageName, field.Type)

			if extractedType == nil && len(extractedName) <= 0 {
				fmt.Printf("WARN: no type name or type extracted from '%s', so not adding any attribute details, likely a function pointer\n", field.Type)
				continue
			} else if extractedType == nil {
				dt := types.NewStructDescriptor(defaultPackageName, extractedName)
				var dataTypesCast types.TypeDescriptor = dt
				extractedType = &dataTypesCast
				master.addNewDataType(extractedType)
			}

			var name string
			if len(field.Names) > 0 {
				name = field.Names[0].Name
			}
			var tag string

			if field.Tag != nil {
				tag = field.Tag.Value
			}

			attribute := types.NewAttributeDescriptor(name, *extractedType, tag)

			attributes = append(attributes, attribute)
		}
	}

	return attributes
}

func (master *MasterTypeCollection) extractReceiverTypeFromFieldList(defaultPackageName string, fieldList *ast.FieldList) (string, *types.TypeDescriptor) {
	var receiverTypeName string
	var receiverType *types.TypeDescriptor

	if fieldList != nil && fieldList.NumFields() > 0 {
		field := fieldList.List[0]

		extractedName, extractedType := master.extractTypeFromExpr(defaultPackageName, field.Type)
		if extractedType == nil {
			fmt.Printf("ERROR: Unable to extract receiver type from field list")
		}

		receiverTypeName = extractedName
		receiverType = extractedType
	}

	return receiverTypeName, receiverType
}

// ExtractFieldType - returns the TypeDecriptor for the provided Field and default package
func (master *MasterTypeCollection) ExtractFieldType(defaultPackageName string, field *ast.Field) (string, *types.TypeDescriptor) {
	return master.extractTypeFromExpr(defaultPackageName, field.Type)
}

func (master *MasterTypeCollection) updateFunctionDescriptorParamsFromFieldList(fd *types.FunctionDescriptor, fieldList *ast.FieldList) {
	attributes := master.extractAttributesFromFieldList(fd.Package, fieldList)

	for _, a := range attributes {
		fd.AddParameter(a)
	}
}

func (master *MasterTypeCollection) updateFunctionDescriptorReturnsFromFieldList(fd *types.FunctionDescriptor, fieldList *ast.FieldList) {
	attributes := master.extractAttributesFromFieldList(fd.Package, fieldList)

	for _, a := range attributes {
		fd.AddReturn(a)
	}
}

func (master *MasterTypeCollection) updateStructDescriptorAttributesFromFieldList(defaultPackage string, dt *types.TypeWithAttributes, fieldList *ast.FieldList) {

	if fieldList != nil && fieldList.NumFields() > 0 {
		for _, field := range fieldList.List {
			fieldName := ExtractFieldName(field)

			attribute := types.AttributeDescriptor{}
			attribute.Name = fieldName

			fieldTypeName, fieldType := master.ExtractFieldType(defaultPackage, field)
			if fieldType == nil {
				newFieldDataType := types.NewStructDescriptor(defaultPackage, fieldTypeName)
				var dataTypesCast types.TypeDescriptor = newFieldDataType
				fieldType = &dataTypesCast
				master.addNewDataType(fieldType)
			}

			if len(fieldName) == 0 {
				(*dt).AddPromotedType(fieldType)
			} else {
				attribute.CannonicalTypeName = (*fieldType).GetCannonicalName()
				(*dt).AddAttribute(&attribute)
				if master.verbose {
					fmt.Printf("added attribute '%s' typed as '%s' to '%s'\n", attribute.Name, attribute.CannonicalTypeName, (*dt).GetCannonicalName())
				}
			}
		}
	}
}

func (master *MasterTypeCollection) updateInterfaceFromFieldList(intfc *types.InterfaceDescriptor, methods *ast.FieldList) {
	if methods != nil && methods.NumFields() > 0 {
		for _, field := range methods.List {

			functionName := ExtractFieldName(field)

			switch funcField := field.Type.(type) {
			case *ast.FuncType:
				params := make([]*types.AttributeDescriptor, 0)
				returns := make([]*types.AttributeDescriptor, 0)

				extractedParams := master.extractAttributesFromFieldList(intfc.Package, funcField.Params)
				for _, p := range extractedParams {
					params = append(params, p)
				}
				extractedReturns := master.extractAttributesFromFieldList(intfc.Package, funcField.Results)
				for _, r := range extractedReturns {
					returns = append(returns, r)
				}

				method := types.NewMethodDescriptor(functionName, params, returns)
				intfc.AddMethod(method)

				if master.verbose {
					fmt.Printf("added funcdef method '%s' to '%s' \n", method.Name, intfc.GetSimpleName())
				}

			case *ast.Ident:
				_, promotingType := master.ExtractIdentType(intfc.Package, funcField)
				intfc.AddPromotedType(promotingType)
			case *ast.SelectorExpr:
				_, promotingType := master.ExtractSelectorType(funcField)
				intfc.AddPromotedType(promotingType)
			default:
				fmt.Printf("ERROR: updating interface with method list that isn't a function type\n")
			}
		}
	}
}

// ExtractIdentType returns the name of the type and the created type from an ast.Ident. if the type already exists return that, else create a placeholder
func (master *MasterTypeCollection) ExtractIdentType(defaultPackageName string, ident *ast.Ident) (string, *types.TypeDescriptor) {
	var theType *types.TypeDescriptor
	name := ident.Name

	if existingDtType, dtExists := master.dataTypesMap[name]; dtExists {
		theType = existingDtType
	} else if len(defaultPackageName) > 0 {
		longerName := defaultPackageName + "." + name
		if existingDtType, dtExists = master.dataTypesMap[longerName]; dtExists {
			name = longerName
			theType = existingDtType
		}
	}

	if theType == nil {
		if scopeType, scopeExists := master.scopesMap[name]; scopeExists {
			var scopeCast types.TypeDescriptor = scopeType
			theType = &scopeCast
		} else {
			name, theType = master.getOrCreatePlaceholderForIdent(defaultPackageName, ident)
		}
	}

	return name, theType
}

// ExtractPointerType - extract a PointerDescriptor from the provided StarExpr and package name
func (master *MasterTypeCollection) ExtractPointerType(defaultPackageName string, expression *ast.StarExpr) (string, *types.TypeDescriptor) {
	var pointerType *types.TypeDescriptor
	ptrName, pointerRefType := master.extractTypeFromExpr(defaultPackageName, expression.X)

	if pointerRefType == nil {
		switch xType := expression.X.(type) {
		case *ast.Ident:
			ptrName, pointerRefType = master.getOrCreatePlaceholderForIdent(defaultPackageName, xType)
		default:
			fmt.Printf("ERROR: can't find the type for '%s' for pointer type\n", expression.X)
		}
	}

	if pointerRefType != nil {
		ptrName = types.GetPointerTypeNameForTypename(ptrName)
		existingPtrType, ptrTypeExists := master.dataTypesMap[ptrName]

		if ptrTypeExists {
			pointerType = existingPtrType
		} else {
			ptrDescriptor := types.NewPointerDescriptor(pointerRefType)
			var ptrCast types.TypeDescriptor = ptrDescriptor
			pointerType = &ptrCast
			master.addNewDataType(pointerType)
		}
	}

	return ptrName, pointerType
}

// ExtractArrayType - extact an ArrayDescriptor from the provided ArrayType expression and default package
func (master *MasterTypeCollection) ExtractArrayType(defaultPackageName string, expression *ast.ArrayType) (string, *types.TypeDescriptor) {
	var arrayTypeName string
	var arrayType *types.TypeDescriptor

	arrayRefTypeName, arrayRefType := master.extractTypeFromExpr(defaultPackageName, expression.Elt)

	if arrayRefType == nil {
		fmt.Printf("ERROR: can't find the type for '%s' for array type\n", expression.Elt)
	} else {
		arrayTypeName = types.GetArrayTypeNameForTypename(arrayRefTypeName)

		if existingArrayType, arrayTypeExists := master.dataTypesMap[arrayTypeName]; arrayTypeExists {
			arrayType = existingArrayType
		} else {
			arrayDescriptor := types.NewArrayDescriptor(arrayRefType)
			var arrayCast types.TypeDescriptor = arrayDescriptor
			arrayType = &arrayCast

			master.addNewDataType(arrayType)
		}
	}

	return arrayTypeName, arrayType
}

// ExtractMapType - extract a MapDescriptor from the MapType expression and the default package
func (master *MasterTypeCollection) ExtractMapType(defaultPackageName string, expression *ast.MapType) (string, *types.TypeDescriptor) {
	var mapTypeName string
	var mapType *types.TypeDescriptor

	keyRefTypeName, keyRefType := master.extractTypeFromExpr(defaultPackageName, expression.Key)
	valueRefTypeName, valueRefType := master.extractTypeFromExpr(defaultPackageName, expression.Value)

	if keyRefType == nil {
		fmt.Printf("ERROR: can't find the key type for '%s' for map type\n", expression.Key)
	} else if valueRefType == nil {
		fmt.Printf("ERROR: can't find the value type for '%s' for map type\n", expression.Value)
	} else {
		mapTypeName = types.GetMapTypeNameForTypenames(keyRefTypeName, valueRefTypeName)

		if existingMapType, mapTypeExists := master.dataTypesMap[mapTypeName]; mapTypeExists {
			mapType = existingMapType
		} else {
			mapDescriptor := types.NewMapDescriptor(keyRefType, valueRefType)
			var mapCast types.TypeDescriptor = mapDescriptor
			mapType = &mapCast

			master.addNewDataType(mapType)
		}
	}

	return mapTypeName, mapType
}

// ExtractChannelType - extract a ChannelDescriptor from the provided ChanType expression and the default package
func (master *MasterTypeCollection) ExtractChannelType(defaultPackageName string, expression *ast.ChanType) (string, *types.TypeDescriptor) {
	var chanTypeName string
	var chanType *types.TypeDescriptor
	channelDir := types.ChannelDirectionTx
	if expression.Dir&ast.RECV > 0 && expression.Dir&ast.SEND > 0 {
		channelDir = types.ChannelDirectionNone
	} else if expression.Dir&ast.RECV > 0 {
		channelDir = types.ChannelDirectionRx
	} else {
		channelDir = types.ChannelDirectionTx
	}

	channelOfTypeName, channelOfRefType := master.extractTypeFromExpr(defaultPackageName, expression.Value)

	if channelOfRefType == nil {
		fmt.Printf("ERROR: can't find the type for '%s' for channel type\n", expression.Value)
	} else {
		chanTypeName = types.GetChannelTypeNameForTypename(channelDir, channelOfTypeName)

		if existingChannelType, chanTypeExists := master.dataTypesMap[chanTypeName]; chanTypeExists {
			chanType = existingChannelType
		} else {
			channelDescriptor := types.NewChannelDescriptor(channelDir, channelOfRefType)
			var channelCast types.TypeDescriptor = channelDescriptor
			chanType = &channelCast

			master.addNewDataType(chanType)
		}
	}

	return chanTypeName, chanType
}

// ExtractEllipsis - extract EllipsisDescriptor from the provided Ellipsis expression and default package
func (master *MasterTypeCollection) ExtractEllipsis(defaultPackageName string, expression *ast.Ellipsis) (string, *types.TypeDescriptor) {
	var ellipsisTypeName string
	var ellipsisType *types.TypeDescriptor
	ellipsisRefTypeName, ellipsisRefType := master.extractTypeFromExpr(defaultPackageName, expression.Elt)

	if ellipsisRefType == nil {
		fmt.Printf("ERROR: can't find the type for '%s' for ellipsis ref type\n", expression.Elt)
	} else {
		ellipsisTypeName = types.GetEllipsisTypeNameForTypename(ellipsisRefTypeName)

		if existingEllipsisType, ellipsisTypeExists := master.dataTypesMap[ellipsisTypeName]; ellipsisTypeExists {
			ellipsisType = existingEllipsisType
		} else {
			ellipsisDescriptor := types.NewEllipsisDescriptor(ellipsisRefType)
			var ellipsisCast types.TypeDescriptor = ellipsisDescriptor
			ellipsisType = &ellipsisCast

			master.addNewDataType(ellipsisType)
		}
	}

	return ellipsisTypeName, ellipsisType
}

// ExtractSelectorType - returns the existing TypeDescriptor for the provided expression, or creates the scope if necessary
// and the associated structure
func (master *MasterTypeCollection) ExtractSelectorType(expression *ast.SelectorExpr) (string, *types.TypeDescriptor) {
	var selectorType *types.TypeDescriptor
	var selectorTypeName string

	scopeName, scopeType := master.extractTypeFromExpr("", expression.X)
	if scopeType == nil {
		// unknown scope, so add it
		newAssemblyScope := types.NewScopeDescriptor(scopeName)
		master.addNewScopeType(newAssemblyScope)
		var assemblyCast types.TypeDescriptor = newAssemblyScope
		scopeType = &assemblyCast
	}

	className := expression.Sel.Name
	selectorTypeName = (*scopeType).GetCannonicalName() + "." + className
	dtType, classExists := master.dataTypesMap[selectorTypeName]
	if classExists {
		selectorType = dtType
	} else {
		packageName := (*scopeType).GetCannonicalName()
		theDataType := types.NewStructDescriptor(packageName, className)
		var dataTypeCast types.TypeDescriptor = theDataType
		selectorType = &dataTypeCast
		master.addNewDataType(selectorType)
	}

	return selectorTypeName, selectorType
}

//ExtractFuncType - returns FunctionDescriptor for the provided FuncType expression and package name
func (master *MasterTypeCollection) ExtractFuncType(defaultPackageName string, expression *ast.FuncType) (string, *types.TypeDescriptor) {
	var funcType *types.TypeDescriptor
	var funcTypeName string

	tmpFunDef := types.NewFunctionDescriptor(defaultPackageName, "")

	master.updateFunctionDescriptorParamsFromFieldList(tmpFunDef, expression.Params)
	master.updateFunctionDescriptorReturnsFromFieldList(tmpFunDef, expression.Results)

	funcTypeName = tmpFunDef.GetCannonicalName()
	if existingType, typeExists := master.dataTypesMap[funcTypeName]; typeExists {
		funcType = existingType
	} else {
		var typeCast types.TypeDescriptor = tmpFunDef
		funcType = &typeCast
		master.addNewDataType(funcType)
	}

	return funcTypeName, funcType
}

// ExtractInterfaceType - returns an anonymous InterfaceDescriptor for the provided InterfaceType expression
func (master *MasterTypeCollection) ExtractInterfaceType(expression *ast.InterfaceType) (string, *types.TypeDescriptor) {
	var interfaceTypeName string
	var interfaceType types.TypeDescriptor

	tmpInteraface := types.NewInterfaceDescriptor("", "")

	master.updateInterfaceFromFieldList(tmpInteraface, expression.Methods)

	interfaceTypeName = tmpInteraface.GetCannonicalName()

	// if existingInterface, interfaceExists := master.DataTypesMap[interfaceTypeName]; interfaceExists {
	// 	interfaceType = existingInterface
	// } else {
	// 	var typeCast types.TypeDescriptor = tmpInteraface
	// 	interfaceType = &typeCast
	// 	master.addNewDataType(interfaceType)
	// }

	interfaceType = tmpInteraface

	return interfaceTypeName, &interfaceType
}

// ExtractStructType - returns an anonoymous StructDescriptor for the provided StructType expression and default package
func (master *MasterTypeCollection) ExtractStructType(defaultPackageName string, expression *ast.StructType) (string, *types.TypeDescriptor) {
	var structTypeName string
	var structType types.TypeDescriptor

	dt := types.NewStructDescriptor("", "")

	structTypeName = dt.GetCannonicalName()

	// if existingType, typeExists := master.DataTypesMap[structTypeName]; typeExists {
	// 	structType = *existingType
	// } else {
	// 	structType = dt
	// 	master.addNewDataType(&structType)
	// }
	structType = dt
	var structTypeWithAttributes types.TypeWithAttributes

	switch actual := structType.(type) {
	case types.TypeWithAttributes:
		structTypeWithAttributes = actual
		master.updateStructDescriptorAttributesFromFieldList(defaultPackageName, &structTypeWithAttributes, expression.Fields)
	default:
		fmt.Printf("ERROR: extractStructType is trying to add attributes to '%s' a non attribute type '%T\n", actual.GetCannonicalName(), actual)
	}

	return structTypeName, &structType
}
