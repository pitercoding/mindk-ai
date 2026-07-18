import { apiClient } from "@/api/client";
import type { Note } from "@/types/note";

export async function createNote(note: Omit<Note, "id">): Promise<Note> {

    return apiClient<Note>("/notes", {
        method: "POST",
        body: JSON.stringify(note),
    });
}

export async function getNotes(): Promise<Note[]> {
    return apiClient<Note[]>("/notes");
}

export async function deleteNote(id: number): Promise<void> {
    await apiClient(`/notes/${id}`, { method: "DELETE", },);
}

export async function updateNote(
    id: number,
    note: Omit<Note, "id">
): Promise<Note> {

    return apiClient<Note>(`/notes/${id}`, {
        method: "PUT",
        body: JSON.stringify(note),
    }
    );
}