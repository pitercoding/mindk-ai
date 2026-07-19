import type { ChatContext } from "./chat-context";

export interface ChatRequest {
    message: string;
    context?: ChatContext;
}

export interface ChatResponse {
    answer: string;
}
