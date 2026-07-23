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
