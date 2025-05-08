package models

// CodeStructure представляет структуру файла кода
type CodeStructure struct {
	// Метаданные файла
	Metadata *FileMetadata

	// Импорты файла
	Imports []*Import

	// Экспорты файла (публичные интерфейсы)
	Exports []*Export

	// Функции верхнего уровня
	Functions []*Function

	// Методы классов
	Methods []*Method

	// Типы/классы, определенные в файле
	Types []*Type

	// Переменные верхнего уровня
	Variables []*Variable

	// Константы
	Constants []*Constant
}

// Import представляет импорт в файле
type Import struct {
	// Путь импорта
	Path string

	// Псевдоним импорта (если есть)
	Alias string

	// Является ли импорт динамическим (import())
	IsDynamic bool

	// Является ли импорт namespace-импортом
	IsNamespace bool

	// Является ли импорт type-импортом (TypeScript)
	IsTypeImport bool

	// Позиция импорта в файле
	Position Position
}

// Export представляет экспортируемый элемент
type Export struct {
	// Имя экспортируемого элемента
	Name string

	// Тип экспортируемого элемента (function, class, variable и т.д.)
	Type string

	// Является ли экспорт default-экспортом
	IsDefault bool

	// Является ли экспорт type-экспортом (TypeScript)
	IsTypeExport bool

	// Является ли экспорт namespace-экспортом
	IsNamespace bool

	// Позиция экспортируемого элемента в файле
	Position Position
}

// Function представляет функцию верхнего уровня
type Function struct {
	// Имя функции
	Name string

	// Параметры функции
	Parameters []*Parameter

	// Тип возвращаемого значения
	ReturnType string

	// Является ли функция публичной
	IsPublic bool

	// Является ли функция асинхронной (async)
	IsAsync bool

	// Является ли функция генератором (function*)
	IsGenerator bool

	// Является ли функция стрелочной функцией
	IsArrow bool

	// Является ли функция IIFE (Immediately Invoked Function Expression)
	IsIIFE bool

	// Позиция функции в файле
	Position Position

	// Описание функции
	Description string
}

// Method представляет метод класса
type Method struct {
	// Имя метода
	Name string

	// Параметры метода
	Parameters []*Parameter

	// Тип возвращаемого значения
	ReturnType string

	// Является ли метод публичным
	IsPublic bool

	// Является ли метод статическим
	IsStatic bool

	// Является ли метод асинхронным (async)
	IsAsync bool

	// Является ли метод генератором (function*)
	IsGenerator bool

	// Является ли метод декоратором
	IsDecorator bool

	// Является ли метод конструктором
	IsConstructor bool

	// Тип метода (method, getter, setter)
	Kind string

	// Позиция метода в файле
	Position Position

	// Описание метода
	Description string

	// Принадлежность к классу/типу
	BelongsTo string
}

// Parameter представляет параметр метода или функции
type Parameter struct {
	// Имя параметра
	Name string

	// Тип параметра
	Type string

	// Значение по умолчанию (если есть)
	DefaultValue string

	// Является ли параметр обязательным
	IsRequired bool

	// Является ли параметр вариативным (rest parameter)
	IsVariadic bool

	// Является ли параметр деструктуризированным объектом
	IsDestructuredObject bool

	// Является ли параметр деструктуризированным массивом
	IsDestructuredArray bool
}

// Type представляет тип или класс
type Type struct {
	// Имя типа
	Name string

	// Тип сущности (class, interface, struct, enum, etc.)
	Kind string

	// Является ли публичным
	IsPublic bool

	// Является ли абстрактным классом
	IsAbstract bool

	// Является ли интерфейсом
	IsInterface bool

	// Является ли миксином
	IsMixin bool

	// Является ли дженериком
	IsGeneric bool

	// Является ли перечислением (enum)
	IsEnum bool

	// Позиция в файле
	Position Position

	// Свойства типа
	Properties []*Property

	// Методы типа
	Methods []*Method

	// Родительский класс/тип (для наследования)
	Parent string

	// Реализуемые интерфейсы
	Implements []string

	// Дженерик параметры
	GenericParameters []string
}

