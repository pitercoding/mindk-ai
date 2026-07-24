export type MessageRole = "user" | "assistant";

export interface Message {
    id: number | string;
    role: MessageRole;
    content: string;
}