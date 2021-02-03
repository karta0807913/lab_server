package main

var CreateUserData Document = Document{
    Path: "/",
    Comment: "",
	Mode: "Create",
    Fields: []Field{
        {
            Required: true,
            Comment: "",
            Name: "Nickname",
            Alias: "nickname",
            Type: "string",
        },{
            Required: true,
            Comment: "",
            Name: "Account",
            Alias: "-",
            Type: "string",
        },{
            Required: true,
            Comment: "",
            Name: "Password",
            Alias: "-",
            Type: "string",
        },{
            Required: true,
            Comment: "",
            Name: "Status",
            Alias: "-",
            Type: "uint",
        },
    },
}