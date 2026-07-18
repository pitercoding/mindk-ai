import type { Note } from "@/types/note";

interface NoteListProps {
    notes: Note[];
    onDelete: (id: number) => void;
    onEdit: (note: Note) => void;
}

export default function NoteList({
    notes,
    onDelete,
    onEdit,
}: NoteListProps) {

    return (
        <section className="notes-grid">

            {notes.map((note) => (
                <article
                    className="note-card"
                    key={note.id}
                >

                    <div className="note-card-header">

                        <h3>{note.title}</h3>

                        <button
                            className="note-action-button"
                            type="button"
                            onClick={() => onEdit(note)}
                        >
                            ✏️
                        </button>

                        <button
                            className="note-action-button"
                            type="button"
                            onClick={() => onDelete(note.id)}
                        >
                            🗑️
                        </button>

                    </div>

                    <p>{note.content}</p>

                </article>
            ))}

        </section>
    );
}
