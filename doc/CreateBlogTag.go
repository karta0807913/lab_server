package main

var CreateBlogTag Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Create",
    Fields: []Field{
        {
            Required: true,
            Comment: "",
            Name: "BlogID",
            Alias: "blog_id",
            Type: "uint",
        },{
            Required: true,
            Comment: "",
            Name: "TagID",
            Alias: "tag_id",
            Type: "uint",
        },
    },
}