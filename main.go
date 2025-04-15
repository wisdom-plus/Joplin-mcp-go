package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
    // Create a new MCP server
    s := server.NewMCPServer(
        "Joplin MCP Server",
        "0.1.0",
        server.WithResourceCapabilities(true, true),
        server.WithLogging(),
        server.WithRecovery(),
    )

    // Add a joplin get note tool
    getnoteTool := mcp.NewTool("get note",
        mcp.WithDescription("Get a note"),
        mcp.WithString("note_id",
            mcp.Required(),
            mcp.Description("The ID of the note to get"),
        ),
    )

    // Add the tool to the server
    s.AddTool(getnoteTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        note_id := request.Params.Arguments["note_id"].(string)

        token := os.Getenv("JOPLIN_TOKEN")
        // get note from Joplin
        client := NewJoplinClient("", token)
        note, err := client.GetNote(note_id)
        if err != nil {
            log.Fatal(err)
        }

        return mcp.NewToolResultText(note.Body), nil
    })

    // Start the server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}
