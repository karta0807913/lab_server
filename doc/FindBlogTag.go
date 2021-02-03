package main

var FindBlogTag Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Search",
    Fields: []Field{
        {
            Required: false,
            Comment: "",
            Name: "BlogID",
            Alias: "blog_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "BlogTagID",
            Alias: "blog_tag_id",
            Type: "uint",
        },
    },
}