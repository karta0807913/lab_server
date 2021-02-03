package main

var UpdateTagInfo Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Updates",
    Fields: []Field{
        {
            Required: true,
            Comment: "",
            Name: "ID",
            Alias: "id",
            Type: "uint",
        },{
            Required: true,
            Comment: "",
            Name: "Name",
            Alias: "name",
            Type: "string",
        },
    },
}