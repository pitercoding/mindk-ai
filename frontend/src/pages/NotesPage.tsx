import { useEffect, useState } from "react";

import NoteList from "@/components/notes/NoteList";
import NoteForm from "@/components/notes/NoteForm";

import { createNote, deleteNote, getNotes, updateNote } from "@/services/noteService";

import type { Note } from "@/types/note";

export default function NotesPage() {

    const [notes, setNotes] = useState<Note[]>([]);
    const [selectedNote, setSelectedNote] = useState<Note | null>(null);

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
                createdNote,
                ...previousNotes,
            ]);

        } catch (error) {
            console.error(
                "Failed to create note:",
                error
            );
        }
    }

    async function handleDeleteNote(id: number) {

        try {

            if (!confirm("Delete this note?")) {
                return;
            }

            await deleteNote(id);

            setNotes(previous =>
                previous.filter(note => note.id !== id)
            );

        } catch (error) {
            console.error("Failed to delete note:", error,);
        }
    }

    async function handleUpdateNote(
        id: number,
        note: Omit<Note, "id">
    ) {

        try {

            const updated = await updateNote(
                id,
                note
            );

            setNotes(previous =>
                previous.map(item =>
                    item.id === id
                        ? updated
                        : item
                )
            );

            setSelectedNote(null);

        } catch (error) {

            console.error(
                "Failed to update note",
                error
            );
        }
    }

    function handleEditNote(note: Note) {
        setSelectedNote(note);

        window.scrollTo({
            top: 0,
            behavior: "smooth",
        });
    }

    return (
        <main className="notes-page">

            <header className="notes-header">
                <h1>Notes</h1>
            </header>


            <section className="notes-form-section">

                <NoteForm
                    selectedNote={selectedNote}
                    onCreated={handleCreateNote}
                    onUpdated={handleUpdateNote}
                />

            </section>


            <section className="notes-list-section">

                <h2>
                    Your Notes
                </h2>

                <NoteList
                    notes={notes}
                    onDelete={handleDeleteNote}
                    onEdit={handleEditNote}
                />

            </section>

        </main>
    );
}
