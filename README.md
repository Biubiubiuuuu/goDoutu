# goDoutu
收集、整理各类表情包，直接保存，让斗图更简单

# vs code 基本的代码提示配置
```
{
    "go.autocompleteUnimportedPackages": true,
    "go.useCodeSnippetsOnFunctionSuggest": true,
    "go.inferGopath": true,
    "go.gopath":"H:\\GoWork",
    "go.useCodeSnippetsOnFunctionSuggestWithoutType": true,
    "go.formatTool": "gofmt",
    "go.vetOnSave": "package",
    "go.buildOnSave": "package",
    "go.coverOnSave": true,
    "go.gocodeFlags": [
        "-builtin",
        "-ignore-case",
        "-unimported-packages"
    ],
    "go.gocodePackageLookupMode": "go",
    "go.gotoSymbol.includeGoroot": true,
    "go.gotoSymbol.includeImports": true,
    "go.addTags": {
        "tags": "json",
        "options": "json=omitempty",
        "promptForTags": true,
        "transform": "snakecase"
    },
    "go.useLanguageServer": true
}
```