import EmptyState from "@/components/common/EmptyState";
import { EditIcon, NotesIcon, PlusIcon, TrashIcon } from "@/components/common/icons";

import { countWords, stripTitleHeading, toPlainText } from "@/utils/text";

import type { Note } from "@/types/note";

interface NoteListProps {
    notes: Note[];
    deletingId: number | null;
    onDelete: (id: number) => void;
    onEdit: (note: Note) => void;
    onCreate: () => void;
}

export default function NoteList({
    notes,
    deletingId,
    onDelete,
    onEdit,
    onCreate,
}: NoteListProps) {

    if (notes.length === 0) {

        return (
            <EmptyState
                icon={<NotesIcon size={20} />}
                title="No notes yet"
                description="Capture ideas, meeting notes or study material. Everything you write can be searched and used in chat."
                actions={
                    <button
                        type="button"
                        className="btn btn-primary"
                        onClick={onCreate}
                    >
                        <PlusIcon size={16} />
                        Write your first note
                    </button>
                }
            />
        );
    }

    return (
        <section className="notes-grid">

            {notes.map((note) => {

                const isDeleting = deletingId === note.id;

                const preview = toPlainText(
                    stripTitleHeading(note.content, note.title),
                );
                const words = countWords(preview);

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
                                className="btn btn-secondary btn-icon"
                                type="button"
                                aria-label={`Edit ${note.title}`}
                                title="Edit"
                                disabled={isDeleting}
                                onClick={() => onEdit(note)}
                            >
                                <EditIcon size={16} />
                            </button>

                            <button
                                className={
                                    isDeleting
                                        ? "btn btn-secondary btn-sm"
                                        : "btn btn-secondary btn-icon btn-danger-hover"
                                }
                                type="button"
                                aria-label={`Delete ${note.title}`}
                                title="Delete"
                                disabled={isDeleting}
                                onClick={() => onDelete(note.id)}
                            >
                                {isDeleting ? "Deleting..." : <TrashIcon size={16} />}
                            </button>

                        </div>

                        {preview ? (
                            <p className="card-preview">{preview}</p>
                        ) : (
                            <p className="card-preview card-preview-empty">
                                No content yet
                            </p>
                        )}

                        <footer className="card-meta">
                            {words} {words === 1 ? "word" : "words"}
                        </footer>

                    </article>
                );
            })}

        </section>
    );
}
