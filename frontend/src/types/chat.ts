export interface ChatMessage {
    id: number;
    session_id: number;
    role: "user" | "assistant";
    content: string;
    created_at: string;
}

export interface ChatRequest {
    session_id: number;
    message: string;

    // Used only when session_id is 0, to auto-create a session.
    mode?: string;
    note_id?: number;
    title?: string;
}

export interface ChatSource {
    document_id: number;
    name: string;
    score: number;
}

export interface ChatResponse {
    answer: string;
    session_id: number;
    sources?: ChatSource[];
}
