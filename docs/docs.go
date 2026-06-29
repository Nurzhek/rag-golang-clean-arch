// Package docs provides the OpenAPI (Swagger 2.0) specification for the RAG
// Backend API and registers it with swaggo at init time, so the /docs endpoint
// can serve interactive documentation.
//
// This file mirrors what `swag init` generates from the @-annotations on the
// handlers and the main package. Regenerate it after changing those annotations:
//
//	swag init -g cmd/server/main.go -o docs
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
        "/api/v1/documents": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "documents"
                ],
                "summary": "List ingested documents",
                "description": "Returns all ingested documents with their chunk counts and metadata.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ListDocumentsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "text/plain"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "documents"
                ],
                "summary": "Upload a file for chunked async ingestion",
                "description": "Accepts a raw request body (the file content) and queues chunked background ingestion. Query parameters are stored as document metadata.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Optional metadata: source name",
                        "name": "source",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Optional metadata: title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "description": "Raw file content",
                        "name": "file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/dto.JobAcceptedResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "413": {
                        "description": "Request Entity Too Large",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "documents"
                ],
                "summary": "Queue inline document ingestion",
                "description": "Accepts inline text content, queues asynchronous embedding + storage, and returns a job ID to poll.",
                "parameters": [
                    {
                        "description": "Document content and optional metadata",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.IngestRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.JobAcceptedResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/documents/{id}": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "documents"
                ],
                "summary": "Delete a document",
                "description": "Deletes a document and all of its chunks by source document ID.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Document ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.DeleteResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/jobs/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jobs"
                ],
                "summary": "Poll ingestion job status",
                "description": "Returns the current status and progress of an asynchronous ingestion job.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Job ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.JobResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/query": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "query"
                ],
                "summary": "Ask a question (RAG)",
                "description": "Embeds the question, retrieves the most relevant chunks, and generates a grounded answer with sources.",
                "parameters": [
                    {
                        "description": "Question and optional top_k",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.QueryRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.QueryResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/dto.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Liveness check",
                "description": "Reports service liveness.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.DeleteResponse": {
            "type": "object",
            "properties": {
                "deleted_chunks": {
                    "type": "integer"
                },
                "document_id": {
                    "type": "string"
                }
            }
        },
        "dto.DocumentSummary": {
            "type": "object",
            "properties": {
                "chunk_count": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            }
        },
        "dto.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "dto.IngestRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            }
        },
        "dto.JobAcceptedResponse": {
            "type": "object",
            "properties": {
                "job_id": {
                    "type": "string"
                },
                "poll_url": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "dto.JobResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "document_id": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "processed_chunks": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "total_chunks": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "dto.ListDocumentsResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "documents": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.DocumentSummary"
                    }
                }
            }
        },
        "dto.QueryRequest": {
            "type": "object",
            "properties": {
                "question": {
                    "type": "string"
                },
                "top_k": {
                    "type": "integer"
                }
            }
        },
        "dto.QueryResponse": {
            "type": "object",
            "properties": {
                "answer": {
                    "type": "string"
                },
                "sources": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.Source"
                    }
                }
            }
        },
        "dto.Source": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "document_id": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "score": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds the exported Swagger spec so it can be referenced and, if
// needed, mutated at runtime (e.g. to set Host).
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "RAG Backend API",
	Description:      "Retrieval-Augmented Generation backend (Clean Architecture, langchaingo + OpenAI). Ingest documents, run chunked async uploads, poll job status, and ask grounded questions.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
