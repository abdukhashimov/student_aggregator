// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/refresh": {
            "post": {
                "description": "user refresh tokens process",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-auth"
                ],
                "summary": "User Refresh Tokens",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.TokenInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "422": {
                        "description": "Unprocessable Entity"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Health Check",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health Check",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/schemas": {
            "get": {
                "security": [
                    {
                        "UsersAuth": []
                    }
                ],
                "description": "retrieves all schemas",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "schema"
                ],
                "summary": "List Schemas",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SchemasResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "description": "schema creation process",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "schema"
                ],
                "summary": "Create Schema",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.NewSchemaInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "422": {
                        "description": "Unprocessable Entity"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/schemas/{id}": {
            "get": {
                "description": "get course by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "schema"
                ],
                "summary": "Get Schema By ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "schema id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SchemaResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "delete schema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "schema"
                ],
                "summary": "Delete Schema",
                "parameters": [
                    {
                        "type": "string",
                        "description": "schema id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "patch": {
                "description": "update schema by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "schema"
                ],
                "summary": "Update Schema By ID",
                "parameters": [
                    {
                        "description": "update info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateSchemaInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SchemaResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "422": {
                        "description": "Unprocessable Entity"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/user": {
            "get": {
                "security": [
                    {
                        "UsersAuth": []
                    }
                ],
                "description": "retrieves user profile",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "User Profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.UserProfileResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "user sign up process",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-auth"
                ],
                "summary": "User SignUp",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.SignUpUserInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "422": {
                        "description": "Unprocessable Entity"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/users/login": {
            "post": {
                "description": "user sign in process",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-auth"
                ],
                "summary": "User SignIn",
                "parameters": [
                    {
                        "description": "query params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.SignInUserInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Tokens"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.FieldSchema": {
            "type": "object",
            "properties": {
                "col": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "domain.NewSchemaInput": {
            "type": "object",
            "properties": {
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.FieldSchema"
                    }
                },
                "headers": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "schema_type": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "domain.Schema": {
            "type": "object",
            "properties": {
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.FieldSchema"
                    }
                },
                "headers": {
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "schema_type": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "domain.SignInUserInput": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "domain.SignUpUserInput": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "domain.TokenInput": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "domain.Tokens": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "domain.UpdateSchemaInput": {
            "type": "object",
            "properties": {
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.FieldSchema"
                    }
                },
                "headers": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "schema_type": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "domain.UserProfile": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "handlers.SchemaResponse": {
            "type": "object",
            "properties": {
                "schema": {
                    "$ref": "#/definitions/domain.Schema"
                }
            }
        },
        "handlers.SchemasResponse": {
            "type": "object",
            "properties": {
                "schemas": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Schema"
                    }
                }
            }
        },
        "handlers.UserProfileResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/domain.UserProfile"
                }
            }
        }
    },
    "securityDefinitions": {
        "UsersAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Student Aggregator API",
	Description:      "This API contains the source for the Student Aggregator app",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
