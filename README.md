# CCNext

> Claude Code 图形化终端包装器 — 语音输入 · 反卡住检测 · CC Switch 集成

CCNext 是一个 Windows 桌面应用，用图形界面包装 Claude Code CLI，提供：

- **终端仿真** — 通过 xterm.js 完整渲染 Claude Code 的 TUI 界面
- **语音输入** — 内置 Web Speech API，支持语音转文字
- **反卡住检测** — 注入随机口令判定任务是否完成，超时自动续接或提醒
- **快捷按键栏** — 一键模拟 ESC / Enter / Space / Shift+Tab / Ctrl+O / Ctrl+B / 上下方向键
- **CC Switch 集成** — 直接复用 CC Switch 的全局 API 配置，无需额外设置

## 技术栈

| 层面 | 技术 |
|---|---|
| 语言 | Go |
| GUI 框架 | Wails v2 |
| 前端 | Svelte 4 + TypeScript + xterm.js |
| 终端桥接 | Windows ConPTY |
| 语音 | Web Speech API |
| 存储 | SQLite (消息历史) + JSON (配置) |

## 快速开始

### 环境要求

- Windows 10/11
- Go 1.22+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- [Claude Code](https://claude.ai/code) CLI (在 PATH 中)
- [CC Switch](https://github.com/CyanBukkit/CC-Next) (处理 API 配置)

### 开发模式

```bash
cd frontend && npm install
cd .. && wails dev
```

### 生产构建

```bash
wails build -platform windows/amd64
# 输出: build/bin/ccnext.exe
```

## 架构

```
┌──────────────────────────────────────────┐
│              CCNext GUI                   │
│    Wails WebView2 (Svelte + xterm.js)    │
│                                           │
│  ┌────────────────────────────────────┐  │
│  │         xterm.js 终端              │  │
│  │    (Claude Code TUI 完整渲染)      │  │
│  └────────────────────────────────────┘  │
│  ┌────────────────────────────────────┐  │
│  │  快捷栏 + 输入框 + 语音按钮         │  │
│  └────────────────────────────────────┘  │
└──────────────┬───────────────────────────┘
               │ Wails Bridge
┌──────────────┴───────────────────────────┐
│              Go 后端                       │
│                                           │
│  ConPTY ←→ claude 进程                   │
│  口令注入 → 流式检测 → 超时判定           │
│  SQLite 会话/消息持久化                    │
└──────────────────────────────────────────┘
```

## 反卡住机制

每次用户发消息时，CCNext 在后台：

1. 生成随机口令（如 `CuAyX4XKsy`）
2. 在消息末尾追加隐藏指令：*"完成后必须以这个短语结尾：`CuAyX4XKsy`"*
3. 监控 PTY 输出，检测口令是否出现
4. `/` 开头的命令不注入、不检测
5. Agent 运行中的进度条自动续命超时计时器

| 场景 | 行为 |
|---|---|
| 口令出现 | 任务完成 |
| 超时无口令 + 无 agent 进度 | 判定卡住 → 弹窗 / 自动续接 |
| 超时但检测到 agent 进度 | 自动重置计时器 |
| `/` 斜杠指令 | 完全跳过 |

## 快捷键栏

```
[返回] [回车] [空格] | [后台] [模式] [思考] [▲] [▼] | [📁] [🔄]
  ESC   Enter  Space   Ctrl+B Sh+Tab Ctrl+O  ↑   ↓   目录  刷新
```

所有按钮直接发送键盘信号到 Claude Code TUI。

## 设置

- **卡住超时** — 30 秒 ~ 30 分钟（默认 5 分钟）
- **自动续接模式** — 检测到卡住时自动发 continue
- **隐藏指令模板** — 支持 `%random%` 变量
- **深色/浅色主题**

## 项目结构

```
ccnext/
├── main.go                  # Wails 入口
├── app.go                   # 核心业务逻辑 + 前端绑定
├── wails.json               # Wails 配置
│
├── internal/
│   ├── claude/manager.go    # ConPTY 管理 + ANSI 清理
│   ├── config/              # 配置持久化
│   ├── passphrase/          # 口令生成 + 流式匹配
│   ├── stuck/               # 超时检测 + 恢复处理
│   └── message/             # 消息路由 + SQLite 存储
│
├── frontend/
│   └── src/
│       ├── App.svelte              # 主应用
│       ├── components/
│       │   ├── TerminalView.svelte # xterm.js 终端
│       │   ├── InputBar.svelte     # 输入栏 + 快捷按键
│       │   ├── VoiceButton.svelte  # 语音输入
│       │   ├── StatusBar.svelte    # 状态栏
│       │   ├── StuckBanner.svelte  # 卡住提醒
│       │   └── SettingsPanel.svelte# 设置面板
│       └── lib/                    # 工具库
│
└── scripts/build.bat        # 构建脚本
```

## License

MIT
