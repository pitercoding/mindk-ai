import { useEffect, useState } from "react";

import ConfirmDialog from "@/components/common/ConfirmDialog";
import NoteList from "@/components/notes/NoteList";
import NoteForm from "@/components/notes/NoteForm";

import { useToast } from "@/context/toast-context";

import { createNote, deleteNote, getNotes, updateNote } from "@/services/noteService";

import type { Note } from "@/types/note";

export default function NotesPage() {

    const toast = useToast();

    const [notes, setNotes] = useState<Note[]>([]);
    const [selectedNote, setSelectedNote] = useState<Note | null>(null);

    const [isLoading, setIsLoading] = useState(true);
    const [loadError, setLoadError] = useState<string | null>(null);
    const [isSaving, setIsSaving] = useState(false);
    const [deletingId, setDeletingId] = useState<number | null>(null);
    const [noteToDelete, setNoteToDelete] = useState<Note | null>(null);

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

            toast.success("Note created.");

            return true;

        } catch (error) {
            console.error(
                "Failed to create note:",
                error
            );

            toast.error("Failed to create note.");

            return false;

        } finally {
            setIsSaving(false);
        }
    }

    function requestDeleteNote(id: number) {
        setNoteToDelete(notes.find((note) => note.id === id) ?? null);
    }

    async function handleDeleteNote(note: Note) {

        setNoteToDelete(null);
        setDeletingId(note.id);

        try {

            await deleteNote(note.id);

            setNotes(previous =>
                previous.filter(item => item.id !== note.id)
            );

            toast.success("Note deleted.");

        } catch (error) {
            console.error("Failed to delete note:", error,);

            toast.error("Failed to delete note.");

        } finally {
            setDeletingId(null);
        }
    }

    async function handleUpdateNote(
        id: number,
        note: Omit<Note, "id">
    ): Promise<boolean> {

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

            toast.success("Note updated.");

            return true;

        } catch (error) {

            console.error(
                "Failed to update note",
                error
            );

            toast.error("Failed to update note.");

            return false;

        } finally {
            setIsSaving(false);
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
                    isSaving={isSaving}
                    onCreated={handleCreateNote}
                    onUpdated={handleUpdateNote}
                />

            </section>

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
                            onDelete={requestDeleteNote}
                            onEdit={handleEditNote}
                        />
                    )
                }

            </section>

            {noteToDelete && (
                <ConfirmDialog
                    title="Delete note?"
                    message={`"${noteToDelete.title}" will be permanently deleted.`}
                    onConfirm={() => handleDeleteNote(noteToDelete)}
                    onCancel={() => setNoteToDelete(null)}
                />
            )}

        </main>
    );
}
