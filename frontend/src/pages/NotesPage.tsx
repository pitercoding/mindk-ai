import { useEffect, useState } from "react";

import NoteList from "@/components/notes/NoteList";
import NoteForm from "@/components/notes/NoteForm";

import { createNote, getNotes } from "@/services/noteService";

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

    async function handleCreateNote(note: Omit<Note, "id">) {

        try {

            const createdNote = await createNote({
                title: note.title,
                content: note.content,
            });

            setNotes((previousNotes) => [
                ...previousNotes,
                createdNote,
            ]);

        } catch (error) {
            console.error(
                "Failed to create note:",
                error
            );
        }
    }

    return (
        <main className="notes-page">

            <header className="notes-header">
                <h1>Notes</h1>
            </header>


            <section className="notes-form-section">

                <NoteForm
                    onCreated={handleCreateNote}
                />

            </section>


            <section className="notes-list-section">

                <h2>
                    Your Notes
                </h2>

                <NoteList
                    notes={notes}
                />

            </section>

        </main>
    );
}