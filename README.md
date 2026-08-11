# `cscope_lsp`

This LSP implementation uses cscope to get results quickly.

## Installation

```shell
go install github.com/dhananjaylatkar/cscope_lsp@latest
```

## Neovim config

```lua
local function start_cscope_lsp()
  local root_files =
    { "cscope.out", "cscope.files", "cscope.in.out", "cscope.out.in", "cscope.out.po", "cscope.po.out" }
  local paths = vim.fs.find(root_files, { stop = vim.env.HOME })
  local root_dir = vim.fs.dirname(paths[1])

  if root_dir then
    vim.lsp.start({
      name = "cscope_lsp",
      cmd = { "cscope_lsp" },
      root_dir = root_dir,
      filetypes = { "c", "h", "cpp", "hpp" },
    })
  end
end

vim.api.nvim_create_autocmd("FileType", {
  pattern = { "c", "h", "cpp", "hpp" },
  desc = "Start cscope_lsp",
  callback = start_cscope_lsp,
})
```

## Helix config

```toml
[language-server.cscope_lsp]
command = "cscope_lsp"
config = { replaceMethods = true }
```

## Requirements

1. `cscope` is installed.
2. `cscope.out` is created and updated.

## Supported Capabilities

1. `textDocument/definition` or `textDocument/typeDefinition`
2. `textDocument/references` or `textDocument/typeDefinition`

If the server is configured with
```json
{ "replaceMethods": true }
```
It will replace the following methods:

1. `textDocument/definition` with `textDocument/typeDefinition`
2. `textDocument/references` with `textDocument/implementation`

There will be conflict with clangd anyway, because clangd implements
all four methods, however, since cscope is mostly used for functions,
methods, etc. `typeDefinition` won't break much.

however, `implementation` for cscope references may not be the best,
but it is what it is...

## Thanks

Used [educationalsp](https://github.com/tjdevries/educationalsp) as starter template.
