package main

var CreateTagInfo Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Create",
    Fields: []Field{
        {
            Required: true,
            Comment: "",
            Name: "Name",
            Alias: "name",
            Type: "string",
        },
    },
}