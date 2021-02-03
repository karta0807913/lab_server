package main

var UpdateUserData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Updates",
    Fields: []Field{
        {
            Required: true,
            Comment: " gorm.Model",
            Name: "ID",
            Alias: "user_id",
            Type: "uint",
        },{
            Required: false,
            Comment: "",
            Name: "Nickname",
            Alias: "nickname",
            Type: "string",
        },{
            Required: false,
            Comment: "",
            Name: "Account",
            Alias: "-",
            Type: "string",
        },{
            Required: false,
            Comment: "",
            Name: "Password",
            Alias: "-",
            Type: "string",
        },{
            Required: false,
            Comment: "",
            Name: "IsAdmin",
            Alias: "is_admin",
            Type: "bool",
        },{
            Required: false,
            Comment: "",
            Name: "Status",
            Alias: "-",
            Type: "uint",
        },
    },
}