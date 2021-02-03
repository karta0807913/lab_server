package main

var CreateBlogData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Create",
    Fields: []Field{
        {
            Required: true,
            Comment: "",
            Name: "Title",
            Alias: "title",
            Type: "string",
        },{
            Required: true,
            Comment: "",
            Name: "Context",
            Alias: "context",
            Type: "string",
        },{
            Required: false,
            Comment: "",
            Name: "FileList",
            Alias: "file_list",
            Type: "*[]FileData",
        },{
            Required: false,
            Comment: "",
            Name: "TagList",
            Alias: "tag_list",
            Type: "*[]BlogTag",
        },
    },
}