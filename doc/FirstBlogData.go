package main

var FirstBlogData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Search",
    Fields: []Field{
        {
            Required: true,
            Comment: "",
            Name: "ID",
            Alias: "blog_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "Deleted",
            Alias: "deleted",
            Type: "uint",
        },
    },
}