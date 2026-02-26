go test -v -count=1 ./...
=== RUN   TestClaudeAdapter_Build
--- PASS: TestClaudeAdapter_Build (0.00s)
=== RUN   TestClaudeAdapter_NoPermissions
--- PASS: TestClaudeAdapter_NoPermissions (0.00s)
=== RUN   TestCursorAdapter_InlinesSystemPrompt
--- PASS: TestCursorAdapter_InlinesSystemPrompt (0.00s)
=== RUN   TestOpencodeAdapter_Build
--- PASS: TestOpencodeAdapter_Build (0.00s)
=== RUN   TestOpencodeAdapter_NoModel
--- PASS: TestOpencodeAdapter_NoModel (0.00s)
=== RUN   TestGetAdapter_Defaults
--- PASS: TestGetAdapter_Defaults (0.00s)
=== RUN   TestClaudeAdapter_AllowedTools
--- PASS: TestClaudeAdapter_AllowedTools (0.00s)
=== RUN   TestClaudeAdapter_DisallowedTools
--- PASS: TestClaudeAdapter_DisallowedTools (0.00s)
=== RUN   TestClaudeAdapter_ToolsIgnoredWithAllowAll
--- PASS: TestClaudeAdapter_ToolsIgnoredWithAllowAll (0.00s)
=== RUN   TestExitCode
--- PASS: TestExitCode (0.00s)
=== RUN   TestParseArgs_PromptFromPositional
--- PASS: TestParseArgs_PromptFromPositional (0.00s)
=== RUN   TestParseArgs_Flags
--- PASS: TestParseArgs_Flags (0.00s)
=== RUN   TestParseArgs_AllowAllDefault
--- PASS: TestParseArgs_AllowAllDefault (0.00s)
=== RUN   TestParseArgs_WithPermissions
--- PASS: TestParseArgs_WithPermissions (0.00s)
=== RUN   TestParseArgs_UnknownFlag
--- PASS: TestParseArgs_UnknownFlag (0.00s)
=== RUN   TestParseArgs_EmptyPrompt
--- PASS: TestParseArgs_EmptyPrompt (0.00s)
=== RUN   TestParseArgs_AllowedTools_Single
--- PASS: TestParseArgs_AllowedTools_Single (0.00s)
=== RUN   TestParseArgs_AllowedTools_Multiple_CommaSeparated
--- PASS: TestParseArgs_AllowedTools_Multiple_CommaSeparated (0.00s)
=== RUN   TestParseArgs_AllowedTools_Multiple_SpaceSeparated
--- PASS: TestParseArgs_AllowedTools_Multiple_SpaceSeparated (0.00s)
=== RUN   TestParseArgs_DisallowedTools_Single
--- PASS: TestParseArgs_DisallowedTools_Single (0.00s)
=== RUN   TestParseArgs_DisallowedTools_Multiple
--- PASS: TestParseArgs_DisallowedTools_Multiple (0.00s)
=== RUN   TestResolveProvider_Explicit
--- PASS: TestResolveProvider_Explicit (0.00s)
=== RUN   TestResolveProvider_NoProviders
--- PASS: TestResolveProvider_NoProviders (0.00s)
=== RUN   TestProviderBinaries
--- PASS: TestProviderBinaries (0.00s)
PASS
ok  	github.com/gisikw/rep	0.003s
