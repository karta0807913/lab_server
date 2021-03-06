package main

var FindFileData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Search",
    Fields: []Field{
        {
            Required: false,
            Comment: " gorm.Model",
            Name: "ID",
            Alias: "file_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "UserID",
            Alias: "user_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "BlogID",
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