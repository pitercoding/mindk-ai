export type ChatMode = "knowledge" | "note";

export interface ChatSession {
    id: number;
    title: string;
    mode: ChatMode;
    note_id?: number;
    created_at: string;
    updated_at: string;
}
