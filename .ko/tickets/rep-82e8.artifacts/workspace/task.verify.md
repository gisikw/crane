warning: Git tree '/home/dev/Projects/rep' is dirty
rep dev shell
  just test    — run tests
  just build   — build binary
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
=== RUN   TestClaudeAdapter_CaptureUsage
--- PASS: TestClaudeAdapter_CaptureUsage (0.00s)
=== RUN   TestCursorAdapter_CaptureUsage
--- PASS: TestCursorAdapter_CaptureUsage (0.00s)
=== RUN   TestOpencodeAdapter_CaptureUsage
--- PASS: TestOpencodeAdapter_CaptureUsage (0.00s)
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
=== RUN   TestParseArgs_ChargeCode
--- PASS: TestParseArgs_ChargeCode (0.00s)
=== RUN   TestParseArgs_ChargeCodeMissingValue
--- PASS: TestParseArgs_ChargeCodeMissingValue (0.00s)
=== RUN   TestParseArgs_ChargeCodeInvalidChars
--- PASS: TestParseArgs_ChargeCodeInvalidChars (0.00s)
=== RUN   TestParseArgs_DisallowedTools_Multiple
--- PASS: TestParseArgs_DisallowedTools_Multiple (0.00s)
=== RUN   TestStateDir_Default
--- PASS: TestStateDir_Default (0.00s)
=== RUN   TestStateDir_Custom
--- PASS: TestStateDir_Custom (0.00s)
=== RUN   TestAppendChargeRecord_CreatesAndAppends
--- PASS: TestAppendChargeRecord_CreatesAndAppends (0.00s)
=== RUN   TestProcessClaudeOutput_Full
--- PASS: TestProcessClaudeOutput_Full (0.00s)
=== RUN   TestProcessClaudeOutput_Fallback
--- PASS: TestProcessClaudeOutput_Fallback (0.00s)
=== RUN   TestProcessCursorOutput_Full
--- PASS: TestProcessCursorOutput_Full (0.00s)
=== RUN   TestProcessCursorOutput_FlatTokens
--- PASS: TestProcessCursorOutput_FlatTokens (0.00s)
=== RUN   TestProcessCursorOutput_Fallback
--- PASS: TestProcessCursorOutput_Fallback (0.00s)
=== RUN   TestProcessOpencodeOutput_Events
--- PASS: TestProcessOpencodeOutput_Events (0.00s)
=== RUN   TestProcessOpencodeOutput_Fallback
--- PASS: TestProcessOpencodeOutput_Fallback (0.00s)
=== RUN   TestValidateChargeCode_Valid
--- PASS: TestValidateChargeCode_Valid (0.00s)
=== RUN   TestValidateChargeCode_Invalid
--- PASS: TestValidateChargeCode_Invalid (0.00s)
=== RUN   TestResolveProvider_Explicit
--- PASS: TestResolveProvider_Explicit (0.00s)
=== RUN   TestResolveProvider_NoProviders
--- PASS: TestResolveProvider_NoProviders (0.00s)
=== RUN   TestProviderBinaries
--- PASS: TestProviderBinaries (0.00s)
=== RUN   TestRollup_Empty
--- PASS: TestRollup_Empty (0.00s)
=== RUN   TestRollup_SingleRecord
--- PASS: TestRollup_SingleRecord (0.00s)
=== RUN   TestRollup_MultipleProviderModel
--- PASS: TestRollup_MultipleProviderModel (0.00s)
=== RUN   TestRollup_NilTokenFields
--- PASS: TestRollup_NilTokenFields (0.00s)
=== RUN   TestReadChargeRecords_FileNotExist
--- PASS: TestReadChargeRecords_FileNotExist (0.00s)
=== RUN   TestReadChargeRecords_ReadsAndParsesLines
--- PASS: TestReadChargeRecords_ReadsAndParsesLines (0.00s)
=== RUN   TestReadChargeRecords_SkipsBadLines
--- PASS: TestReadChargeRecords_SkipsBadLines (0.00s)
PASS
ok  	github.com/gisikw/rep	0.006s
