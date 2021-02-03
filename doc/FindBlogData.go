package main

var FindBlogData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Search",
    Fields: []Field{
        {
            Required: false,
            Comment: "",
            Name: "Title",
            Alias: "title",
            Type: "string",
        },{
            Required: false,
            Comment: "",
            Name: "OwnerID",
            Alias: "user_id",
            Type: "uint",
        },
    },
}