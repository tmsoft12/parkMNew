// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/createcar": {
            "post": {
                "description": "Registers a new car entering the parking lot",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Create a new car entry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Parking spot number",
                        "name": "parkno",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Car details",
                        "name": "car",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/modelscar.Car_Model"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created car details",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Invalid request or car already inside",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Database error",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/getallcars": {
            "get": {
                "description": "Get list of cars with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Get list of cars",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 5,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Parking spot number",
                        "name": "parkno",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.GetCarsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/getcar/{id}": {
            "get": {
                "description": "Get a car by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Get a car by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Car ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/modelscar.Car_Model"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Authenticates a user and starts a session.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Login User",
                "parameters": [
                    {
                        "description": "User Login Data",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usercontrol.LoginInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Login successful",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "message: Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "message: Invalid credentials",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "message: Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/logout": {
            "post": {
                "description": "Ends the session of a logged-in user.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Logout User",
                "responses": {
                    "200": {
                        "description": "message: Logout successful",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "message: Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Creates a new user and stores their hashed password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register User",
                "parameters": [
                    {
                        "description": "User Registration Data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/modelsuser.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "message: User Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "message: Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "message: Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/searchcar": {
            "get": {
                "description": "Search for a car by plate number, parking number, and other optional filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Search for a car by plate number and optional filters",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Car plate number",
                        "name": "car_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Parking spot number",
                        "name": "parkno",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Car status (Inside, Exited)",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 5,
                        "description": "Number of items per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.GetCarsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/updatecar/{plate}": {
            "put": {
                "description": "Updates a car's status and calculates payment and duration based on start and end times.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Update a car by plate number",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Car plate number",
                        "name": "plate",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Car details to update",
                        "name": "car",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/modelscar.Car_Model"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Updated car details",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Car already exited or invalid request",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Car not found",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Error parsing time",
                        "schema": {
                            "$ref": "#/definitions/carcontrol.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "carcontrol.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "carcontrol.GetCarsResponse": {
            "type": "object",
            "properties": {
                "cars": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/modelscar.Car_Model"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "page": {
                    "type": "integer"
                },
                "total_count": {
                    "type": "integer"
                }
            }
        },
        "modelscar.Car_Model": {
            "type": "object",
            "properties": {
                "car_number": {
                    "type": "string"
                },
                "duration": {
                    "type": "integer"
                },
                "end_time": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image_url": {
                    "type": "string"
                },
                "park_no": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "total_payment": {
                    "type": "number"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "modelsuser.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "usercontrol.LoginInput": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "192.168.100.192:3000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Airline API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
