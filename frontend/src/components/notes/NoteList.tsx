import type { Note } from "@/types/note";

interface NoteListProps {
    notes: Note[];
    deletingId: number | null;
    onDelete: (id: number) => void;
    onEdit: (note: Note) => void;
}

export default function NoteList({
    notes,
    deletingId,
    onDelete,
    onEdit,
}: NoteListProps) {

    if (notes.length === 0) {

        return (
            <p>
                No notes yet.
            </p>
        );
    }

    return (
        <section className="notes-grid">

            {notes.map((note) => {

                const isDeleting = deletingId === note.id;

                return (
                    <article
                        className={
                            isDeleting
                                ? "note-card note-card-deleting"
                                : "note-card"
                        }
                        key={note.id}
                    >

                        <div className="note-card-header">

                            <h3>{note.title}</h3>

                            <button
                                className="note-action-button"
                                type="button"
                                disabled={isDeleting}
                                onClick={() => onEdit(note)}
                            >
                                ✏️
                            </button>

                            <button
                                className={
                                    isDeleting
                                        ? "note-action-button note-action-button-busy"
                                        : "note-action-button"
                                }
                                type="button"
                                disabled={isDeleting}
                                onClick={() => onDelete(note.id)}
                            >
                                {isDeleting ? "Deleting..." : "🗑️"}
                            </button>

                        </div>

                        <p>{note.content}</p>

                    </article>
                );
            })}

        </section>
    );
}
