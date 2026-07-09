import {
  SendMessage,
  GetSettings,
  UpdateSettings,
  GetStatus,
  RestartClaude,
  GetHistory,
  ClearHistory,
  ListSessions,
  SwitchSession,
  NewSession,
  DeleteSession,
  GetCurrentSession,
  ResizePTY,
  SendPTYKey,
  PickWorkDir,
  GetActiveMode,
  SetActiveMode,
} from "../../wailsjs/go/main/App";

import type { config, main, message } from "../../wailsjs/go/models";

export {
  EventsOn,
  EventsOff,
  EventsOffAll,
} from "../../wailsjs/runtime/runtime";

// Re-export auto-generated types
export type Settings = config.Settings;
export type SessionInfo = main.SessionInfo;
export type Message = message.Message;

export async function sendMessage(text: string): Promise<void> {
  return SendMessage(text);
}

export async function getSettings(): Promise<config.Settings> {
  return GetSettings();
}

export async function updateSettings(s: config.Settings): Promise<void> {
  return UpdateSettings(s);
}

export async function getStatus(): Promise<string> {
  return GetStatus();
}

export async function restartClaude(): Promise<void> {
  return RestartClaude();
}

export async function getHistory(): Promise<message.Message[]> {
  return GetHistory();
}

export async function clearHistory(): Promise<string> {
  return ClearHistory();
}

export async function listSessions(): Promise<main.SessionInfo[]> {
  return ListSessions();
}

export async function switchSession(sessionID: string): Promise<message.Message[]> {
  return SwitchSession(sessionID);
}

export async function newSession(title: string): Promise<main.SessionInfo> {
  return NewSession(title);
}

export async function deleteSession(sessionID: string): Promise<void> {
  return DeleteSession(sessionID);
}

export async function getCurrentSession(): Promise<main.SessionInfo> {
  return GetCurrentSession();
}

export function resizeTerminal(cols: number, rows: number): void {
  ResizePTY(cols, rows);
}

// PTY key sequences for Claude Code shortcuts
export const Keys = {
  esc:      '\x1b',
  enter:    '\r',
  space:    ' ',
  shiftTab: '\x1b[Z',
  ctrlL:    '\x0c',
  ctrlO:    '\x0f',
  ctrlB:    '\x02',
  up:       '\x1b[A',
  down:     '\x1b[B',
} as const;

export function sendPTYKey(seq: string): void {
  SendPTYKey(seq);
}

export async function pickWorkDir(): Promise<string> {
  return PickWorkDir();
}

export async function getActiveMode(): Promise<string> {
  return GetActiveMode();
}

export async function setActiveMode(mode: string): Promise<void> {
  return SetActiveMode(mode);
}
