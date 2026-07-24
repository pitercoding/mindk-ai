export interface ChatMessage {
    id: number;
    note_id: number;
    role: "user" | "assistant";
    content: string;
    created_at: string;
}


export interface ChatContext {
    note_id: number;
    title: string;
    content: string;
}

export interface ChatRequest {
    message: string;
    context?: ChatContext;
}

export interface ChatResponse {
    answer: string;
}
