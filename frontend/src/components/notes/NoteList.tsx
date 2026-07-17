import type { Note } from "@/types/note";

interface NoteListProps {
    notes: Note[];
}

export default function NoteList({
    notes,
}: NoteListProps) {

    return (
        <section className="notes-grid">
            {notes.map((note) => (
                <article className="note-card" key={note.id}>
                    <h3>
                        {note.title}
                    </h3>

                    <p>
                        {note.content}
                    </p>
                </article>
            ))}
        </section>
    );
}