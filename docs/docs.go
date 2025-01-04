package docs

import "github.com/swaggo/swag"

const docTemplate = `
{
  "openapi": "3.0.1",
  "info": {
    "title": "USER-SERVICE API",
    "version": "1.0.1"
  },
  "servers": [
    {
      "url": "http://localhost:8080/",
      "description": "Local server"
    },
    {
      "url": "https://...",
      "description": "k8s dev"
    }
  ],
  "paths": {
    "/public/v1/auth/sign-up": {
      "post": {
        "summary": "регистрация пользователя",
        "tags": [
          "Auth"
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SignUpRequest"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/UserWithJWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "409": {
            "description": "Конфликт",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/public/v1/auth/sign-in": {
      "post": {
        "summary": "авторизация пользователя",
        "tags": [
          "Auth"
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SignInRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/UserWithJWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "409": {
            "description": "Конфликт",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/public/v1/auth/refresh/{token}": {
      "get": {
        "summary": "получение новых токенов",
        "tags": [
          "Auth"
        ],
        "parameters": [
          {
            "in": "path",
            "name": "token",
            "schema": {
              "type": "string",
              "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5"
            },
            "description": "рефреш-токен",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/JWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "409": {
            "description": "Конфликт",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/public/v1/users": {
      "get": {
        "security": [
          {
            "bearerAuth": []
          }
        ],
        "summary": "получение списка пользователей",
        "tags": [
          "Users"
        ],
        "parameters": [
          {
            "in": "query",
            "name": "role",
            "schema": {
              "type": "string",
              "enum": [
                "developer",
                "admin",
                "backend",
                "frontend",
                "designer",
                "devops",
                "project-manager"
              ]
            }
          },
          {
            "in": "query",
            "name": "sort",
            "schema": {
              "type": "string",
              "enum": [
                "name",
                "surname",
                "createdDate",
                "email"
              ]
            },
            "description": "Сортировка. Дефолтный - createDate",
            "required": false
          },
          {
            "in": "query",
            "name": "order",
            "schema": {
              "type": "string",
              "enum": [
                "asc",
                "desc"
              ]
            },
            "description": "Порядок сортировки. Дефолтный - desc",
            "required": false
          },
          {
            "in": "query",
            "name": "limit",
            "schema": {
              "type": "integer"
            },
            "required": false,
            "example": 20
          },
          {
            "in": "query",
            "name": "offset",
            "schema": {
              "type": "integer"
            },
            "required": false,
            "example": 0
          }
        ],
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/ListUsersWithoutJWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/public/v1/users/{id}": {
      "get": {
        "security": [
          {
            "bearerAuth": []
          }
        ],
        "summary": "получение пользователя",
        "tags": [
          "Users"
        ],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "schema": {
              "type": "string",
              "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5"
            },
            "description": "идентификатор пользователя",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/UserWithoutJWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      },
      "patch": {
        "security": [
          {
            "bearerAuth": []
          }
        ],
        "summary": "редактирование пользователя",
        "tags": [
          "Users"
        ],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "schema": {
              "type": "string",
              "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5"
            },
            "description": "идентификатор пользователя",
            "required": true
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UpdateUserRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/UserWithoutJWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      },
      "delete": {
        "security": [
          {
            "bearerAuth": []
          }
        ],
        "summary": "удаление пользователя",
        "tags": [
          "Users"
        ],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "schema": {
              "type": "string",
              "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5"
            },
            "description": "идентификатор пользователя",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/private/v1/users/{id}": {
      "patch": {
        "security": [
          {
            "bearerAuth": []
          }
        ],
        "summary": "редактирование пользователя (доступно только админам)",
        "tags": [
          "Users Private"
        ],
        "parameters": [
          {
            "in": "path",
            "name": "id",
            "schema": {
              "type": "string",
              "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5"
            },
            "description": "идентификатор пользователя",
            "required": true
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/UpdatePrivateUserRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Успешный ответ",
            "content": {
              "application/json": {
                "schema": {
                  "allOf": [
                    {
                      "$ref": "#/components/schemas/SuccessResponse"
                    },
                    {
                      "type": "object",
                      "properties": {
                        "result": {
                          "$ref": "#/components/schemas/UserPrivateWithoutJWT"
                        }
                      }
                    }
                  ]
                }
              }
            }
          },
          "400": {
            "description": "Не получилось обработать данные",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "401": {
            "description": "Не авторизован",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          },
          "500": {
            "description": "Внутренняя проблема сервера",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/ErrorResponse"
                }
              }
            }
          }
        }
      }
    },
    "/health/live": {
      "get": {
        "tags": [
          "HEALTH"
        ],
        "summary": "Health live",
        "responses": {
          "200": {
            "description": "Healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "string",
                  "example": "Healthy"
                }
              }
            }
          }
        }
      }
    },
    "/health/readiness": {
      "get": {
        "tags": [
          "HEALTH"
        ],
        "summary": "Health readiness",
        "responses": {
          "200": {
            "description": "Healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "string",
                  "example": "Healthy"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "ErrorResponse": {
        "properties": {
          "code": {
            "type": "integer",
            "example": 500,
            "nullable": false,
            "enum": [
              500,
              400,
              404,
              401,
              409
            ]
          },
          "result": {
            "type": "object",
            "example": null,
            "nullable": true
          },
          "error": {
            "type": "string",
            "example": "message error",
            "nullable": false
          },
          "errorType": {
            "type": "string",
            "nullable": false,
            "example": "message error type"
          }
        }
      },
      "SuccessResponse": {
        "properties": {
          "code": {
            "type": "integer",
            "example": 200,
            "nullable": false,
            "enum": [
              200,
              201
            ]
          },
          "result": {
            "type": "object",
            "nullable": false
          },
          "error": {
            "type": "string",
            "example": "",
            "nullable": false
          },
          "errorType": {
            "type": "string",
            "nullable": false,
            "example": ""
          }
        }
      },
      "SignInRequest": {
        "type": "object",
        "description": "модель авторизации пользователя",
        "properties": {
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": false
          },
          "password": {
            "type": "string",
            "description": "пароль",
            "example": "qwerty12345",
            "nullable": false
          }
        }
      },
      "SignUpRequest": {
        "type": "object",
        "description": "модель создания регистрации пользователя",
        "properties": {
          "name": {
            "type": "string",
            "description": "имя",
            "example": "German",
            "nullable": false
          },
          "surname": {
            "type": "string",
            "description": "фамилия",
            "example": "Bogatov",
            "nullable": false
          },
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": false
          },
          "password": {
            "type": "string",
            "description": "пароль",
            "example": "qwerty12345",
            "nullable": false
          }
        }
      },
      "UpdateUserRequest": {
        "type": "object",
        "description": "модель обновления пользователя",
        "properties": {
          "name": {
            "type": "string",
            "description": "имя",
            "example": "German",
            "nullable": true
          },
          "surname": {
            "type": "string",
            "description": "фамилия",
            "example": "Bogatov",
            "nullable": true
          },
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": true
          },
          "password": {
            "type": "string",
            "description": "пароль",
            "example": "newqwerty12345",
            "nullable": true
          }
        }
      },
      "UpdatePrivateUserRequest": {
        "type": "object",
        "description": "модель обновления пользователя",
        "properties": {
          "name": {
            "type": "string",
            "description": "имя",
            "example": "German",
            "nullable": true
          },
          "surname": {
            "type": "string",
            "description": "фамилия",
            "example": "Bogatov",
            "nullable": true
          },
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": true
          },
          "role": {
            "type": "string",
            "description": "роль",
            "enum": [
              "developer",
              "admin"
            ],
            "nullable": true
          }
        }
      },
      "JWT": {
        "type": "object",
        "description": "модель jwt токена",
        "properties": {
          "token": {
            "type": "string",
            "description": "идентификатор пользователя",
            "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VycyIsImV4cCI6MTcyNzU0OTUxOSwianRpIjoiYTEzZjA0ODQtZGU0Yi00MDdhLWI5NGQtMWUyMTk0YzQ4MTE3IiwiZW1haWwiOiJib2dhdG92Z3JtbkBtYWlsLnJ1Iiwicm9sZSI6ImRldmVsb3BlciJ9.pFJSZwFe0zpnmRuMQp7S4OvOknh1mTyuhfj6z0f4Afw",
            "nullable": false
          },
          "refreshToken": {
            "type": "string",
            "description": "идентификатор пользователя",
            "example": "909c6a00-76f1-491f-9b07-982705c6d68b",
            "nullable": false
          }
        }
      },
      "ListUsersWithoutJWT": {
        "type": "object",
        "description": "список пользователей",
        "properties": {
          "total": {
            "type": "integer",
            "description": "количество всех пользователей",
            "example": 100,
            "nullable": false
          },
          "items": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/UserWithoutJWT"
            }
          }
        }
      },
      "UserPrivateWithoutJWT": {
        "type": "object",
        "description": "модель пользователя без jwt",
        "properties": {
          "id": {
            "type": "string",
            "description": "идентификатор пользователя",
            "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5",
            "nullable": false
          },
          "name": {
            "type": "string",
            "description": "имя",
            "example": "German",
            "nullable": false
          },
          "surname": {
            "type": "string",
            "description": "фамилия",
            "example": "Bogatov",
            "nullable": false
          },
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": false
          },
          "role": {
            "type": "string",
            "description": "роль",
            "enum": [
              "developer",
              "admin"
            ],
            "nullable": false
          },
          "createdDate": {
            "type": "string",
            "description": "дата создания пользователя",
            "example": "2024-05-22T13:03:25Z",
            "nullable": false
          },
          "updatedDate": {
            "type": "string",
            "description": "дата изменения пользователя",
            "example": "2024-05-22T13:03:25Z",
            "nullable": true
          }
        }
      },
      "UserWithoutJWT": {
        "type": "object",
        "description": "модель пользователя без jwt",
        "properties": {
          "id": {
            "type": "string",
            "description": "идентификатор пользователя",
            "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5",
            "nullable": false
          },
          "name": {
            "type": "string",
            "description": "имя",
            "example": "German",
            "nullable": false
          },
          "surname": {
            "type": "string",
            "description": "фамилия",
            "example": "Bogatov",
            "nullable": false
          },
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": false
          },
          "role": {
            "type": "string",
            "description": "роль пользователя",
            "example": "developer",
            "nullable": false
          },
          "createdDate": {
            "type": "string",
            "description": "дата создания пользователя",
            "example": "2024-05-22T13:03:25Z",
            "nullable": false
          },
          "updatedDate": {
            "type": "string",
            "description": "дата изменения пользователя",
            "example": "2024-05-22T13:03:25Z",
            "nullable": true
          }
        }
      },
      "UserWithJWT": {
        "type": "object",
        "description": "модель пользователя с JWT",
        "properties": {
          "id": {
            "type": "string",
            "description": "идентификатор пользователя",
            "example": "c1cfe4b9-f7c2-423c-abfa-6ed1c05a15c5",
            "nullable": false
          },
          "name": {
            "type": "string",
            "description": "имя",
            "example": "German",
            "nullable": false
          },
          "surname": {
            "type": "string",
            "description": "фамилия",
            "example": "Bogatov",
            "nullable": false
          },
          "email": {
            "type": "string",
            "description": "электронная почта",
            "example": "bogatovgrmn@gmail.com",
            "nullable": false
          },
          "role": {
            "type": "string",
            "description": "роль пользователя",
            "example": "developer",
            "nullable": false
          },
          "createdDate": {
            "type": "string",
            "description": "дата создания пользователя",
            "example": "2024-05-22T13:03:25Z",
            "nullable": false
          },
          "updatedDate": {
            "type": "string",
            "description": "дата изменения пользователя",
            "example": "2024-05-22T13:03:25Z",
            "nullable": false
          },
          "jwt": {
            "type": "object",
            "description": "структура с jwt токеном и рефреш токеном",
            "nullable": false,
            "properties": {
              "token": {
                "type": "string",
                "description": "идентификатор пользователя",
                "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VycyIsImV4cCI6MTcyNzU0OTUxOSwianRpIjoiYTEzZjA0ODQtZGU0Yi00MDdhLWI5NGQtMWUyMTk0YzQ4MTE3IiwiZW1haWwiOiJib2dhdG92Z3JtbkBtYWlsLnJ1Iiwicm9sZSI6ImRldmVsb3BlciJ9.pFJSZwFe0zpnmRuMQp7S4OvOknh1mTyuhfj6z0f4Afw",
                "nullable": false
              },
              "refreshToken": {
                "type": "string",
                "description": "идентификатор пользователя",
                "example": "909c6a00-76f1-491f-9b07-982705c6d68b",
                "nullable": false
              }
            }
          }
        }
      }
    }
  }
}
`

var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "common-toxic-message service api",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
