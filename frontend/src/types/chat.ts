export interface ChatContext {
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
