package parser

// JSQuery 是针对 JavaScript/TypeScript 的 Tree-sitter 查询语句
// 注意：Tree-sitter Query 语言中的注释使用分号 (;)
const JSQuery = `
; 1. 捕获类的声明
(class_declaration name: (identifier) @class.name)

; 2. 捕获函数的声明
(function_declaration name: (identifier) @function.name)

; 3. 捕获类里面的方法定义
(method_definition name: (property_identifier) @method.name)

; 4. 捕获函数调用
(call_expression
  function: [
    (identifier) @call.function
    (member_expression property: (property_identifier) @call.method)
  ])

; 5. 捕获导入依赖
(import_statement source: (string) @import.source)
`

// PYQuery 是针对 Python 的 Tree-sitter 查询语句
const PYQuery = `
; 1. 捕获类的声明 (Class Definitions)
(class_definition name: (identifier) @class.name)

; 2. 捕获函数的声明 (Function Definitions)
(function_definition name: (identifier) @function.name)

; 3. Python 里的方法定义其实也是 function_definition 嵌套在 block 里的，
; 为了简化和 Python AST 的兼容性，Python 原版也会把方法统一视为 function_definition。
; 但如果需要专门区分，可以在 AST 层提取父节点是否是 class_definition。
; 在这里我们暂且不写单独的 @method.name 捕获，统一走 @function.name

; 4. 捕获函数调用 (Call)
(call function: [
    (identifier) @call.function
    (attribute attribute: (identifier) @call.method)
  ])

; 5. 捕获导入依赖 (Import & Import From)
(import_statement name: (dotted_name) @import.source)
(import_from_statement module_name: (dotted_name) @import.source)
`

// GOQuery 是针对 Go 语言的 Tree-sitter 查询语句
const GOQuery = `
; 1. 捕获结构体/类型的声明 (Type Declarations)
(type_declaration (type_spec name: (type_identifier) @class.name))

; 2. 捕获普通函数的声明
(function_declaration name: (identifier) @function.name)

; 3. 捕获带 Receiver 的方法定义
(method_declaration name: (field_identifier) @method.name)

; 4. 捕获函数调用
(call_expression
  function: [
    (identifier) @call.function
    (selector_expression field: (field_identifier) @call.method)
  ])

; 5. 捕获导入依赖
(import_declaration (import_spec path: (interpreted_string_literal) @import.source))
`
