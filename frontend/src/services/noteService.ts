import { apiClient } from "@/api/client";
import type { Note } from "@/types/note";

export async function createNote(note: Omit<Note, "id">): Promise<Note> {

    return apiClient<Note>(
        "/notes",
        { method: "POST", body: JSON.stringify(note) }
    );
}