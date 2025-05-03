package models

// CodeStructure представляет структуру файла кода
type CodeStructure struct {
	// Метаданные файла
	Metadata *FileMetadata

	// Импорты файла
	Imports []*Import

	// Экспорты файла (публичные интерфейсы)
	Exports []*Export

	// Все методы/функции файла
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

	// Позиция импорта в файле
	Position Position
}

// Export представляет экспортируемый элемент
type Export struct {
	// Имя экспортируемого элемента
	Name string

	// Тип экспортируемого элемента (функция, класс, переменная и т.д.)
	Type string

	// Позиция экспортируемого элемента в файле
	Position Position
}

// Method представляет метод или функцию
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
}

// Type представляет тип или класс
type Type struct {
	// Имя типа
	Name string

	// Тип сущности (класс, интерфейс, структура и т.д.)
	Kind string

	// Является ли публичным
	IsPublic bool

	// Позиция в файле
	Position Position

	// Свойства типа
	Properties []*Property

	// Методы типа
	Methods []*Method

	// Родительский класс/тип (для наследования)
	Parent string
}

// Property представляет свойство класса или поле структуры
type Property struct {
	// Имя свойства
	Name string

	// Тип свойства
	Type string

	// Является ли публичным
	IsPublic bool

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

// AddMethod добавляет метод в структуру кода
func (cs *CodeStructure) AddMethod(method *Method) {
	cs.Methods = append(cs.Methods, method)
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
