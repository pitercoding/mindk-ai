import { useEffect, useState } from "react";

import NoteList from "@/components/notes/NoteList";

import { getNotes } from "@/services/noteService";

import type { Note } from "@/types/note";

export default function NotesPage() {

    const [notes, setNotes] = useState<Note[]>([]);

    async function loadNotes() {
        try {
            const response = await getNotes();

            setNotes(response);

        } catch (error) {
            console.error(
                "Failed to load notes:",
                error
            );
        }
    }

    useEffect(() => {
        loadNotes();
    }, []);

    return (
        <main>
            <h1>Notes</h1>

            <NoteList notes={notes} />
        </main>
    );
}