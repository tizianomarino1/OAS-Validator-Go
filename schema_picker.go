package main


import (
"errors"
"fmt"
"strings"


"github.com/getkin/kin-openapi/openapi3"
)


// preferredName ritorna true per nomi di schema "preferiti" quando dobbiamo scegliere automaticamente.
func nomePref(name string) bool {
n := strings.ToLower(name)
return n == "instancedescriptor" || n == "body"
}


// getSchemaByName cerca uno schema in components.schemas per nome (case-insensitive).
func getSchemaByName(doc *openapi3.T, name string) *openapi3.SchemaRef {
if doc.Components.Schemas == nil {
return nil
}
if s, ok := doc.Components.Schemas[name]; ok && s != nil {
return s
}
lname := strings.ToLower(name)
for k, v := range doc.Components.Schemas {
if strings.ToLower(k) == lname {
return v
}
}
return nil
}


// findUniqueReqBodyJSONSchema scandisce tutte le operations e prova a trovare un unico schema di requestBody application/json.
func findUniqueReqBodyJSONSchema(doc *openapi3.T) *openapi3.SchemaRef {
found := map[*openapi3.SchemaRef]struct{}{}
var picked *openapi3.SchemaRef


collect := func(op *openapi3.Operation) {
if op == nil || op.RequestBody == nil || op.RequestBody.Value == nil {
return
}
if mt, ok := op.RequestBody.Value.Content["application/json"]; ok && mt != nil && mt.Schema != nil {
if _, seen := found[mt.Schema]; !seen {
found[mt.Schema] = struct{}{}
picked = mt.Schema
}
}
}


for _, p := range doc.Paths {
if p == nil { continue }
collect(p.Get)
collect(p.Post)
collect(p.Put)
collect(p.Patch)
collect(p.Delete)
collect(p.Options)
collect(p.Head)
collect(p.Trace)
}


if len(found) == 1 {
return picked
}
return nil
}


// findSchemaFromPathMethod estrae lo schema del requestBody application/json da una data path+method.
func findSchemaFromPathMethod(doc *openapi3.T, path, method string) *openapi3.SchemaRef {
if doc.Paths == nil { return nil }
item, ok := doc.Paths[path]
if !ok || item == nil { return nil }
var op *openapi3.Operation
switch strings.ToUpper(method) {
case "GET":
op = item.Get
case "POST":
op = item.Post
}