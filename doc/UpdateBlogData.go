package main

var UpdateBlogData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Updates",
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
            Name: "Title",
            Alias: "title",
            Type: "string",
        },{
            Required: false,
            Comment: "",
            Name: "OwnerID",
            Alias: "user_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "Owner",
            Alias: "owner",
            Type: "*UserData",
        },{
            Required: false,
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
        },{
            Required: false,
            Comment: "",
            Name: "Deleted",
            Alias: "deleted",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "CreatedAt",
            Alias: "create_time",
            Type: "time.Time",
        },
    },
}