// Property представляет свойство класса или поле структуры
type Property struct {
	// Имя свойства
	Name string

	// Тип свойства
	Type string

	// Является ли публичным
	IsPublic bool

	// Является ли статическим
	IsStatic bool

	// Является ли вычисляемым свойством
	IsComputed bool

	// Является ли приватным полем
	IsPrivate bool

	// Является ли readonly
	IsReadonly bool

	// Позиция в файле
	Position Position
}

// Variable представляет переменную
type Variable struct {
	// Имя переменной
	Name string

	// Тип переменной
	Type string

	// Является ли публичной
	IsPublic bool

	// Позиция в файле
	Position Position
}

// Constant представляет константу
type Constant struct {
	// Имя константы
	Name string

	// Тип константы
	Type string

	// Значение константы
	Value string

	// Позиция в файле
	Position Position
}

// Position представляет позицию в файле
type Position struct {
	// Начальная строка
	StartLine int

	// Начальная колонка
	StartColumn int

	// Конечная строка
	EndLine int

	// Конечная колонка
	EndColumn int
}

// NewCodeStructure создает новую структуру кода для файла
func NewCodeStructure(metadata *FileMetadata) *CodeStructure {
	return &CodeStructure{
		Metadata:  metadata,
		Imports:   make([]*Import, 0),
		Exports:   make([]*Export, 0),
		Functions: make([]*Function, 0),
		Methods:   make([]*Method, 0),
		Types:     make([]*Type, 0),
		Variables: make([]*Variable, 0),
		Constants: make([]*Constant, 0),
	}
}

// AddImport добавляет импорт в структуру кода
func (cs *CodeStructure) AddImport(imp *Import) {
	cs.Imports = append(cs.Imports, imp)
}

// AddExport добавляет экспорт в структуру кода
func (cs *CodeStructure) AddExport(exp *Export) {
	cs.Exports = append(cs.Exports, exp)
}

// AddFunction добавляет функцию в структуру кода
func (cs *CodeStructure) AddFunction(fn *Function) {
	cs.Functions = append(cs.Functions, fn)
}

// AddType добавляет тип в структуру кода
func (cs *CodeStructure) AddType(typ *Type) {
	cs.Types = append(cs.Types, typ)
}

// AddVariable добавляет переменную в структуру кода
func (cs *CodeStructure) AddVariable(variable *Variable) {
	cs.Variables = append(cs.Variables, variable)
}

// AddConstant добавляет константу в структуру кода
func (cs *CodeStructure) AddConstant(constant *Constant) {
	cs.Constants = append(cs.Constants, constant)
}

// AddMethod добавляет метод в структуру кода
func (cs *CodeStructure) AddMethod(method *Method) {
	cs.Methods = append(cs.Methods, method)
}

// GetPublicFunctions возвращает только публичные функции
func (cs *CodeStructure) GetPublicFunctions() []*Function {
	var publicFunctions []*Function
	for _, fn := range cs.Functions {
		if fn.IsPublic {
			publicFunctions = append(publicFunctions, fn)
		}
	}
	return publicFunctions
}

// GetPublicMethods возвращает только публичные методы
func (cs *CodeStructure) GetPublicMethods() []*Method {
	var publicMethods []*Method
	for _, method := range cs.Methods {
		if method.IsPublic {
			publicMethods = append(publicMethods, method)
		}
	}
	return publicMethods
}

// GetPublicTypes возвращает только публичные типы
func (cs *CodeStructure) GetPublicTypes() []*Type {
	var publicTypes []*Type
	for _, typ := range cs.Types {
		if typ.IsPublic {
			publicTypes = append(publicTypes, typ)
		}
	}
	return publicTypes
}
