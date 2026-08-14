import { useEffect, useState } from "react";

import NoteList from "@/components/notes/NoteList";
import NoteForm from "@/components/notes/NoteForm";

import { createNote, deleteNote, getNotes, updateNote } from "@/services/noteService";

import type { Note } from "@/types/note";

export default function NotesPage() {

    const [notes, setNotes] = useState<Note[]>([]);
    const [selectedNote, setSelectedNote] = useState<Note | null>(null);

    const [isLoading, setIsLoading] = useState(true);
    const [loadError, setLoadError] = useState<string | null>(null);
    const [actionError, setActionError] = useState<string | null>(null);
    const [isSaving, setIsSaving] = useState(false);
    const [deletingId, setDeletingId] = useState<number | null>(null);

    async function loadNotes() {
        setIsLoading(true);
        setLoadError(null);

        try {
            const response = await getNotes();

            setNotes(response);

        } catch (error) {
            console.error(
                "Failed to load notes:",
                error
            );

            setLoadError("Failed to load notes.");

        } finally {
            setIsLoading(false);
        }
    }

    useEffect(() => {
        loadNotes();
    }, []);

    async function handleCreateNote(
        note: Omit<Note, "id">
    ): Promise<boolean> {

        setActionError(null);
        setIsSaving(true);

        try {

            const createdNote = await createNote({
                title: note.title,
                content: note.content,
            });

            setNotes((previousNotes) => [
                createdNote,
                ...previousNotes,
            ]);

            return true;

        } catch (error) {
            console.error(
                "Failed to create note:",
                error
            );

            setActionError("Failed to create note.");

            return false;

        } finally {
            setIsSaving(false);
        }
    }

    async function handleDeleteNote(id: number) {

        if (!confirm("Delete this note?")) {
            return;
        }

        setActionError(null);
        setDeletingId(id);

        try {

            await deleteNote(id);

            setNotes(previous =>
                previous.filter(note => note.id !== id)
            );

        } catch (error) {
            console.error("Failed to delete note:", error,);

            setActionError("Failed to delete note.");

        } finally {
            setDeletingId(null);
        }
    }

    async function handleUpdateNote(
        id: number,
        note: Omit<Note, "id">
    ): Promise<boolean> {

        setActionError(null);
        setIsSaving(true);

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

            return true;

        } catch (error) {

            console.error(
                "Failed to update note",
                error
            );

            setActionError("Failed to update note.");

            return false;

        } finally {
            setIsSaving(false);
        }
    }

    function handleEditNote(note: Note) {
        setSelectedNote(note);
        setActionError(null);

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
                    isSaving={isSaving}
                    onCreated={handleCreateNote}
                    onUpdated={handleUpdateNote}
                />

            </section>

            {actionError && (
                <p className="notes-error">
                    {actionError}
                </p>
            )}

            <section className="notes-list-section">

                <h2>
                    Your Notes
                </h2>

                {
                    isLoading ? (
                        <p>
                            Loading notes...
                        </p>
                    ) : loadError ? (
                        <p className="notes-error">
                            {loadError}
                        </p>
                    ) : (
                        <NoteList
                            notes={notes}
                            deletingId={deletingId}
                            onDelete={handleDeleteNote}
                            onEdit={handleEditNote}
                        />
                    )
                }

            </section>

        </main>
    );
}
