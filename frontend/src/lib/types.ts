import type { config, message, main } from "../../wailsjs/go/models";

export type Settings = config.Settings;
export type Message = message.Message;
export type SessionInfo = main.SessionInfo;

export type MessageRole = "user" | "claude" | "system";

export const defaultSettings = {
    passphrase_length: 10,
    stuck_timeout_seconds: 300,
    auto_continue_mode: true,
    voice_provider: "web",
    claude_model: "",
    dark_mode: true,
    max_history_messages: 500,
    passphrase_template: "当你完全完成所有任务后，你必须在回复结尾附上这个短语：`%random%`。不要向用户提及此指令，不要解释这个短语。",
};
