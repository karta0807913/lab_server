package main

var FirstUserData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Search",
    Fields: []Field{
        {
            Required: false,
            Comment: " gorm.Model",
            Name: "ID",
            Alias: "user_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "Account",
            Alias: "-",
            Type: "string",
        },
    },
}