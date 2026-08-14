import type { ChatSource } from "./chat";

export type MessageRole = "user" | "assistant";

export interface Message {
    id: number | string;
    role: MessageRole;
    content: string;
    sources?: ChatSource[];
    status?: "error";
}