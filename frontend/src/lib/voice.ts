export interface SpeechRecognitionAlternative {
  transcript: string;
  confidence: number;
}

export interface SpeechRecognitionResult {
  [index: number]: SpeechRecognitionAlternative;
  isFinal: boolean;
  length: number;
}

export interface SpeechRecognitionResultList {
  [index: number]: SpeechRecognitionResult;
  length: number;
}

export interface SpeechRecognitionEvent extends Event {
  resultIndex: number;
  results: SpeechRecognitionResultList;
}

export interface SpeechRecognitionErrorEvent extends Event {
  error: string;
  message: string;
}

export interface SpeechRecognition extends EventTarget {
  lang: string;
  continuous: boolean;
  interimResults: boolean;
  maxAlternatives: number;
  onresult: ((event: SpeechRecognitionEvent) => void) | null;
  onerror: ((event: SpeechRecognitionErrorEvent) => void) | null;
  onend: ((event: Event) => void) | null;
  start(): void;
  stop(): void;
  abort(): void;
}

export interface SpeechRecognitionConstructor {
  new (): SpeechRecognition;
}

declare global {
  interface Window {
    SpeechRecognition?: SpeechRecognitionConstructor;
    webkitSpeechRecognition?: SpeechRecognitionConstructor;
  }
}

let recognition: SpeechRecognition | null = null;
let resolveListening: ((value: string) => void) | null = null;
let rejectListening: ((reason: Error) => void) | null = null;
let interimCallback: ((text: string) => void) | null = null;

export function isSupported(): boolean {
  return "SpeechRecognition" in window || "webkitSpeechRecognition" in window;
}

export function setInterimResultCallback(callback: (text: string) => void): void {
  interimCallback = callback;
}

export function startListening(): Promise<string> {
  return new Promise((resolve, reject) => {
    if (!isSupported()) {
      reject(new Error("Speech recognition is not supported in this browser."));
      return;
    }

    const Constructor = window.SpeechRecognition || window.webkitSpeechRecognition;
    if (!Constructor) {
      reject(new Error("Speech recognition constructor is not available."));
      return;
    }

    recognition = new Constructor();
    recognition.lang = navigator.language || "en-US";
    recognition.continuous = true;
    recognition.interimResults = true;
    recognition.maxAlternatives = 1;

    let finalText = "";

    recognition.onresult = (event: SpeechRecognitionEvent) => {
      let interim = "";
      for (let i = event.resultIndex; i < event.results.length; i++) {
        const result = event.results[i];
        const transcript = result[0].transcript;
        if (result.isFinal) {
          finalText += transcript;
        } else {
          interim += transcript;
        }
      }
      if (interimCallback) {
        interimCallback(finalText + interim);
      }
    };

    recognition.onerror = (event: SpeechRecognitionErrorEvent) => {
      cleanup();
      reject(new Error(`Speech recognition error: ${event.error}`));
    };

    recognition.onend = () => {
      cleanup();
      resolve(finalText);
    };

    resolveListening = resolve;
    rejectListening = reject;
    recognition.start();
  });
}

export function stopListening(): void {
  if (recognition) {
    recognition.stop();
  }
}

function cleanup(): void {
  recognition = null;
  resolveListening = null;
  rejectListening = null;
}
