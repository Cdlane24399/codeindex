package indexer_test

import (
	"testing"

	"github.com/01x-in/codeindex/internal/indexer"
	"github.com/stretchr/testify/assert"
)

func TestParseSwiftFunc(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "func generateId() -> String {\n    return UUID().uuidString\n}",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 24, Column: 4}, End: indexer.Position{Line: 26, Column: 5}},
			File:   "/repo/Sources/app/Service.swift",
			Lines:  "func generateId() -> String {",
			RuleID: "swift-func-def",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Service.swift", "swift")

	assert.Len(t, result.Nodes, 1)
	assert.Equal(t, "generateId", result.Nodes[0].Name)
	assert.Equal(t, "fn", result.Nodes[0].Kind)
	assert.False(t, result.Nodes[0].Exported, "func without public/open should not be exported")
	assert.Equal(t, "swift", result.Nodes[0].Language)
	assert.Equal(t, 25, result.Nodes[0].LineStart)
}

func TestParseSwiftStruct(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "public struct User {\n    public let id: String\n    public var name: String\n}",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 1, Column: 0}, End: indexer.Position{Line: 4, Column: 1}},
			File:   "/repo/Sources/app/Models.swift",
			Lines:  "public struct User {",
			RuleID: "swift-class-def",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Models.swift", "swift")

	assert.Len(t, result.Nodes, 1)
	assert.Equal(t, "User", result.Nodes[0].Name)
	assert.Equal(t, "class", result.Nodes[0].Kind)
	assert.True(t, result.Nodes[0].Exported)
	assert.Equal(t, "swift", result.Nodes[0].Language)
}

func TestParseSwiftClass(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "public class UserService {\n    private var users: [String: User] = [:]\n}",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 2, Column: 0}, End: indexer.Position{Line: 4, Column: 1}},
			File:   "/repo/Sources/app/Service.swift",
			Lines:  "public class UserService {",
			RuleID: "swift-class-def",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Service.swift", "swift")

	assert.Len(t, result.Nodes, 1)
	assert.Equal(t, "UserService", result.Nodes[0].Name)
	assert.Equal(t, "class", result.Nodes[0].Kind)
	assert.True(t, result.Nodes[0].Exported)
}

func TestParseSwiftEnum(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "public enum Status {\n    case active\n    case inactive\n}",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 14, Column: 0}, End: indexer.Position{Line: 17, Column: 1}},
			File:   "/repo/Sources/app/Models.swift",
			Lines:  "public enum Status {",
			RuleID: "swift-class-def",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Models.swift", "swift")

	assert.Len(t, result.Nodes, 1)
	assert.Equal(t, "Status", result.Nodes[0].Name)
	assert.Equal(t, "class", result.Nodes[0].Kind)
	assert.True(t, result.Nodes[0].Exported)
}

func TestParseSwiftExtensionSkipped(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "extension User {\n    func describe() -> String { return name }\n}",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 30, Column: 0}, End: indexer.Position{Line: 32, Column: 1}},
			RuleID: "swift-class-def",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Models.swift", "swift")
	assert.Len(t, result.Nodes, 0, "extensions do not define a new symbol and should be skipped")
}

func TestParseSwiftProtocol(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "public protocol Repository {\n    func findById(_ id: String) -> User?\n}",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 19, Column: 0}, End: indexer.Position{Line: 22, Column: 1}},
			File:   "/repo/Sources/app/Models.swift",
			Lines:  "public protocol Repository {",
			RuleID: "swift-protocol-def",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Models.swift", "swift")

	assert.Len(t, result.Nodes, 1)
	assert.Equal(t, "Repository", result.Nodes[0].Name)
	assert.Equal(t, "interface", result.Nodes[0].Kind)
	assert.True(t, result.Nodes[0].Exported)
	assert.Equal(t, "swift", result.Nodes[0].Language)
}

func TestParseSwiftImport(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "import Foundation",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 0, Column: 0}, End: indexer.Position{Line: 0, Column: 17}},
			File:   "/repo/Sources/app/Service.swift",
			Lines:  "import Foundation",
			RuleID: "swift-import",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/Service.swift", "swift")

	assert.Len(t, result.Edges, 1)
	assert.Equal(t, "imports", result.Edges[0].Kind)
	assert.Equal(t, "Foundation", result.Edges[0].TargetName)
	assert.Equal(t, "Sources/app/Service.swift", result.Edges[0].FilePath)
}

func TestParseSwiftExportDetection(t *testing.T) {
	t.Run("public func is exported", func(t *testing.T) {
		matches := []indexer.AstGrepMatch{
			{
				Text:   "public func add(_ user: User) {\n    users[user.id] = user\n}",
				Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 7, Column: 4}, End: indexer.Position{Line: 9, Column: 5}},
				RuleID: "swift-func-def",
			},
		}
		result := indexer.ParseMatches(matches, "Sources/app/Service.swift", "swift")
		assert.Len(t, result.Nodes, 1)
		assert.Equal(t, "add", result.Nodes[0].Name)
		assert.True(t, result.Nodes[0].Exported)
	})

	t.Run("internal func is not exported", func(t *testing.T) {
		matches := []indexer.AstGrepMatch{
			{
				Text:   "func generateId() -> String {\n    return UUID().uuidString\n}",
				Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 24, Column: 0}, End: indexer.Position{Line: 26, Column: 1}},
				RuleID: "swift-func-def",
			},
		}
		result := indexer.ParseMatches(matches, "Sources/app/Service.swift", "swift")
		assert.Len(t, result.Nodes, 1)
		assert.Equal(t, "generateId", result.Nodes[0].Name)
		assert.False(t, result.Nodes[0].Exported)
	})
}

func TestParseSwiftBuiltinCallsFiltered(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "print(\"hello\")",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 0, Column: 4}, End: indexer.Position{Line: 0, Column: 18}},
			RuleID: "swift-call-expr",
		},
		{
			Text:   "fatalError(\"boom\")",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 1, Column: 4}, End: indexer.Position{Line: 1, Column: 22}},
			RuleID: "swift-call-expr",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/main.swift", "swift")
	assert.Len(t, result.Edges, 0, "Swift built-in calls should be filtered out")
}

func TestParseSwiftCall(t *testing.T) {
	matches := []indexer.AstGrepMatch{
		{
			Text:   "service.createUser(name: \"Alice\", email: \"a@b.com\")",
			Range:  indexer.AstGrepRange{Start: indexer.Position{Line: 3, Column: 4}, End: indexer.Position{Line: 3, Column: 50}},
			RuleID: "swift-call-expr",
		},
	}

	result := indexer.ParseMatches(matches, "Sources/app/main.swift", "swift")
	assert.Len(t, result.Edges, 1)
	assert.Equal(t, "calls", result.Edges[0].Kind)
	assert.Equal(t, "service.createUser", result.Edges[0].TargetName)
}
